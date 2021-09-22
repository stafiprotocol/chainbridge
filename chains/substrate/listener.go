// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stafiprotocol/chainbridge/config"

	"github.com/ChainSafe/log15"
	"github.com/itering/substrate-api-rpc/websocket"
	"github.com/stafiprotocol/chainbridge/chains"
	"github.com/stafiprotocol/chainbridge/utils/blockstore"
	"github.com/stafiprotocol/chainbridge/utils/msg"
)

type listener struct {
	name          string
	chainId       msg.ChainId
	startBlock    uint64
	blockstore    blockstore.Blockstorer
	conn          *Connection
	subscriptions map[eventName]eventHandler // Handlers for specific events
	router        chains.Router
	log           log15.Logger
	stop          <-chan int
	sysErr        chan<- error
	wsconn        websocket.WsConn
	decimals      map[string]decimal.Decimal
}

var (
	// Frequency of polling for a new block
	BlockRetryInterval = 15 * time.Second
	BlockRetryLimit    = 35

	EventRetryLimit    = 20
	EventRetryInterval = 100 * time.Millisecond
)

func NewListener(conn *Connection, name string, id msg.ChainId, startBlock uint64, log log15.Logger, bs blockstore.Blockstorer, stop <-chan int, sysErr chan<- error, decimals map[string]decimal.Decimal) *listener {
	return &listener{
		name:          name,
		chainId:       id,
		startBlock:    startBlock,
		blockstore:    bs,
		conn:          conn,
		subscriptions: make(map[eventName]eventHandler),
		log:           log,
		stop:          stop,
		sysErr:        sysErr,
		decimals:      decimals,
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

	for _, sub := range Subscriptions {
		err := l.registerEventHandler(sub.name, sub.handler)
		if err != nil {
			return err
		}
	}

	go func() {
		err = l.pollBlocks()
		if err != nil {
			l.log.Error("Polling blocks failed", "err", err)
		}
	}()

	return nil
}

// registerEventHandler enables a handler for a given event. This cannot be used after Start is called.
func (l *listener) registerEventHandler(name eventName, handler eventHandler) error {
	if l.subscriptions[name] != nil {
		return fmt.Errorf("event %s already registered", name)
	}
	l.subscriptions[name] = handler
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

	for _, evt := range evts {
		if evt.ModuleId != config.BridgeCommon || evt.EventId != config.FungibleTransferEventId {
			continue
		}

		data, err := FungibleTransferEventData(evt, l.decimals)
		if err != nil {
			return err
		}

		if l.subscriptions[FungibleTransfer] != nil {
			l.submitMessage(l.subscriptions[FungibleTransfer](data, l.log))
		}
	}

	return nil
}

// submitMessage inserts the chainId into the msg and sends it to the router
func (l *listener) submitMessage(m msg.Message, err error) {
	if err != nil {
		l.log.Error("Critical error processing event", "err", err)
		return
	}
	m.Source = l.chainId
	err = l.router.Send(m)
	if err != nil {
		l.log.Error("failed to process event", "err", err)
	}
}
