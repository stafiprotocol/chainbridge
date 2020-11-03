package substrate

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ChainSafe/log15"
	"github.com/ethereum/go-ethereum/common/hexutil"
	scalecodec "github.com/itering/scale.go"
	"github.com/itering/scale.go/types"
	"github.com/itering/scale.go/utiles"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/util"
	"github.com/itering/substrate-api-rpc/websocket"
	"github.com/stafiprotocol/chainbridge-utils/msg"
	"github.com/stafiprotocol/chainbridge/config"
)

const (
	wsId                = 1
	storageKey          = "0x26aa394eea5630e07c48ae0c9558cef780d41e5e16056765bc8461851072c9d7"
	DefaultTypeFilePath = "./network/stafi.json"
)

func (l *listener) GetEventsAt(blockNum uint64) ([]*EventFungibleTransfer, error) {
	//l.log.Info("GetEventsAt", "CurrentBlockNum", blockNum)
	evts := make([]*EventFungibleTransfer, 0)
	v := &rpc.JsonRpcResult{}
	// Block Hash
	err := websocket.SendWsRequest(nil, v, rpc.ChainGetBlockHash(wsId, int(blockNum)))
	if err != nil {
		return nil, err
	}
	blockHash, err := v.ToString()
	if err != nil {
		return nil, err
	}

	// event
	err = websocket.SendWsRequest(nil, v, rpc.StateGetStorage(wsId, storageKey, blockHash))
	eventRaw, err := v.ToString()
	if err != nil {
		return nil, err
	}

	// metadata
	err = websocket.SendWsRequest(nil, v, rpc.StateGetMetadata(wsId, blockHash))
	metaRaw, err := v.ToString()
	if err != nil {
		return nil, err
	}
	m := scalecodec.MetadataDecoder{}
	m.Init(utiles.HexToBytes(metaRaw))
	err = m.Process()
	if err != nil {
		return nil, err
	}

	e := scalecodec.EventsDecoder{}
	option := types.ScaleDecoderOption{Metadata: &m.Metadata}
	e.Init(types.ScaleBytes{Data: util.HexToBytes(eventRaw)}, &option)
	e.Process()

	var events []ChainEvent
	b, err := json.Marshal(e.Value)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &events)
	if err != nil {
		return nil, err
	}

	for _, ev := range events {
		if ev.ModuleId != config.BridgeCommon || ev.EventId != config.FungibleTransferEventId {
			continue
		}

		evt := new(EventFungibleTransfer)
		for _, p := range ev.Params {
			switch p.Type {
			case "ChainId":
				var cp ChainIdParam
				x, _ := json.Marshal(p)
				err := json.Unmarshal(x, &cp)
				if err != nil {
					return nil, err
				}
				evt.Destination = cp.Value
			case "DepositNonce":
				var dn DepositNonceParam
				x, _ := json.Marshal(p)
				err := json.Unmarshal(x, &dn)
				if err != nil {
					return nil, err
				}
				evt.DepositNonce = dn.Value
			case "ResourceId":
				val := p.Value.(string)
				r, _ := hexutil.Decode(val)
				var rId [32]byte
				copy(rId[:], r)
				evt.ResourceId = rId
			case "U256":
				amount := new(big.Int)
				b := utiles.HexToBytes(p.Value.(string))
				a := utiles.ReverseBytes(b)
				amount.SetBytes(a)
				evt.Amount = amount
			case "Vec<u8>":
				evt.Recipient, _ = hexutil.Decode(utiles.AddHex(p.Value.(string)))
			}
		}
		evts = append(evts, evt)
	}
	return evts, nil
}

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
		evt.Amount.Mul(evt.Amount, config.DecimalFactor),
		resourceId,
		evt.Recipient,
	), nil
}
