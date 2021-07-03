// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package utils

import (
	"github.com/ethereum/go-ethereum/common"

	"github.com/stafiprotocol/chainbridge/utils/keystore"
)

var (
	RelayerAddresses = []common.Address{
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.AliceKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.BobKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.CharlieKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.DaveKey].Address()),
		common.HexToAddress(keystore.TestKeyRing.EthereumKeys[keystore.EveKey].Address()),
	}

	ZeroAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")
)

type DeployedContracts struct {
	BridgeAddress       common.Address
	ERC20HandlerAddress common.Address
}
