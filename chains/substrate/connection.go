// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"fmt"
	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge/config"
	"github.com/stafiprotocol/chainbridge/utils/msg"
	gsrpc "github.com/stafiprotocol/go-substrate-rpc-client"
)

type Connection struct {
	api    *gsrpc.SubstrateAPI
	log    log15.Logger
	url    string            // API endpoint
	name   string            // Chain name
	opts   map[string]string // opts
	stop   <-chan int        // Signals system shutdown, should be observed in all selects and loops
	sysErr chan<- error      // Propagates fatal errors to core

}

func NewConnection(url string, name string, opts map[string]string, log log15.Logger, stop <-chan int, sysErr chan<- error) *Connection {
	return &Connection{url: url, name: name, opts: opts, log: log, stop: stop, sysErr: sysErr}
}

func (c *Connection) Connect() error {
	c.log.Info("Connecting to substrate chain...", "url", c.url)
	api, err := gsrpc.NewSubstrateAPI(c.url)
	if err != nil {
		return err
	}
	c.api = api
	return nil
}

func (c *Connection) checkChainId(expected msg.ChainId) error {
	var actual msg.ChainId
	err := c.api.RPC.State.GetConst(config.BridgeCommon, config.ChainIdentity, &actual)
	if err != nil {
		return err
	}

	if actual != expected {
		return fmt.Errorf("ChainID is incorrect, Expected chainId: %d, got chainId: %d", expected, actual)
	}
	return nil
}

func (c *Connection) Close() {
}
