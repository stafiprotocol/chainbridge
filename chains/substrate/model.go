package substrate

import (
	scalecodec "github.com/itering/scale.go"
	"github.com/stafiprotocol/chainbridge/utils/msg"
	"math/big"
)

type ChainEvent struct {
	ModuleId string                  `json:"module_id" `
	EventId  string                  `json:"event_id" `
	Params   []scalecodec.EventParam `json:"params"`
}

type EventFungibleTransfer struct {
	Destination  uint8
	DepositNonce uint64
	ResourceId   msg.ResourceId
	Amount       *big.Int
	Recipient    []byte
}

type ChainIdParam struct {
	Type  string `json:"type"`
	Value uint8  `json:"value"`
}

type DepositNonceParam struct {
	Type  string `json:"type"`
	Value uint64 `json:"value"`
}
