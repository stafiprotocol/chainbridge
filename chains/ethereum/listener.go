// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package ethereum

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/stafiprotocol/chainbridge/utils/msg"
	"math/big"
	"time"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge/bindings/Bridge"
	"github.com/stafiprotocol/chainbridge/bindings/ERC20Handler"
	"github.com/stafiprotocol/chainbridge/chains"
	ethconn "github.com/stafiprotocol/chainbridge/connections/ethereum"
	utils "github.com/stafiprotocol/chainbridge/shared/ethereum"
	"github.com/stafiprotocol/chainbridge/utils/blockstore"
)

var (
	BlockDelay         = big.NewInt(10)
	BlockRetryInterval = time.Second * 15
	BlockRetryLimit    = 20
	ErrFatalPolling    = errors.New("listener block polling failed")
	logInterval        = uint64(100)
)

type listener struct {
	cfg                  *ethconn.Config
	conn                 Connection
	router               chains.Router
	bridgeContract       *Bridge.Bridge // instance of bound bridge contract
	erc20HandlerContract *ERC20Handler.ERC20Handler
	log                  log15.Logger
	blockstore           blockstore.Blockstorer
	stop                 <-chan int
	sysErr               chan<- error // Reports fatal error to core
}

// NewListener creates and returns a listener
func NewListener(conn Connection, cfg *ethconn.Config, log log15.Logger, bs blockstore.Blockstorer, stop <-chan int, sysErr chan<- error) *listener {
	return &listener{
		cfg:        cfg,
		conn:       conn,
		log:        log,
		blockstore: bs,
		stop:       stop,
		sysErr:     sysErr,
	}
}

// setContracts sets the listener with the appropriate contracts
func (l *listener) setContracts(bridge *Bridge.Bridge, erc20Handler *ERC20Handler.ERC20Handler) {
	l.bridgeContract = bridge
	l.erc20HandlerContract = erc20Handler
}

// sets the router
func (l *listener) setRouter(r chains.Router) {
	l.router = r
}

// start registers all subscriptions provided by the config
func (l *listener) start() error {
	l.log.Debug("Starting listener...")

	go func() {
		err := l.pollBlocks()
		if err != nil {
			l.log.Error("Polling blocks failed", "err", err)
		}
	}()

	return nil
}

// pollBlocks will poll for the latest block and proceed to parse the associated events as it sees new blocks.
// Polling begins at the block defined in `l.cfg.startBlock`. Failed attempts to fetch the latest block or parse
// a block will be retried up to BlockRetryLimit times before continuing to the next block.
func (l *listener) pollBlocks() error {
	l.log.Info("Polling Blocks...")
	var currentBlock = l.cfg.StartBlock()
	var retry = BlockRetryLimit
	if l.cfg.ChainId() == 3 {
		logInterval = 200
	}
	for {
		select {
		case <-l.stop:
			return errors.New("polling terminated")
		default:
			// No more retries, goto next block
			if retry == 0 {
				l.log.Error("Polling failed, retries exceeded")
				l.sysErr <- ErrFatalPolling
				return nil
			}

			latestBlock, err := l.conn.LatestBlock()
			if err != nil {
				l.log.Error("Unable to get latest block", "block", currentBlock, "err", err)
				retry--
				time.Sleep(BlockRetryInterval)
				continue
			}

			if currentBlock.Uint64()%logInterval == 0 {
				l.log.Debug("pollBlocks", "target", currentBlock, "latest", latestBlock)
			}

			// Sleep if the difference is less than BlockDelay; (latest - current) < BlockDelay
			if big.NewInt(0).Sub(latestBlock, currentBlock).Cmp(BlockDelay) == -1 {
				time.Sleep(BlockRetryInterval)
				continue
			}

			// Parse out events
			err = l.getDepositEventsForBlock(currentBlock)
			if err != nil {
				l.log.Error("Failed to get events for block", "block", currentBlock, "err", err)
				retry--
				continue
			}

			// Write to block store. Not a critical operation, no need to retry
			err = l.blockstore.StoreBlock(currentBlock)
			if err != nil {
				l.log.Error("Failed to write latest block to blockstore", "block", currentBlock, "err", err)
			}

			// Goto next block and reset retry counter
			currentBlock.Add(currentBlock, big.NewInt(1))
			retry = BlockRetryLimit
		}
	}
}

// getDepositEventsForBlock looks for the deposit event in the latest block
func (l *listener) getDepositEventsForBlock(latestBlock *big.Int) error {
	query := buildQuery(l.cfg.BridgeContract(), utils.Deposit, latestBlock, latestBlock)

	// querying for logs
	logs, err := l.conn.Client().FilterLogs(context.Background(), query)
	if err != nil {
		return fmt.Errorf("unable to Filter Logs: %s", err)
	}

	// read through the log events and handle their deposit event if handler is recognized
	for _, log := range logs {
		var m msg.Message
		destId := msg.ChainId(log.Topics[1].Big().Uint64())
		rId := msg.ResourceIdFromSlice(log.Topics[2].Bytes())
		nonce := msg.Nonce(log.Topics[3].Big().Uint64())

		//skip event if not support chainId
		if !l.router.SupportChainId(destId) {
			continue
		}

		addr, err := l.bridgeContract.ResourceIDToHandlerAddress(&bind.CallOpts{From: l.conn.Keypair().CommonAddress()}, rId)
		if err != nil {
			return fmt.Errorf("failed to get handler from resource ID %x", rId)
		}

		if addr == l.cfg.Erc20HandlerContract() {
			m, err = l.handleErc20DepositedEvent(destId, nonce)
		} else {
			l.log.Error("event has unrecognized handler", "handler", addr.Hex())
			return nil
		}

		if err != nil {
			return err
		}

		err = l.router.Send(m)
		if err != nil {
			l.log.Error("subscription error: failed to route message", "err", err)
		}
	}

	return nil
}
