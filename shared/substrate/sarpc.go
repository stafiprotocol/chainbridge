package substrate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/ChainSafe/log15"
	"github.com/gorilla/websocket"
	scalecodec "github.com/itering/scale.go"
	"github.com/itering/scale.go/source"
	"github.com/itering/scale.go/types"
	"github.com/itering/scale.go/utiles"
	"github.com/itering/substrate-api-rpc/pkg/recws"
	"github.com/itering/substrate-api-rpc/rpc"
	"github.com/itering/substrate-api-rpc/util"
	wbskt "github.com/stafiprotocol/chainbridge/shared/substrate/websocket"
	gsrpc "github.com/stafiprotocol/go-substrate-rpc-client"
)

const (
	wsId       = 1
	storageKey = "0x26aa394eea5630e07c48ae0c9558cef780d41e5e16056765bc8461851072c9d7"
)

type SarpcClient struct {
	endpoint           string
	wsPool             wbskt.Pool
	log                log15.Logger
	chainType          string
	metaRaw            string
	typesPath          string
	currentSpecVersion int
	metaDecoder        scalecodec.MetadataDecoder
	eventDecoder       scalecodec.EventsDecoder
}

func NewSarpcClient(endpoint, typesPath string, log log15.Logger) (*SarpcClient, error) {
	api, err := gsrpc.NewSubstrateAPI(endpoint)
	if err != nil {
		return nil, err
	}

	latestHash, err := api.RPC.Chain.GetFinalizedHead()
	if err != nil {
		return nil, err
	}
	log.Info("NewSarpcClient", "latestHash", latestHash.Hex())

	sc := &SarpcClient{
		endpoint:           endpoint,
		wsPool:             nil,
		log:                log,
		metaRaw:            "",
		typesPath:          typesPath,
		currentSpecVersion: 0,
		metaDecoder:        scalecodec.MetadataDecoder{},
		eventDecoder:       scalecodec.EventsDecoder{},
	}

	sc.regCustomTypes()

	err = sc.UpdateMeta(latestHash.Hex())
	if err != nil {
		return nil, err
	}

	return sc, nil
}

func (sc *SarpcClient) regCustomTypes() {
	content, err := ioutil.ReadFile(sc.typesPath)
	if err != nil {
		panic(err)
	}

	types.RuntimeType{}.Reg()
	types.RegCustomTypes(source.LoadTypeRegistry(content))
}

func (sc *SarpcClient) initial() (*wbskt.PoolConn, error) {
	var err error
	if sc.wsPool == nil {
		factory := func() (*recws.RecConn, error) {
			SubscribeConn := &recws.RecConn{KeepAliveTimeout: 10 * time.Second}
			SubscribeConn.Dial(sc.endpoint, nil)
			return SubscribeConn, err
		}
		if sc.wsPool, err = wbskt.NewChannelPool(1, 25, factory); err != nil {
			fmt.Println("NewChannelPool", err)
		}
	}
	if err != nil {
		return nil, err
	}
	conn, err := sc.wsPool.Get()
	return conn, err
}

func (sc *SarpcClient) sendWsRequest(p wbskt.WsConn, v interface{}, action []byte) (err error) {
	if p == nil {
		var pool *wbskt.PoolConn
		if pool, err = sc.initial(); err == nil {
			defer pool.Close()
			p = pool.Conn
		} else {
			return
		}

	}

	if err = p.WriteMessage(websocket.TextMessage, action); err != nil {
		if p != nil {
			p.MarkUnusable()
		}
		return fmt.Errorf("websocket send error: %v", err)
	}
	if err = p.ReadJSON(v); err != nil {
		if p != nil {
			p.MarkUnusable()
		}
		return
	}
	return nil
}

func (sc *SarpcClient) UpdateMeta(blockHash string) error {
	v := &rpc.JsonRpcResult{}
	// runtime version
	if err := sc.sendWsRequest(nil, v, rpc.ChainGetRuntimeVersion(wsId, blockHash)); err != nil {
		return err
	}

	r := v.ToRuntimeVersion()
	if r == nil {
		return fmt.Errorf("runtime version nil")
	}

	// metadata raw
	if sc.metaRaw == "" || r.SpecVersion > sc.currentSpecVersion {
		if err := sc.sendWsRequest(nil, v, rpc.StateGetMetadata(wsId, blockHash)); err != nil {
			return err
		}
		metaRaw, err := v.ToString()
		if err != nil {
			return err
		}
		sc.metaRaw = metaRaw
		sc.currentSpecVersion = r.SpecVersion
		sc.metaDecoder.Init(utiles.HexToBytes(metaRaw))
		if err := sc.metaDecoder.Process(); err != nil {
			return err
		}
	}

	return nil
}

func (sc *SarpcClient) GetBlock(blockHash string) (*rpc.Block, error) {
	v := &rpc.JsonRpcResult{}
	if err := sc.sendWsRequest(nil, v, rpc.ChainGetBlock(wsId, blockHash)); err != nil {
		return nil, err
	}
	rpcBlock := v.ToBlock()
	return &rpcBlock.Block, nil
}

func (sc *SarpcClient) GetBlockHash(blockNum uint64) (string, error) {
	v := &rpc.JsonRpcResult{}
	if err := sc.sendWsRequest(nil, v, rpc.ChainGetBlockHash(wsId, int(blockNum))); err != nil {
		return "", fmt.Errorf("websocket get block hash error: %v", err)
	}

	blockHash, err := v.ToString()
	if err != nil {
		return "", err
	}

	return blockHash, nil
}

func (sc *SarpcClient) GetChainEvents(blockHash string) ([]*ChainEvent, error) {
	err := sc.UpdateMeta(blockHash)
	if err != nil {
		return nil, err
	}

	v := &rpc.JsonRpcResult{}
	if err := sc.sendWsRequest(nil, v, rpc.StateGetStorage(wsId, storageKey, blockHash)); err != nil {
		return nil, fmt.Errorf("websocket get event raw error: %v", err)
	}
	eventRaw, err := v.ToString()
	if err != nil {
		return nil, err
	}

	var events []*ChainEvent
	option := types.ScaleDecoderOption{Metadata: &sc.metaDecoder.Metadata}
	sc.eventDecoder.Init(types.ScaleBytes{Data: util.HexToBytes(eventRaw)}, &option)
	sc.eventDecoder.Process()
	b, err := json.Marshal(sc.eventDecoder.Value)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (sc *SarpcClient) GetEvents(blockNum uint64) ([]*ChainEvent, error) {
	blockHash, err := sc.GetBlockHash(blockNum)
	if err != nil {
		return nil, err
	}

	evts, err := sc.GetChainEvents(blockHash)
	if err != nil {
		return nil, err
	}

	return evts, nil
}
