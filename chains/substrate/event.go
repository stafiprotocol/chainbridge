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
	"github.com/stafiprotocol/chainbridge/config"
	"github.com/stafiprotocol/chainbridge/utils/ethereum"
	"github.com/stafiprotocol/chainbridge/utils/msg"
)

const (
	wsId       = 1
	storageKey = "0x26aa394eea5630e07c48ae0c9558cef780d41e5e16056765bc8461851072c9d7"
)

var (
	CurrentRuntimeSpecVersion = 6
	metaRaw                   string
	metaDecoder               = scalecodec.MetadataDecoder{}
)

func (l *listener) GetEventsAt(blockNum uint64) ([]*EventFungibleTransfer, error) {
	evts := make([]*EventFungibleTransfer, 0)
	v := &rpc.JsonRpcResult{}
	// Block Hash
	if err := websocket.SendWsRequest(l.wsconn, v, rpc.ChainGetBlockHash(wsId, int(blockNum))); err != nil {
		return nil, fmt.Errorf("websocket get block hash error: %v", err)
	}
	blockHash, err := v.ToString()
	if err != nil {
		return nil, err
	}

	// runtime version
	if err := websocket.SendWsRequest(l.wsconn, v, rpc.ChainGetRuntimeVersion(wsId, blockHash)); err != nil {
		return nil, fmt.Errorf("websocket get runtime version error: %v", err)
	}
	r := v.ToRuntimeVersion()
	if r == nil {
		return nil, fmt.Errorf("runtime version nil")
	}

	// metadata raw
	if metaRaw == "" || r.SpecVersion > CurrentRuntimeSpecVersion {
		if err := websocket.SendWsRequest(l.wsconn, v, rpc.StateGetMetadata(wsId, blockHash)); err != nil {
			return nil, fmt.Errorf("websocket get metadata error: %v", err)
		}
		metaRaw, err = v.ToString()
		if err != nil {
			return nil, err
		}
		metaDecoder.Init(utiles.HexToBytes(metaRaw))
		err = metaDecoder.Process()
		if err != nil {
			return nil, err
		}
		CurrentRuntimeSpecVersion = r.SpecVersion
	}

	// event raw
	if err := websocket.SendWsRequest(l.wsconn, v, rpc.StateGetStorage(wsId, storageKey, blockHash)); err != nil {
		return nil, fmt.Errorf("websocket get event raw error: %v", err)
	}
	eventRaw, err := v.ToString()
	if err != nil {
		return nil, err
	}

	// parse event raw into []ChainEvent
	e := scalecodec.EventsDecoder{}
	option := types.ScaleDecoderOption{Metadata: &metaDecoder.Metadata}
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

	l.log.Trace("block", "eventsNum", len(events), "blockNum", blockNum)

	for _, ev := range events {
		if ev.ModuleId != config.BridgeCommon || ev.EventId != config.FungibleTransferEventId {
			continue
		}

		evt := new(EventFungibleTransfer)
		skip := false
		for _, p := range ev.Params {
			switch p.Type {
			case "ChainId":
				var cp ChainIdParam
				x, _ := json.Marshal(p)
				json.Unmarshal(x, &cp)
				evt.Destination = cp.Value
			case "DepositNonce":
				var dn DepositNonceParam
				x, _ := json.Marshal(p)
				json.Unmarshal(x, &dn)
				evt.DepositNonce = dn.Value
			case "ResourceId":
				val := p.Value.(string)
				r, _ := hexutil.Decode(val)
				evt.ResourceId = msg.ResourceIdFromSlice(r)
			case "U256":
				amount := new(big.Int)
				b := utiles.HexToBytes(p.Value.(string))
				rb := utiles.ReverseBytes(b)
				amount.SetBytes(rb)
				evt.Amount = amount
			case "Vec<u8>":
				addr := utiles.AddHex(p.Value.(string))
				if !ethereum.IsAddressValid(addr) {
					l.log.Warn("GetEventsAt", "Recipient address is not valid: ", addr, "blockNum", l.latestBlock.Height)
					skip = true
					break
				}
				evt.Recipient, _ = hexutil.Decode(addr)
			case "AccountId":
				l.log.Info("GetEventsAt", "from", p.Value.(string))
			default:
				l.log.Warn("GetEventsAt", "EventFungibleTransfer got an unexpected type", p.Type, "blockNum", l.latestBlock.Height)
				skip = true
				break
			}
		}
		if !skip {
			evts = append(evts, evt)
		}

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
