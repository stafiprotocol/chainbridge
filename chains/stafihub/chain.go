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
package stafihub

import (
	"strconv"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge/utils/blockstore"
	"github.com/stafiprotocol/chainbridge/utils/core"
	"github.com/stafiprotocol/chainbridge/utils/msg"
)

type Chain struct {
	cfg      *core.ChainConfig // The config of the chain
	conn     *Connection       // THe chains connection
	listener *listener         // The listener of this chain
	writer   *writer           // The writer of the chain
	stop     chan<- int
}

func InitializeChain(cfg *core.ChainConfig, logger log15.Logger, sysErr chan<- error) (*Chain, error) {
	stop := make(chan int)
	conn, err := NewConnection(cfg, logger, stop)
	if err != nil {
		return nil, err
	}

	// Attempt to load latest block
	bs, err := blockstore.NewBlockstore(cfg.BlockstorePath, cfg.Id, conn.Address())
	if err != nil {
		return nil, err
	}
	startBlock := parseStartBlock(cfg)
	if !cfg.FreshStart {
		startBlock, err = checkBlockstore(bs, startBlock)
		if err != nil {
			return nil, err
		}
	}

	if cfg.LatestBlock {
		curr, err := conn.LatestBlockNumber()
		if err != nil {
			return nil, err
		}
		startBlock = curr
	}

	// Setup listener & writer
	l := NewListener(conn, cfg.Name, cfg.Id, startBlock, logger, bs, stop, sysErr)
	w := NewWriter(conn, logger, sysErr, stop)
	return &Chain{cfg: cfg, conn: conn, listener: l, writer: w, stop: stop}, nil
}

func (c *Chain) Start() error {
	err := c.listener.start()
	if err != nil {
		return err
	}

	if len(c.conn.client.GetFromName()) > 0 {
		err = c.writer.start()
		if err != nil {
			return err
		}
	}

	c.conn.log.Debug("Successfully started chain", "chainId", c.cfg.Id)
	return nil
}

func (c *Chain) SetRouter(r *core.Router) {
	r.Listen(c.cfg.Id, c.writer)
	c.listener.setRouter(r)
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

// checkBlockstore queries the blockstore for the latest known block. If the latest block is
// greater than startBlock, then the latest block is returned, otherwise startBlock is.
func checkBlockstore(bs *blockstore.Blockstore, startBlock uint64) (uint64, error) {
	latestBlock, err := bs.TryLoadLatestBlock()
	if err != nil {
		return 0, err
	}

	if latestBlock.Uint64() > startBlock {
		return latestBlock.Uint64(), nil
	} else {
		return startBlock, nil
	}
}

func parseStartBlock(cfg *core.ChainConfig) uint64 {
	if blk, ok := cfg.Opts["startBlock"]; ok {
		res, err := strconv.ParseUint(blk, 10, 32)
		if err != nil {
			panic(err)
		}
		return res
	}
	return 0
}
