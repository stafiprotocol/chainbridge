package events

import (
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
)

type SubEvents struct {
	types.EventRecords
	Events
}

type Events struct {
	BridgeCommon_FungibleTransfer []EventFungibleTransfer //nolint:stylecheck,golint
	BridgeCommon_ChainWhitelisted []EventChainWhitelisted //nolint:stylecheck,golint
}

type EventFungibleTransfer struct {
	Phase        types.Phase
	AccountID    types.AccountID
	Destination  types.U8
	DepositNonce types.U64
	ResourceId   types.Bytes32
	Amount       types.U256
	Recipient    types.Bytes
	Topics       []types.Hash
}

type EventChainWhitelisted struct {
	Phase   types.Phase
	ChainId types.U8
	Topics  []types.Hash
}
