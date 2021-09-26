// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package ethereum

import (
	"fmt"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge/bindings/Bridge"
	ethconn "github.com/stafiprotocol/chainbridge/connections/ethereum"
	"github.com/stafiprotocol/chainbridge/utils/msg"
)

const (
	TransferredStatus uint8 = 2
	CancelledStatus   uint8 = 3

	msgLimit = 4096
)

type writer struct {
	cfg            *ethconn.Config
	conn           Connection
	bridgeContract *Bridge.Bridge // instance of bound receiver bridgeContract
	log            log15.Logger
	msgChan        chan msg.Message
	stop           <-chan int
	sysErr         chan<- error // Reports fatal error to core
}

// NewWriter creates and returns writer
func NewWriter(conn Connection, cfg *ethconn.Config, log log15.Logger, stop <-chan int, sysErr chan<- error) *writer {
	return &writer{
		cfg:     cfg,
		conn:    conn,
		log:     log,
		msgChan: make(chan msg.Message, msgLimit),
		stop:    stop,
		sysErr:  sysErr,
	}
}

func (w *writer) start() error {
	w.log.Debug("Starting ethereum writer...")
	go func() {
		for {
			select {
			case <-w.stop:
				close(w.msgChan)
				w.log.Info("writer stopped")
				return
			case msg := <-w.msgChan:
				result := w.processMessage(msg)
				w.log.Info("processMessage", "result", result)
				if !result {
					w.sysErr <- fmt.Errorf("processMessage failed")
				}
			}
		}
	}()

	return nil
}

// setContract adds the bound receiver bridgeContract to the writer
func (w *writer) setContract(bridge *Bridge.Bridge) {
	w.bridgeContract = bridge
}

func (w *writer) ResolveMessage(m msg.Message) bool {
	w.log.Info("Attempting to resolve message", "type", m.Type, "src", m.Source, "dst", m.Destination, "nonce", m.DepositNonce, "rId", m.ResourceId.Hex())
	w.log.Info("ResolveMessage: size of msgChan", "size", len(w.msgChan))
	w.msgChan <- m
	return true
}

// ResolveMessage handles any given message based on type
// A bool is returned to indicate failure/success, this should be ignored except for within tests.
func (w *writer) processMessage(m msg.Message) bool {
	w.log.Info("Attempting to process message", "type", m.Type, "src", m.Source, "dst", m.Destination, "nonce", m.DepositNonce, "rId", m.ResourceId.Hex())
	switch m.Type {
	case msg.FungibleTransfer:
		result := make(chan bool)
		defer close(result)
		go w.createErc20Proposal(m, result)
		select {
		case <-w.stop:
			return false
		case re := <-result:
			return re
		}
	default:
		w.log.Error("Unknown message type received", "type", m.Type)
		return false
	}
}
