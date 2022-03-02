// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package stafihub

import (
	"fmt"
	"os"

	"github.com/ChainSafe/log15"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stafiprotocol/chainbridge/shared/stafihub"
	"github.com/stafiprotocol/chainbridge/utils/core"
)

type Connection struct {
	url    string // API endpoint
	name   string // Chain name
	client *stafihub.Client
	stop   <-chan int // Signals system shutdown, should be observed in all selects and loops
	log    log15.Logger
}

func NewConnection(cfg *core.ChainConfig, log log15.Logger, stop <-chan int) (*Connection, error) {
	log.Info("NewConnection", "name", cfg.Name, "KeystorePath", cfg.KeystorePath, "Endpoint", cfg.Endpoint)
	fmt.Printf("Will open stafihub wallet from <%s>. \nPlease ", cfg.KeystorePath)
	key, err := keyring.New(types.KeyringServiceName(), keyring.BackendFile, cfg.KeystorePath, os.Stdin)
	if err != nil {
		return nil, err
	}
	account := cfg.From
	gasPrice := cfg.Opts["gasPrice"]
	client, err := stafihub.NewClient(key, account, gasPrice, cfg.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("hubClient.NewClient err: %s", err)
	}

	return &Connection{
		url:    cfg.Endpoint,
		name:   cfg.Name,
		client: client,
		stop:   stop,
		log:    log,
	}, nil
}

func (c *Connection) Address() string {
	return c.client.GetFromAddress().String()
}

func (c *Connection) LatestBlockNumber() (uint64, error) {
	blockHeight, err := c.client.GetCurrentBlockHeight()
	return uint64(blockHeight), err
}

func (c *Connection) FinalizedBlockNumber() (uint64, error) {
	blockHeight, err := c.client.GetCurrentBlockHeight()
	return uint64(blockHeight - 6), err
}

func (c *Connection) GetEvents(blockNum uint64) ([]*types.TxResponse, error) {
	return c.client.GetBlockTxs(int64(blockNum))
}

func (c *Connection) Close() {
}
