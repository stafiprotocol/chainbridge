// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package ethereum

import (
	"errors"
	"math/big"
	"time"

	"github.com/ChainSafe/log15"
	eth "github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stafiprotocol/chainbridge-utils/blockstore"
	metrics "github.com/stafiprotocol/chainbridge-utils/metrics/types"
	"github.com/stafiprotocol/chainbridge/bindings/Bridge"
	"github.com/stafiprotocol/chainbridge/bindings/ERC20Handler"
	"github.com/stafiprotocol/chainbridge/chains"
	utils "github.com/stafiprotocol/chainbridge/shared/ethereum"
)

var BlockDelay = big.NewInt(10)
var BlockRetryInterval = time.Second * 5
var BlockRetryLimit = 5
var ErrFatalPolling = errors.New("listener block polling failed")

type listener struct {
	cfg                    Config
	conn                   Connection
	router                 chains.Router
	bridgeContract         *Bridge.Bridge // instance of bound bridge contract
	erc20HandlerContract   *ERC20Handler.ERC20Handler
	log                    log15.Logger
	blockstore             blockstore.Blockstorer
	stop                   <-chan int
	sysErr                 chan<- error // Reports fatal error to core
	latestBlock            metrics.LatestBlock
	metrics                *metrics.ChainMetrics
}

// NewListener creates and returns a listener
func NewListener(conn Connection, cfg *Config, log log15.Logger, bs blockstore.Blockstorer, stop <-chan int, sysErr chan<- error, m *metrics.ChainMetrics) *listener {
	return &listener{
		cfg:         *cfg,
		conn:        conn,
		log:         log,
		blockstore:  bs,
		stop:        stop,
		sysErr:      sysErr,
		latestBlock: metrics.LatestBlock{LastUpdated: time.Now()},
		metrics:     m,
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
	var currentBlock = l.cfg.startBlock
	var retry = BlockRetryLimit
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

			// Sleep if the difference is less than BlockDelay; (latest - current) < BlockDelay
			if big.NewInt(0).Sub(latestBlock, currentBlock).Cmp(BlockDelay) == -1 {
				l.log.Debug("Block not ready, will retry", "target", currentBlock, "latest", latestBlock)
				time.Sleep(BlockRetryInterval)
				continue
			}

			// Write to block store. Not a critical operation, no need to retry
			err = l.blockstore.StoreBlock(currentBlock)
			if err != nil {
				l.log.Error("Failed to write latest block to blockstore", "block", currentBlock, "err", err)
			}

			// Goto next block and reset retry counter
			currentBlock.Add(currentBlock, big.NewInt(1))
			l.latestBlock.Height = big.NewInt(0).Set(latestBlock)
			l.latestBlock.LastUpdated = time.Now()
			retry = BlockRetryLimit
			if l.metrics != nil {
				l.metrics.BlocksProcessed.Inc()
			}
		}
	}
}

// buildQuery constructs a query for the bridgeContract by hashing sig to get the event topic
func buildQuery(contract ethcommon.Address, sig utils.EventSig, startBlock *big.Int, endBlock *big.Int) eth.FilterQuery {
	query := eth.FilterQuery{
		FromBlock: startBlock,
		ToBlock:   endBlock,
		Addresses: []ethcommon.Address{contract},
		Topics: [][]ethcommon.Hash{
			{sig.GetTopic()},
		},
	}
	return query
}
