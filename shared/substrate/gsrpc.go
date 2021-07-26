package substrate

import (
	"errors"
	"fmt"
	"github.com/stafiprotocol/chainbridge/config"
	"time"

	"github.com/ChainSafe/log15"
	gsrpc "github.com/stafiprotocol/go-substrate-rpc-client"
	gsrpcConfig "github.com/stafiprotocol/go-substrate-rpc-client/config"
	"github.com/stafiprotocol/go-substrate-rpc-client/rpc/author"
	"github.com/stafiprotocol/go-substrate-rpc-client/signature"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
)

type GsrpcClient struct {
	endpoint    string
	addressType string
	api         *gsrpc.SubstrateAPI
	key         *signature.KeyringPair
	genesisHash types.Hash
	stop        <-chan int
	log         log15.Logger
}

func NewGsrpcClient(endpoint, addressType string, key *signature.KeyringPair, log log15.Logger, stop <-chan int) (*GsrpcClient, error) {
	log.Info("Connecting to substrate chain with Gsrpc", "endpoint", endpoint)

	if addressType != AddressTypeAccountId && addressType != AddressTypeMultiAddress {
		return nil, errors.New("addressType not supported")
	}

	api, err := gsrpc.NewSubstrateAPI(endpoint)
	if err != nil {
		return nil, err
	}

	gsrpcConfig.SetSubscribeTimeout(15 * time.Second)

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		return nil, err
	}

	return &GsrpcClient{
		endpoint:    endpoint,
		addressType: addressType,
		api:         api,
		key:         key,
		genesisHash: genesisHash,
		stop:        stop,
		log:         log,
	}, nil
}

func (gc *GsrpcClient) FlashApi() (*gsrpc.SubstrateAPI, error) {
	_, err := gc.api.RPC.Chain.GetBlockHashLatest()
	if err != nil {
		var api *gsrpc.SubstrateAPI
		for i := 0; i < 3; i++ {
			api, err = gsrpc.NewSubstrateAPI(gc.endpoint)
			if err == nil {
				break
			} else {
				time.Sleep(time.Millisecond * 100)
			}
		}
		if api != nil {
			gc.api = api
		}
	}
	return gc.api, nil
}

func (gc *GsrpcClient) Address() string {
	return gc.key.Address
}

func (gc *GsrpcClient) GetLatestBlockNumber() (uint64, error) {
	h, err := gc.GetHeaderLatest()
	if err != nil {
		return 0, err
	}

	return uint64(h.Number), nil
}

func (gc *GsrpcClient) GetFinalizedBlockNumber() (uint64, error) {
	hash, err := gc.GetFinalizedHead()
	if err != nil {
		return 0, err
	}

	header, err := gc.GetHeader(hash)
	if err != nil {
		return 0, err
	}

	return uint64(header.Number), nil
}

func (gc *GsrpcClient) GetHeaderLatest() (*types.Header, error) {
	api, err := gc.FlashApi()
	if err != nil {
		return nil, err
	}
	return api.RPC.Chain.GetHeaderLatest()
}

func (gc *GsrpcClient) GetFinalizedHead() (types.Hash, error) {
	api, err := gc.FlashApi()
	if err != nil {
		return types.NewHash([]byte{}), err
	}
	return api.RPC.Chain.GetFinalizedHead()
}

func (gc *GsrpcClient) GetHeader(blockHash types.Hash) (*types.Header, error) {
	api, err := gc.FlashApi()
	if err != nil {
		return nil, err
	}
	return api.RPC.Chain.GetHeader(blockHash)
}

func (gc *GsrpcClient) GetBlockNumber(blockHash types.Hash) (uint64, error) {
	head, err := gc.GetHeader(blockHash)
	if err != nil {
		return 0, err
	}

	return uint64(head.Number), nil
}

// queryStorage performs a storage lookup. Arguments may be nil, result must be a pointer.
func (gc *GsrpcClient) QueryStorage(prefix, method string, arg1, arg2 []byte, result interface{}) (bool, error) {
	meta, err := gc.GetLatestMetadata()
	if err != nil {
		return false, err
	}

	key, err := types.CreateStorageKey(meta, prefix, method, arg1, arg2)
	if err != nil {
		return false, err
	}

	api, err := gc.FlashApi()
	if err != nil {
		return false, err
	}

	ok, err := api.RPC.State.GetStorageLatest(key, result)
	if err != nil {
		return false, err
	}

	return ok, nil
}

