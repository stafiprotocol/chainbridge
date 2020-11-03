package config

import "math/big"

const (
	BridgeCommon            = "BridgeCommon"
	ChainIdentity           = "ChainIdentity"
	FungibleTransferEventId = "FungibleTransfer"
)

var (
	// DecimalFactor = ERC20Decimal / StafiDecimal
	DecimalFactor = big.NewInt(1000000)
)
