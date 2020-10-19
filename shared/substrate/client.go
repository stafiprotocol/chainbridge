// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package utils

import (
	gsrpc "github.com/stafiprotocol/go-substrate-rpc-client"
	"github.com/stafiprotocol/go-substrate-rpc-client/signature"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
)

// Client is a container for all the components required to submit extrinsics
// TODO: Perhaps this would benefit an interface so we can interchange Connection and a client like this
type Client struct {
	Api     *gsrpc.SubstrateAPI
	Meta    *types.Metadata
	Genesis types.Hash
	Key     *signature.KeyringPair
}

func CreateClient(key *signature.KeyringPair, endpoint string) (*Client, error) {
	c := &Client{Key: key}
	api, err := gsrpc.NewSubstrateAPI(endpoint)
	if err != nil {
		return nil, err
	}
	c.Api = api

	// Fetch metadata
	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, err
	}
	c.Meta = meta

	// Fetch genesis hash
	genesisHash, err := c.Api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return nil, err
	}
	c.Genesis = genesisHash

	return c, nil
}

func (c *Client) LatestBlock() (uint64, error) {
	head, err := c.Api.RPC.Chain.GetHeaderLatest()
	if err != nil {
		return 0, err
	}
	return uint64(head.Number), nil
}

func (c *Client) GetConst(prefix, name string, res interface{}) error {
	return c.Api.RPC.State.GetConstWithMetadata(c.Meta, prefix, name, res)
}
