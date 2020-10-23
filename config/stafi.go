package config

import "math/big"

const (
	BridgeCommon = "BridgeCommon"
	ChainIdentity = "ChainIdentity"
)

var (
	// DecimalFactor = ERC20Decimal / StafiDecimal
	DecimalFactor = big.NewInt(1000000)
)
