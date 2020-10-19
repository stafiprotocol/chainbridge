package config

import "math/big"

const (
	BridgeCommon = "BridgeCommon"
	BridgeSwap = "BridgeSwap"

	ChainIdentity = "ChainIdentity"
	NativeTokenId = "NativeTokenId"
	ChainNonces = "ChainNonces"
)

var (
	// DecimalFactor = ERC20Decimal / StafiDecimal
	DecimalFactor = big.NewInt(1000000)
)
