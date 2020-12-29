// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"errors"
	"fmt"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge/config"
	"github.com/stafiprotocol/chainbridge/utils/msg"
	gsrpc "github.com/stafiprotocol/go-substrate-rpc-client"
	"github.com/stafiprotocol/go-substrate-rpc-client/rpc/author"
	"github.com/stafiprotocol/go-substrate-rpc-client/signature"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
)

type Connection struct {
	api         *gsrpc.SubstrateAPI
	log         log15.Logger
	url         string                 // API endpoint
	name        string                 // Chain name
	opts        map[string]string      // opts
	stop        <-chan int             // Signals system shutdown, should be observed in all selects and loops
	sysErr      chan<- error           // Propagates fatal errors to core
	key         *signature.KeyringPair // Keyring used for signing
	genesisHash types.Hash             // Chain genesis hash
}

func NewConnection(url string, name string, key *signature.KeyringPair, opts map[string]string, log log15.Logger, stop <-chan int, sysErr chan<- error) *Connection {
	return &Connection{url: url, name: name, key: key, opts: opts, log: log, stop: stop, sysErr: sysErr}
}

func (c *Connection) Connect() error {
	c.log.Info("Connecting to substrate chain...", "url", c.url)
	api, err := gsrpc.NewSubstrateAPI(c.url)
	if err != nil {
		return err
	}

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return err
	}

	c.api = api
	c.genesisHash = genesisHash
	return nil
}

func (c *Connection) SignAndSubmitTx(ext *types.Extrinsic) error {
	err := c.signExtrinsic(ext)
	if err != nil {
		return err
	}

	// Do the transfer and track the actual status
	sub, err := c.api.RPC.Author.SubmitAndWatchExtrinsic(*ext)
	if err != nil {
		return err
	}
	c.log.Trace("Extrinsic submission succeeded")
	defer sub.Unsubscribe()

	return c.watchSubmission(sub)
}

func (c *Connection) watchSubmission(sub *author.ExtrinsicStatusSubscription) error {
	for {
		select {
		case <-c.stop:
			return TerminatedError
		case status := <-sub.Chan():
			switch {
			case status.IsInBlock:
				c.log.Info("Extrinsic included in block", "block", status.AsInBlock.Hex())
				return nil
			case status.IsRetracted:
				return fmt.Errorf("extrinsic retracted: %s", status.AsRetracted.Hex())
			case status.IsDropped:
				return fmt.Errorf("extrinsic dropped from network")
			case status.IsInvalid:
				return fmt.Errorf("extrinsic invalid")
			}
		case err := <-sub.Err():
			c.log.Trace("Extrinsic subscription error", "err", err)
			return err
		}
	}
}

func (c *Connection) NewUnsignedExtrinsic(callMethod string, args ...interface{}) (*types.Extrinsic, error) {
	c.log.Debug("Submitting substrate call...", "callMethod", callMethod, "sender", c.key.Address)
	meta, err := c.GetLatestMetadata()
	if err != nil {
		return nil, err
	}

	call, err := types.NewCall(meta, callMethod, args...)
	if err != nil {
		return nil, err
	}

	unsignedExt := types.NewExtrinsic(call)
	return &unsignedExt, nil
}

func (c *Connection) signExtrinsic(ext *types.Extrinsic) error {
	rv, err := c.GetLatestRuntimeVersion()
	if err != nil {
		return err
	}

	nonce, err := c.GetLatestNonce()
	if err != nil {
		return err
	}

	o := types.SignatureOptions{
		BlockHash:          c.genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        c.genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	err = ext.Sign(*c.key, o)
	if err != nil {
		return err
	}

	return nil
}

func (c *Connection) GetLatestMetadata() (*types.Metadata, error) {
	meta, err := c.api.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, err
	}

	return meta, nil
}

func (c *Connection) GetLatestRuntimeVersion() (*types.RuntimeVersion, error) {
	rv, err := c.api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return nil, err
	}

	return rv, nil
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

func (c *Connection) GetLatestNonce() (types.U32, error) {
	ac, err := c.GetAccountInfo()
	if err != nil {
		return 0, err
	}

	return ac.Nonce, nil
}

func (c *Connection) GetAccountInfo() (*types.AccountInfo, error) {
	ac := new(types.AccountInfo)
	exist, err := c.queryStorage("System", "Account", c.key.PublicKey, nil, &ac)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, errors.New("account not exist")
	}

	return ac, nil
}

// queryStorage performs a storage lookup. Arguments may be nil, result must be a pointer.
func (c *Connection) queryStorage(prefix, method string, arg1, arg2 []byte, result interface{}) (bool, error) {
	meta, err := c.GetLatestMetadata()
	if err != nil {
		return false, err
	}

	key, err := types.CreateStorageKey(meta, prefix, method, arg1, arg2)
	if err != nil {
		return false, err
	}

	ok, err := c.api.RPC.State.GetStorageLatest(key, result)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (c *Connection) Close() {
}
