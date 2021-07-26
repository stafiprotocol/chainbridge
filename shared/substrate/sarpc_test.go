package substrate

import (
	"testing"

	"github.com/ChainSafe/log15"
	"github.com/stretchr/testify/assert"
)

var (
	tlog = log15.Root()
)

const (
	stafiTypesFile = "/Users/fwj/Go/stafi/rtoken-relay/network/stafi.json"
)

func TestSarpcClient_GetChainEvents(t *testing.T) {
	sc, err := NewSarpcClient("wss://stafi-seiya.stafi.io", stafiTypesFile, tlog)
	//sc, err := NewSarpcClient("wss://mainnet-rpc.stafi.io", stafiTypesFile, tlog)
	//sc, err := NewSarpcClient("wss://polkadot-test-rpc.stafi.io", polkaTypesFile, tlog)
	//sc, err := NewSarpcClient(ChainTypeStafi, "ws://127.0.0.1:9944", stafiTypesFile, tlog)
	assert.NoError(t, err)

	for i := 2963170; i < 2963180; i++ {
		evts, err := sc.GetEvents(uint64(i))
		assert.NoError(t, err)
		for _, evt := range evts {
			//if evt.ModuleId != config.RTokenLedgerModuleId {
			//	continue
			//}
			t.Log("i", i)
			t.Log("ModuleId", evt.ModuleId)
			t.Log("EventId", evt.EventId)
			t.Log("Params", evt.Params)
		}
	}
}
