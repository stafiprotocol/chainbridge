package substrate

import (
	"errors"
	scalecodec "github.com/itering/scale.go"
)

const (
	ChainTypeStafi    = "stafi"
	ChainTypePolkadot = "polkadot"

	AddressTypeAccountId    = "AccountId"
	AddressTypeMultiAddress = "MultiAddress"
)

var (
	TerminatedError           = errors.New("terminated")
	BondEqualToUnbondError    = errors.New("BondEqualToUnbondError")
	BondSmallerThanLeastError = errors.New("BondSmallerThanLeastError")
)

type ChainEvent struct {
	ModuleId string                  `json:"module_id" `
	EventId  string                  `json:"event_id" `
	Params   []scalecodec.EventParam `json:"params"`
}
