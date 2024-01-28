// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

/*
The substrate package contains the logic for interacting with substrate chains.
The current supported transfer types are Fungible.

There are 3 major components: the connection, the listener, and the writer.

# Connection

The Connection handles connecting to the substrate client, and submitting transactions to the client.
It also handles state queries. The connection is shared by the writer and listener.

# Listener

The substrate listener polls blocks and parses the associated events for the three transfer types. It then forwards these into the router.

# Writer

As the writer receives messages from the router, nothing happened.
*/
package neutron

import (
	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge/utils/core"
	"github.com/stafiprotocol/chainbridge/utils/msg"
)

type Chain struct {
	cfg    *core.ChainConfig // The config of the chain
	conn   *Connection       // THe chains connection
	writer *writer           // The writer of the chain
	stop   chan<- int
}

func InitializeChain(cfg *core.ChainConfig, logger log15.Logger, sysErr chan<- error) (*Chain, error) {
	stop := make(chan int)
	conn, err := NewConnection(cfg, logger, stop)
	if err != nil {
		return nil, err
	}
	w := NewWriter(conn, logger, sysErr, stop)
	return &Chain{cfg: cfg, conn: conn, writer: w, stop: stop}, nil
}

func (c *Chain) Start() error {
	err := c.writer.start()
	if err != nil {
		return err
	}

	c.conn.log.Debug("Successfully started chain", "chainId", c.cfg.Id)
	return nil
}

func (c *Chain) SetRouter(r *core.Router) {
	r.Listen(c.cfg.Id, c.writer)
}

func (c *Chain) Id() msg.ChainId {
	return c.cfg.Id
}

func (c *Chain) Name() string {
	return c.cfg.Name
}

func (c *Chain) Stop() {
	close(c.stop)
}
