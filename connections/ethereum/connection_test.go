// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package ethereum_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ChainSafe/log15"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethutils "github.com/stafiprotocol/chainbridge/shared/ethereum"
	ethtest "github.com/stafiprotocol/chainbridge/shared/ethereum/testing"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/chainbridge/connections/ethereum"
)

var TestEndpoint = "https://mainnet.infura.io/v3/4cb873af07a84e42a952189eff3a6954"
var AliceKp = keystore.TestKeyRing.EthereumKeys[keystore.AliceKey]
var GasLimit = big.NewInt(ethutils.DefaultGasLimit)
var MaxGasPrice = big.NewInt(ethutils.DefaultMaxGasPrice)
var EtherScanUrl = "https://api-cn.etherscan.com/api?module=gastracker&action=gasoracle&apikey=RFPRRAX9BZGX2SHNNHXIRVPCSDPZUUGDFN"
func TestConnect(t *testing.T) {
	conn := ethereum.NewConnection(TestEndpoint, true, AliceKp, log15.Root(), GasLimit, MaxGasPrice,EtherScanUrl)
	err := conn.Connect()
	if err != nil {
		t.Fatal(err)
	}
	conn.Close()
}

// TestContractCode is used to make sure the contracts are deployed correctly.
// This is probably the least intrusive way to check if the contracts exists
func TestContractCode(t *testing.T) {
	client := ethtest.NewClient(t, TestEndpoint, AliceKp)
	contracts, err := ethutils.DeployContracts(client, 0, big.NewInt(0))
	if err != nil {
		t.Fatal(err)
	}

	conn := ethereum.NewConnection(TestEndpoint, false, AliceKp, log15.Root(), GasLimit, MaxGasPrice,EtherScanUrl)
	err = conn.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// The following section checks if the byteCode exists on the chain at the specificed Addresses
	err = conn.EnsureHasBytecode(contracts.BridgeAddress)
	if err != nil {
		t.Fatal(err)
	}

	err = conn.EnsureHasBytecode(ethcmn.HexToAddress("0x0"))
	if err == nil {
		t.Fatal("should detect no bytecode")
	}

}

func TestConnection_SafeEstimateGas(t *testing.T) {
	// MaxGasPrice is the constant price on the dev network, so we increase it here by 1 to ensure it adjusts
	conn := ethereum.NewConnection(TestEndpoint, false, AliceKp, log15.Root(), GasLimit, 
	MaxGasPrice.Add(MaxGasPrice, big.NewInt(1)),EtherScanUrl)
	err := conn.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	price, err := conn.SafeEstimateGas(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if price.Cmp(MaxGasPrice) == 0 {
		t.Fatalf("Gas price should be less than max. Suggested: %s Max: %s", price.String(), MaxGasPrice.String())
	}
}

func TestConnection_SafeEstimateGasMax(t *testing.T) {
	maxPrice := big.NewInt(300e9)
	conn := ethereum.NewConnection(TestEndpoint, true, AliceKp, log15.Root(), GasLimit, maxPrice,EtherScanUrl)
	err := conn.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	price, err := conn.SafeEstimateGas(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(price.String())

	if price.Cmp(maxPrice) != 0 {
		t.Fatalf("Gas price should equal max. Suggested: %s Max: %s", price.String(), maxPrice.String())
	}
}