func (gc *GsrpcClient) GetLatestMetadata() (*types.Metadata, error) {
	api, err := gc.FlashApi()
	if err != nil {
		return nil, err
	}
	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		return nil, err
	}

	return meta, nil
}

func (gc *GsrpcClient) GetLatestRuntimeVersion() (*types.RuntimeVersion, error) {
	api, err := gc.FlashApi()
	if err != nil {
		return nil, err
	}
	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		return nil, err
	}

	return rv, nil
}

func (gc *GsrpcClient) GetLatestNonce() (types.U32, error) {
	ac, err := gc.GetAccountInfo()
	if err != nil {
		return 0, err
	}

	return ac.Nonce, nil
}

func (gc *GsrpcClient) GetAccountInfo() (*types.AccountInfo, error) {
	ac := new(types.AccountInfo)
	exist, err := gc.QueryStorage("System", "Account", gc.key.PublicKey, nil, &ac)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, errors.New("account not exist")
	}

	return ac, nil
}

func (gc *GsrpcClient) NewUnsignedExtrinsic(callMethod string, args ...interface{}) (interface{}, error) {
	gc.log.Debug("Submitting substrate call...", "callMethod", callMethod, "addressType", gc.addressType, "sender", gc.key.Address)
	meta, err := gc.GetLatestMetadata()
	if err != nil {
		return nil, err
	}

	call, err := types.NewCall(meta, callMethod, args...)
	if err != nil {
		return nil, err
	}

	if gc.addressType == AddressTypeAccountId {
		unsignedExt := types.NewExtrinsic(call)
		return &unsignedExt, nil
	} else if gc.addressType == AddressTypeMultiAddress {
		unsignedExt := types.NewExtrinsicMulti(call)
		return &unsignedExt, nil
	} else {
		return nil, errors.New("addressType not supported")
	}
}

func (gc *GsrpcClient) SignAndSubmitTx(ext interface{}) error {
	err := gc.signExtrinsic(ext)
	if err != nil {
		return err
	}

	api, err := gc.FlashApi()
	if err != nil {
		return err
	}
	// Do the transfer and track the actual status
	sub, err := api.RPC.Author.SubmitAndWatch(ext)
	if err != nil {
		return err
	}
	gc.log.Trace("Extrinsic submission succeeded")
	defer sub.Unsubscribe()

	return gc.watchSubmission(sub)
}

func (gc *GsrpcClient) watchSubmission(sub *author.ExtrinsicStatusSubscription) error {
	for {
		select {
		case <-gc.stop:
			return TerminatedError
		case status := <-sub.Chan():
			switch {
			case status.IsInBlock:
				gc.log.Info("Extrinsic included in block", "block", status.AsInBlock.Hex())
				return nil
			case status.IsRetracted:
				return fmt.Errorf("extrinsic retracted: %s", status.AsRetracted.Hex())
			case status.IsDropped:
				return fmt.Errorf("extrinsic dropped from network")
			case status.IsInvalid:
				return fmt.Errorf("extrinsic invalid")
			}
		case err := <-sub.Err():
			gc.log.Trace("Extrinsic subscription error", "err", err)
			return err
		}
	}
}

func (gc *GsrpcClient) signExtrinsic(xt interface{}) error {
	rv, err := gc.GetLatestRuntimeVersion()
	if err != nil {
		return err
	}

	nonce, err := gc.GetLatestNonce()
	if err != nil {
		return err
	}

	o := types.SignatureOptions{
		BlockHash:          gc.genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        gc.genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(0),
		TransactionVersion: rv.TransactionVersion,
	}

	if ext, ok := xt.(*types.Extrinsic); ok {
		gc.log.Info("signExtrinsic", "addressType", gc.addressType)
		err = ext.Sign(*gc.key, o)
		if err != nil {
			return err
		}
	} else if ext, ok := xt.(*types.ExtrinsicMulti); ok {
		gc.log.Info("signExtrinsic", "addressType1", gc.addressType)
		err = ext.Sign(*gc.key, o)
		if err != nil {
			return err
		}
	} else {
		return errors.New("extrinsic cast error")
	}

	return nil
}

func (gc *GsrpcClient) PublicKey() []byte {
	return gc.key.PublicKey
}

func (gc *GsrpcClient) ChainId() (uint8, error) {
	api, err := gc.FlashApi()
	if err != nil {
		return 0, err
	}
	var chainId uint8
	err = api.RPC.State.GetConst(config.BridgeCommon, config.ChainIdentity, &chainId)
	if err != nil {
		return 0, err
	}
	return chainId, nil
}
