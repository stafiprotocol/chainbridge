// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"fmt"
	"github.com/stafiprotocol/chainbridge/config"
	"sync"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge-utils/msg"
	gsrpc "github.com/stafiprotocol/go-substrate-rpc-client"
	"github.com/stafiprotocol/go-substrate-rpc-client/signature"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
)

type Connection struct {
	api         *gsrpc.SubstrateAPI
	log         log15.Logger
	url         string                 // API endpoint
	name        string                 // Chain name
	meta        types.Metadata         // Latest chain metadata
	metaLock    sync.RWMutex           // Lock metadata for updates, allows concurrent reads
	genesisHash types.Hash             // Chain genesis hash
	key         *signature.KeyringPair // Keyring used for signing
	nonce       types.U32              // Latest account nonce
	nonceLock   sync.Mutex             // Locks nonce for updates
	stop        <-chan int             // Signals system shutdown, should be observed in all selects and loops
	sysErr      chan<- error           // Propagates fatal errors to core
}

func NewConnection(url string, name string, key *signature.KeyringPair, log log15.Logger, stop <-chan int, sysErr chan<- error) *Connection {
	return &Connection{url: url, name: name, key: key, log: log, stop: stop, sysErr: sysErr}
}

func (c *Connection) getMetadata() (meta types.Metadata) {
	c.metaLock.RLock()
	meta = c.meta
	c.metaLock.RUnlock()
	return meta
}

func (c *Connection) updateMetatdata() error {
	c.metaLock.Lock()
	meta, err := c.api.RPC.State.GetMetadataLatest()
	if err != nil {
		c.metaLock.Unlock()
		return err
	}
	c.meta = *meta
	c.metaLock.Unlock()
	return nil
}

func (c *Connection) Connect() error {
	c.log.Info("Connecting to substrate chain...", "url", c.url)
	api, err := gsrpc.NewSubstrateAPI(c.url)
	if err != nil {
		return err
	}
	c.api = api

	// Fetch metadata
	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return err
	}
	c.meta = *meta


	// Fetch genesis hash
	genesisHash, err := c.api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return err
	}
	c.genesisHash = genesisHash
	c.log.Info("Fetched substrate genesis hash", "hash", genesisHash.Hex())
	return nil
}

// queryStorage performs a storage lookup. Arguments may be nil, result must be a pointer.
func (c *Connection) queryStorage(prefix, method string, arg1, arg2 []byte, result interface{}) (bool, error) {
	// Fetch account nonce
	data := c.getMetadata()
	key, err := types.CreateStorageKey(&data, prefix, method, arg1, arg2)
	if err != nil {
		return false, err
	}
	return c.api.RPC.State.GetStorageLatest(key, result)
}

// TODO: Add this to GSRPC
func getConst(meta *types.Metadata, prefix, name string, res interface{}) error {
	for _, mod := range meta.AsMetadataV12.Modules {
		if string(mod.Name) == prefix {
			for _, cons := range mod.Constants {
				if string(cons.Name) == name {
					return types.DecodeFromBytes(cons.Value, res)
				}
			}
		}
	}
	return fmt.Errorf("could not find constant %s.%s", prefix, name)
}

func (c *Connection) getConst(prefix, name string, res interface{}) error {
	meta := c.getMetadata()
	return getConst(&meta, prefix, name, res)
}

func (c *Connection) checkChainId(expected msg.ChainId) error {
	var actual msg.ChainId
	err := c.getConst(config.BridgeCommon, config.ChainIdentity, &actual)
	if err != nil {
		return err
	}

	if actual != expected {
		return fmt.Errorf("ChainID is incorrect, Expected chainId: %d, got chainId: %d", expected, actual)
	}

	return nil
}

func (c *Connection) getLatestNonce() (types.U32, error) {
	var acct types.AccountInfo
	exists, err := c.queryStorage("System", "Account", c.key.PublicKey, nil, &acct)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, nil
	}

	return acct.Nonce, nil
}
func (c *Connection) Close() {
}
