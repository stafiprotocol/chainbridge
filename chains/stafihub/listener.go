// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package stafihub

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/shopspring/decimal"

	"github.com/ChainSafe/log15"
	"github.com/cosmos/cosmos-sdk/types"
	stafiHubXBridgeTypes "github.com/stafihub/stafihub/x/bridge/types"
	"github.com/stafiprotocol/chainbridge/chains"
	"github.com/stafiprotocol/chainbridge/utils/blockstore"
	"github.com/stafiprotocol/chainbridge/utils/msg"
)

var (
	ErrEventAttributeNumberUnMatch = errors.New("ErrEventAttributeNumberTooFew")
)

type listener struct {
	name       string
	chainId    msg.ChainId
	startBlock uint64
	blockstore blockstore.Blockstorer
	conn       *Connection
	router     chains.Router
	log        log15.Logger
	stop       <-chan int
	sysErr     chan<- error
	decimals   map[string]decimal.Decimal
}

var (
	// Frequency of polling for a new block
	BlockRetryInterval = 6 * time.Second
	BlockRetryLimit    = 100

	EventRetryLimit    = 20
	EventRetryInterval = 100 * time.Millisecond
)

func NewListener(conn *Connection, name string, id msg.ChainId, startBlock uint64, log log15.Logger,
	bs blockstore.Blockstorer, stop <-chan int, sysErr chan<- error, decimals map[string]decimal.Decimal) *listener {
	return &listener{
		name:       name,
		chainId:    id,
		startBlock: startBlock,
		blockstore: bs,
		conn:       conn,
		log:        log,
		stop:       stop,
		sysErr:     sysErr,
		decimals:   decimals,
	}
}

func (l *listener) setRouter(r chains.Router) {
	l.router = r
}

// Start creates the initial subscription for all events
func (l *listener) start() error {
	// Check whether latest is less than starting block
	latest, err := l.conn.LatestBlockNumber()
	if err != nil {
		return err
	}

	if latest < l.startBlock {
		return fmt.Errorf("starting block (%d) is greater than latest known block (%d)", l.startBlock, latest)
	}

	go func() {
		err = l.pollBlocks()
		if err != nil {
			l.log.Error("Polling blocks failed", "err", err)
		}
	}()

	return nil
}

// pollBlocks will poll for the latest block and proceed to parse the associated events as it sees new blocks.
// Polling begins at the block defined in `l.startBlock`. Failed attempts to fetch the latest block or parse
// a block will be retried up to BlockRetryLimit times before returning with an error.
func (l *listener) pollBlocks() error {
	var currentBlock = l.startBlock
	var retry = BlockRetryLimit
	for {
		select {
		case <-l.stop:
			return errors.New("terminated")
		default:
			// No more retries, goto next block
			if retry == 0 {
				l.sysErr <- fmt.Errorf("event polling retries exceeded (chain=%d, name=%s)", l.chainId, l.name)
				return nil
			}

			finalized, err := l.conn.FinalizedBlockNumber()
			if err != nil {
				return err
			}

			// Sleep if the block we want comes after the most recently finalized block
			if currentBlock > finalized {
				if currentBlock%100 == 0 {
					l.log.Trace("Block not yet finalized", "target", currentBlock, "finalized", finalized)
				}
				time.Sleep(BlockRetryInterval)
				continue
			}

			err = l.processEvents(currentBlock)
			if err != nil {
				l.log.Error("Failed to process events in block", "block", currentBlock, "err", err)
				retry--
				continue
			}

			// Write to blockstore
			err = l.blockstore.StoreBlock(big.NewInt(0).SetUint64(currentBlock))
			if err != nil {
				l.log.Error("Failed to write to blockstore", "err", err)
			}

			currentBlock++
			retry = BlockRetryLimit
		}
	}
}

// processEvents fetches a block and parses out the events, calling Listener.handleEvents()
func (l *listener) processEvents(blockNum uint64) error {
	if blockNum%100 == 0 {
		l.log.Debug("processEvents", "blockNum", blockNum)
	}

	evts, err := l.conn.GetEvents(blockNum)
	if err != nil {
		l.log.Warn("processEvents GetEvents error, will retry", "err", err)
		for i := 0; i < EventRetryLimit; i++ {
			time.Sleep(EventRetryInterval)
			evts, err = l.conn.GetEvents(blockNum)
			if err == nil {
				break
			}
		}
		if err != nil {
			return err
		}
	}

	for _, tx := range evts {
		for _, log := range tx.Logs {
			for _, event := range log.Events {
				err := l.processStringEvents(event, int64(blockNum))
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// submitMessage inserts the chainId into the msg and sends it to the router
func (l *listener) submitMessage(m msg.Message) (err error) {
	if err != nil {
		l.log.Error("Critical error processing event", "err", err)
		return
	}
	m.Source = l.chainId
	err = l.router.Send(m)
	if err != nil {
		l.log.Error("failed to process event", "err", err)
	}
	return err
}

func (l *listener) processStringEvents(event types.StringEvent, blockNumber int64) error {
	l.log.Debug("processStringEvents", "event", event)
	switch {
	case event.Type == stafiHubXBridgeTypes.EventTypeDeposit:
		if len(event.Attributes) != 5 {
			return ErrEventAttributeNumberUnMatch
		}
		chainId, err := strconv.Atoi(event.Attributes[0].Value)
		if err != nil {
			return err
		}
		resourceBts, err := hex.DecodeString(event.Attributes[1].Value)
		if err != nil {
			return err
		}
		var resource [32]byte
		copy(resource[:], resourceBts)

		depositNonce, err := strconv.Atoi(event.Attributes[2].Value)
		if err != nil {
			return err
		}
		amount, ok := types.NewIntFromString(event.Attributes[3].Value)
		if !ok {
			return fmt.Errorf("amount format err")
		}

		receiver, err := types.AccAddressFromBech32(event.Attributes[4].Value)
		if err != nil {
			return err
		}

		m := msg.NewFungibleTransfer(
			l.chainId,
			msg.ChainId(chainId),
			msg.Nonce(depositNonce),
			amount.BigInt(),
			resource,
			receiver,
		)

		l.log.Info("find event", "msg", m, "block number", blockNumber)
		return l.submitMessage(m)

	default:
		return nil
	}
}
