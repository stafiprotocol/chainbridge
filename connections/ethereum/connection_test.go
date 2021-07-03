// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package ethereum_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ChainSafe/log15"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/stafiprotocol/chainbridge/bindings/Bridge"
	"github.com/stafiprotocol/chainbridge/connections/ethereum"
	ethutils "github.com/stafiprotocol/chainbridge/shared/ethereum"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
)

var (
	TestEndpoint    = "https://mainnet.infura.io/v3/4cb873af07a84e42a952189eff3a6954"
	BscTestEndpoint = "https://data-seed-prebsc-2-s3.binance.org:8545/"
	AliceKp         = keystore.TestKeyRing.EthereumKeys[keystore.AliceKey]
	GasLimit        = big.NewInt(ethutils.DefaultGasLimit)
	MaxGasPrice     = big.NewInt(ethutils.DefaultMaxGasPrice)
)

func TestConnect(t *testing.T) {
	conn := ethereum.NewConnection(TestEndpoint, true, AliceKp, log15.Root(), GasLimit, MaxGasPrice)
	err := conn.Connect()
	if err != nil {
		t.Fatal(err)
	}
	conn.Close()
}

func TestBscConnect(t *testing.T) {
	conn := ethereum.NewConnection(BscTestEndpoint, true, AliceKp, log15.Root(), GasLimit, MaxGasPrice)
	err := conn.Connect()
	if err != nil {
		t.Fatal(err)
	}
	conn.Close()

	bridgeContract := ethcmn.HexToAddress("0xb5Dc44e4e680d9a485fCbe1f465eA0F65287FBc0")
	bridge, err := Bridge.NewBridge(bridgeContract, conn.Client())
	if err != nil {
		t.Fatal(err)
	}

	id, err := bridge.ChainID(nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id)

	threshold, err := bridge.RelayerThreshold(nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(threshold)

	relayer := ethcmn.HexToAddress("0xBca9567A9e8D5F6F58C419d32aF6190F74C880e6")
	flag, err := bridge.IsRelayer(nil, relayer)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(flag)

	relayer1 := ethcmn.HexToAddress("0xaD0bf51f7fc89e262edBbdF53C260088B024D857")
	flag1, err := bridge.IsRelayer(nil, relayer1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(flag1)

	relayer2 := ethcmn.HexToAddress("0xBd39f5936969828eD9315220659cD11129071814")
	flag2, err := bridge.IsRelayer(nil, relayer2)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(flag2)
}

func TestConnection_SafeEstimateGas(t *testing.T) {
	// MaxGasPrice is the constant price on the dev network, so we increase it here by 1 to ensure it adjusts
	conn := ethereum.NewConnection(TestEndpoint, false, AliceKp, log15.Root(), GasLimit,
		MaxGasPrice.Add(MaxGasPrice, big.NewInt(1)))
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
	conn := ethereum.NewConnection(TestEndpoint, true, AliceKp, log15.Root(), GasLimit, maxPrice)
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
