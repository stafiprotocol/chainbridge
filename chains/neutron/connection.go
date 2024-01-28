// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package neutron

import (
	"fmt"
	"os"

	"github.com/ChainSafe/log15"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types"
	neutronClient "github.com/stafihub/neutron-relay-sdk/client"
	commonLog "github.com/stafihub/neutron-relay-sdk/common/log"
	commonCore "github.com/stafihub/rtoken-relay-core/common/core"
	"github.com/stafiprotocol/chainbridge/utils/core"
)

type Connection struct {
	bridgeAddress string
	url           []string // API endpoint
	name          string   // Chain name
	client        *neutronClient.Client
	stop          <-chan int // Signals system shutdown, should be observed in all selects and loops
	log           log15.Logger
	from          string
}

func NewConnection(cfg *core.ChainConfig, log log15.Logger, stop <-chan int) (*Connection, error) {
	log.Info("NewConnection", "name", cfg.Name, "KeystorePath", cfg.KeystorePath, "Endpoint", cfg.EndpointList)
	fmt.Printf("Will open stafihub wallet from <%s>. \nPlease ", cfg.KeystorePath)
	key, err := keyring.New(types.KeyringServiceName(), keyring.BackendFile, cfg.KeystorePath, os.Stdin, neutronClient.MakeEncodingConfig().Marshaler)
	if err != nil {
		return nil, err
	}
	account := cfg.From

	gasPrice := cfg.Opts["gasPrice"]
	client, err := neutronClient.NewClient(key, account, gasPrice, "neutron", cfg.EndpointList, commonLog.NewLog("client", "chain bridge"))
	if err != nil {
		return nil, fmt.Errorf("hubClient.NewClient err: %s", err)
	}
	bridgeAddress := cfg.Opts["bridgeAddress"]

	done := commonCore.UseSdkConfigContext(client.GetAccountPrefix())
	_, err = types.AccAddressFromBech32(bridgeAddress)
	from := client.GetFromAddress().String()
	done()
	if err != nil {
		return nil, err
	}

	return &Connection{
		bridgeAddress: bridgeAddress,
		url:           cfg.EndpointList,
		name:          cfg.Name,
		client:        client,
		stop:          stop,
		log:           log,
		from:          from,
	}, nil
}

func (c *Connection) Address() string {
	return c.from
}

func (c *Connection) LatestBlockNumber() (uint64, error) {
	blockHeight, err := c.client.GetCurrentBlockHeight()
	return uint64(blockHeight), err
}

func (c *Connection) FinalizedBlockNumber() (uint64, error) {
	blockHeight, err := c.client.GetCurrentBlockHeight()
	return uint64(blockHeight - 6), err
}

func (c *Connection) Close() {
}
