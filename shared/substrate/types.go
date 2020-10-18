// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package utils

import (
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
)

const BridgePalletName = "ChainBridge"
const BridgeStoragePrefix = "ChainBridge"

type Erc721Token struct {
	Id       types.U256
	Metadata types.Bytes
}
