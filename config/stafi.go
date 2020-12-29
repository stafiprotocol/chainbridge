package config

import "math/big"

const (
	BridgeCommon            = "BridgeCommon"
	BridgeSwap              = "BridgeSwap"
	ChainIdentity           = "ChainIdentity"
	FungibleTransferEventId = "FungibleTransfer"
	AcknowledgeProposal     = "BridgeCommon.acknowledge_proposal"
)

var (
	// DecimalFactor = ERC20Decimal / StafiDecimal
	DecimalFactor = big.NewInt(1000000)
)
