package substrate

import (
	"fmt"
	"math/big"

	"github.com/stafiprotocol/chainbridge/utils/msg"
	"github.com/stafiprotocol/go-substrate-rpc-client/scale"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
)

var decimalDefault = "Default"

type EventFungibleTransfer struct {
	Destination  uint8
	DepositNonce uint64
	ResourceId   msg.ResourceId
	Amount       *big.Int
	Recipient    []byte
	Decimal      *big.Int
}

type ChainIdParam struct {
	Type  string `json:"type"`
	Value uint8  `json:"value"`
}

type DepositNonceParam struct {
	Type  string `json:"type"`
	Value uint64 `json:"value"`
}

type voteState struct {
	Voted  []types.AccountID
	Status VoteStatus
}

type VoteStatus string

const (
	VoteStatusActive   = "Active"
	VoteStatusPassed   = "Passed"
	VoteStatusExpired  = "Expired"
	VoteStatusExecuted = "Executed"
)

func (m *VoteStatus) Decode(decoder scale.Decoder) error {
	b, err := decoder.ReadOneByte()
	if err != nil {
		return err
	}

	switch b {
	case 0:
		*m = VoteStatusActive
	case 1:
		*m = VoteStatusPassed
	case 2:
		*m = VoteStatusExpired
	case 3:
		*m = VoteStatusExecuted
	default:
		return fmt.Errorf("VoteStatus decode error: %d", b)
	}

	return nil
}

// proposal represents an on-chain proposal
type proposal struct {
	DepositNonce types.U64     `json:"depositNonce"`
	Call         types.Call    `json:"call"`
	SourceId     types.U8      `json:"sourceId"`
	ResourceId   types.Bytes32 `json:"resourceId"`
	Method       string        `json:"method"`
}

// encode takes only nonce and call and encodes them for storage queries
func (p *proposal) encode() ([]byte, error) {
	return types.EncodeToBytes(struct {
		types.U64
		types.Call
	}{p.DepositNonce, p.Call})
}
