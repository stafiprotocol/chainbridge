package substrate

import (
	"math/big"
	"testing"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"github.com/stretchr/testify/assert"
)

var (
	AliceKey       = keystore.TestKeyRing.SubstrateKeys[keystore.AliceKey].AsKeyringPair()
	From           = "31yavGB5CVb8EwpqKQaS9XY7JZcfbK6QpWPn5kkweHVpqcov"
	From1          = "31d96Cq9idWQqPq3Ch5BFY84zrThVE3r98M7vG4xYaSWHwsX"
	KeystorePath   = "/Users/fwj/Go/stafi/chainbridge/keys"
	tlog           = log15.Root()
	stafiTypesFile = "/Users/tpkeeper/gowork/stafi/rtoken-relay/network/stafi.json"
)

func TestBatchTransfer(t *testing.T) {
	stop := make(chan int)
	gc, err := NewGsrpcClient("ws://127.0.0.1:9944", AddressTypeAccountId, AliceKey, tlog, stop)
	assert.NoError(t, err)

	less, _ := types.NewAddressFromHexAccountID("0x3673009bdb664a3f3b6d9f69c9dd37fc0473551a249aa48542408b016ec62b2e")
	jun, _ := types.NewAddressFromHexAccountID("0x765f3681fcc33aba624a09833455a3fd971d6791a8f2c57440626cd119530860")
	wen, _ := types.NewAddressFromHexAccountID("0x26db25c52b007221331a844e5335e59874e45b03e81c3d76ff007377c2c17965")
	bao, _ := types.NewAddressFromHexAccountID("0x9c4189297ad2140c85861f64656d1d1318994599130d98b75ff094176d2ca31e")

	addrs := []types.Address{less, jun, wen, bao}

	amount := big.NewInt(1000000000000000)
	value := types.NewUCompact(amount)

	calls := make([]types.Call, 0)
	meta, err := gc.GetLatestMetadata()
	assert.NoError(t, err)

	for _, addr := range addrs {
		call, err := types.NewCall(
			meta,
			"Balances.transfer_keep_alive",
			addr,
			value,
		)
		assert.NoError(t, err)
		calls = append(calls, call)
	}

	ext, err := gc.NewUnsignedExtrinsic("Utility.batch", calls)
	assert.NoError(t, err)

	err = gc.SignAndSubmitTx(ext)
	assert.NoError(t, err)
}

func TestSarpcClient_GetChainEvents(t *testing.T) {
	//sc, err := NewSarpcClient(ChainTypeStafi, "wss://stafi-seiya.stafi.io", stafiTypesFile, tlog)
	sc, err := NewSarpcClient("wss://mainnet-rpc.stafi.io", stafiTypesFile, tlog)
	//sc, err := NewSarpcClient("wss://polkadot-test-rpc.stafi.io", polkaTypesFile, tlog)
	//sc, err := NewSarpcClient(ChainTypeStafi, "ws://127.0.0.1:9944", stafiTypesFile, tlog)
	// stop := make(chan int)
	// sc, err := NewSarpcClient(ChainTypePolkadot, "wss://kusama-rpc.polkadot.io", polkaTypesFile, AddressTypeMultiAddress, AliceKey, tlog, stop)
	if err != nil {
		t.Fatal(err)
	}
	evt, err := sc.GetEvents(10105347)
	assert.NoError(t, err)
	for _, e := range evt {
		t.Log(e.EventId)
		t.Log(e.ModuleId, e.Params)
	}

	// for i := 10105347; i <= 10105347; i++ {
	// 	if i%10 == 0 {
	// 		t.Log("i", i)
	// 	}

	// 	_, err := sc.GetEvents(uint64(i))
	// 	if err != nil {
	// 		t.Fatal(err)
	// 	}

	// 	//assert.NoError(t, err)
	// 	//for _, evt := range evts {
	// 	//	fmt.Println(evt.ModuleId)
	// 	//	fmt.Println(evt.EventId)
	// 	//	fmt.Println(evt.Params)
	// 	//}
	// }
}
