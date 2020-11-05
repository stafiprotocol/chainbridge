package substrate

import (
	"fmt"
	"github.com/stafiprotocol/chainbridge/utils/blockstore"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const StafiMainetEndpoint = "wss://mainnet-rpc.stafi.io"

func TestListener_GetEventsAt(t *testing.T) {
	errs := make(chan error)
	stop := make(chan int)
	conn := NewConnection(StafiMainetEndpoint, "Alice", nil, AliceTestLogger, stop, errs)

	l := NewListener(conn, "stafi", ThisChain, 1, AliceTestLogger, &blockstore.EmptyStore{}, stop, errs, nil)
	startBlock := 1
	endBlock := 100
	for i := startBlock; i <= endBlock; i++ {
		if i%10 == 0 {
			fmt.Println("i:", i)
		}
		_, err := l.GetEventsAt(uint64(i))
		assert.NoError(t, err)
	}
}

func TestListener_start(t *testing.T) {
	errs := make(chan error)
	stop := make(chan int)
	conn := NewConnection(StafiMainetEndpoint, "Alice", nil, AliceTestLogger, stop, errs)
	err := conn.Connect()
	assert.NoError(t, err)
	l := NewListener(conn, "stafi", ThisChain, 100000, AliceTestLogger, &blockstore.EmptyStore{}, stop, errs, nil)
	err = l.start()
	assert.NoError(t, err)

	<-time.After(time.Minute)
}
