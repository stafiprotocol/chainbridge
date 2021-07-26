package substrate

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stafiprotocol/chainbridge/shared/substrate"
	"math/big"

	"github.com/ChainSafe/log15"
	"github.com/ethereum/go-ethereum/common/hexutil"
	scalecodec "github.com/itering/scale.go"
	"github.com/itering/scale.go/utiles"
	"github.com/stafiprotocol/chainbridge/config"
	"github.com/stafiprotocol/chainbridge/utils/msg"
)

var (
	ValueNotStringError = errors.New("value not string")
)

type eventName string
type eventHandler func(interface{}, log15.Logger) (msg.Message, error)

const FungibleTransfer eventName = config.FungibleTransferEventId

var Subscriptions = []struct {
	name    eventName
	handler eventHandler
}{{FungibleTransfer, fungibleTransferHandler}}

func fungibleTransferHandler(evtI interface{}, log log15.Logger) (msg.Message, error) {
	evt, ok := evtI.(*EventFungibleTransfer)
	if !ok {
		return msg.Message{}, fmt.Errorf("failed to cast EventFungibleTransfer type")
	}

	resourceId := msg.ResourceId(evt.ResourceId)
	log.Info("Got fungible transfer event!", "destination", evt.Destination, "resourceId", resourceId.Hex(), "amount", evt.Amount, "Recipient", hexutil.Encode(evt.Recipient))

	return msg.NewFungibleTransfer(
		0,
		msg.ChainId(evt.Destination),
		msg.Nonce(evt.DepositNonce),
		evt.Amount.Mul(evt.Amount, evt.Decimal),
		resourceId,
		evt.Recipient,
	), nil
}

// FungibleTransfer(AccountId, ChainId, DepositNonce, ResourceId, U256, Vec<u8>)
func FungibleTransferEventData(evt *substrate.ChainEvent, decimals map[string]*big.Int) (*EventFungibleTransfer, error) {
	if len(evt.Params) != 6 {
		return nil, fmt.Errorf("EventFungibleTransfer params number not right: %d, expected: 6", len(evt.Params))
	}

	chainId, err := parseChainId(evt.Params[1])
	if err != nil {
		return nil, fmt.Errorf("EventFungibleTransfer params[1] -> chainId error: %s", err)
	}

	nonce, err := parseDepositNonce(evt.Params[2])
	if err != nil {
		return nil, fmt.Errorf("EventFungibleTransfer params[1] -> chainId error: %s", err)
	}

	resourceId, err := parseBytes(evt.Params[3].Value)
	if err != nil {
		return nil, fmt.Errorf("EventFungibleTransfer params[3] -> resourceId error: %s", err)
	}

	amount, err := parseU256(evt.Params[4].Value)
	if err != nil {
		return nil, fmt.Errorf("EventFungibleTransfer params[4] -> amount error: %s", err)
	}

	rec, err := parseBytes(evt.Params[5].Value)
	if err != nil {
		return nil, fmt.Errorf("EventFungibleTransfer params[5] -> recipient error: %s", err)
	}

	eft := &EventFungibleTransfer{
		Destination:  chainId,
		DepositNonce: nonce,
		ResourceId:   msg.ResourceIdFromSlice(resourceId),
		Amount:       amount,
		Recipient:    rec,
	}

	decimal, ok := decimals[eft.ResourceId.Hex()]
	if !ok {
		decimal, ok = decimals["Default"]
		if !ok {
			return nil, fmt.Errorf("failed to get decimal")
		}
	}
	eft.Decimal = decimal
	return eft, nil
}

func parseChainId(param scalecodec.EventParam) (uint8, error) {
	var cp ChainIdParam
	x, err := json.Marshal(param)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(x, &cp)
	if err != nil {
		return 0, err
	}
	return cp.Value, nil
}

func parseDepositNonce(param scalecodec.EventParam) (uint64, error) {
	var cp DepositNonceParam
	x, err := json.Marshal(param)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(x, &cp)
	if err != nil {
		return 0, err
	}
	return cp.Value, nil
}

func parseBytes(value interface{}) ([]byte, error) {
	val, ok := value.(string)
	if !ok {
		return nil, ValueNotStringError
	}

	bz, err := hexutil.Decode(utiles.AddHex(val))
	if err != nil {
		return nil, err
	}

	return bz, nil
}

func parseU256(value interface{}) (*big.Int, error) {
	val, ok := value.(string)
	if !ok {
		return nil, ValueNotStringError
	}

	amount := new(big.Int)
	b := utiles.HexToBytes(val)
	rb := utiles.ReverseBytes(b)
	amount.SetBytes(rb)

	return amount, nil
}
