// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"fmt"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge/shared/substrate"
	"github.com/stafiprotocol/chainbridge/utils/core"
	"github.com/stafiprotocol/chainbridge/utils/crypto/sr25519"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/chainbridge/utils/msg"
	"github.com/stafiprotocol/go-substrate-rpc-client/signature"
)

type Connection struct {
	url  string // API endpoint
	name string // Chain name
	sc   *substrate.SarpcClient
	gc   *substrate.GsrpcClient
	stop <-chan int // Signals system shutdown, should be observed in all selects and loops
	log  log15.Logger
	key  *signature.KeyringPair
}

const (
	DefaultTypeFilePath = "../../network/stafi.json"
)

func NewConnection(cfg *core.ChainConfig, log log15.Logger, stop <-chan int) (*Connection, error) {
	log.Info("NewConnection", "name", cfg.Name, "KeystorePath", cfg.KeystorePath, "Endpoint", cfg.Endpoint)

	kp, err := keystore.KeypairFromAddress(cfg.From, keystore.SubChain, cfg.KeystorePath, cfg.Insecure)
	if err != nil {
		return nil, err
	}
	krp := kp.(*sr25519.Keypair).AsKeyringPair()
	gc, err := substrate.NewGsrpcClient(cfg.Endpoint, "AccountId", krp, log, stop)
	if err != nil {
		return nil, err
	}

	path := DefaultTypeFilePath
	if file, ok := cfg.Opts["typeRegister"]; ok {
		path = file
	}
	sc, err := substrate.NewSarpcClient(cfg.Endpoint, path, log)
	if err != nil {
		return nil, err
	}

	return &Connection{
		url:  cfg.Endpoint,
		name: cfg.Name,
		sc:   sc,
		gc:   gc,
		stop: stop,
		log:  log,
		key:  krp,
	}, nil
}

func (c *Connection) Address() string {
	return c.gc.Address()
}

func (c *Connection) LatestBlockNumber() (uint64, error) {
	return c.gc.GetLatestBlockNumber()
}

func (c *Connection) FinalizedBlockNumber() (uint64, error) {
	return c.gc.GetFinalizedBlockNumber()
}

func (c *Connection) GetEvents(blockNum uint64) ([]*substrate.ChainEvent, error) {
	return c.sc.GetEvents(blockNum)
}

// queryStorage performs a storage lookup. Arguments may be nil, result must be a pointer.
func (c *Connection) QueryStorage(prefix, method string, arg1, arg2 []byte, result interface{}) (bool, error) {
	return c.gc.QueryStorage(prefix, method, arg1, arg2, result)
}

func (c *Connection) checkChainId(expected msg.ChainId) error {
	actual, err := c.gc.ChainId()
	if err != nil {
		return err
	}

	if actual != uint8(expected) {
		return fmt.Errorf("ChainID is incorrect, Expected chainId: %d, got chainId: %d", expected, actual)
	}
	return nil
}

func (c *Connection) Close() {
}
