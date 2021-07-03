// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package ethereum

import (
	"math/big"
	"testing"
	"time"

	"github.com/ChainSafe/log15"
	"github.com/ethereum/go-ethereum/common"
	connection "github.com/stafiprotocol/chainbridge/connections/ethereum"
	utils "github.com/stafiprotocol/chainbridge/shared/ethereum"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/chainbridge/utils/msg"
)

const TestEndpoint = "ws://localhost:8545"

var TestLogger = newTestLogger("test")
var TestTimeout = time.Second * 30

var AliceKp = keystore.TestKeyRing.EthereumKeys[keystore.AliceKey]
var BobKp = keystore.TestKeyRing.EthereumKeys[keystore.BobKey]

var TestRelayerThreshold = big.NewInt(2)
var TestChainId = msg.ChainId(0)

var aliceTestConfig = createConfig("alice", nil, nil)

func createConfig(name string, startBlock *big.Int, contracts *utils.DeployedContracts) *Config {
	cfg := &Config{
		name:                 name,
		id:                   0,
		endpoint:             TestEndpoint,
		from:                 name,
		keystorePath:         "",
		blockstorePath:       "",
		freshStart:           true,
		bridgeContract:       common.Address{},
		erc20HandlerContract: common.Address{},
		gasLimit:             big.NewInt(DefaultGasLimit),
		maxGasPrice:          big.NewInt(DefaultGasPrice),
		http:                 false,
		startBlock:           startBlock,
	}

	if contracts != nil {
		cfg.bridgeContract = contracts.BridgeAddress
		cfg.erc20HandlerContract = contracts.ERC20HandlerAddress
	}

	return cfg
}

func newTestLogger(name string) log15.Logger {
	tLog := log15.New("chain", name)
	tLog.SetHandler(log15.LvlFilterHandler(log15.LvlError, tLog.GetHandler()))
	return tLog
}

func newLocalConnection(t *testing.T, cfg *Config) *connection.Connection {
	kp := keystore.TestKeyRing.EthereumKeys[cfg.from]
	conn := connection.NewConnection(TestEndpoint, false, kp, TestLogger, big.NewInt(DefaultGasLimit), big.NewInt(DefaultGasPrice))
	err := conn.Connect()
	if err != nil {
		t.Fatal(err)
	}

	return conn
}
