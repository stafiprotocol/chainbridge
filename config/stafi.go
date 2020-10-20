package config

import "math/big"

const (
	BridgeCommon = "BridgeCommon"
	BridgeSwap = "BridgeSwap"

	ChainIdentity = "ChainIdentity"
	NativeTokenId = "NativeTokenId"
)



var (
	// DecimalFactor = ERC20Decimal / StafiDecimal
	DecimalFactor = big.NewInt(1000000)

	EventsToWatch = []EventFullName {
		{
			ModuleName: BridgeCommon,
			EventName:  "FungibleTransfer",
		},
		{
			ModuleName: BridgeCommon,
			EventName:  "ChainWhitelisted",
		},
		{
			ModuleName: "System",
			EventName:  "CodeUpdated",
		},
	}

)

type EventFullName struct {
	ModuleName string
	EventName string
}
