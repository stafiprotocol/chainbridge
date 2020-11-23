package substrate

import (
	"testing"
	"time"

	"github.com/stafiprotocol/chainbridge/utils/blockstore"
	"github.com/stretchr/testify/assert"
)

const StafiMainetEndpoint = "wss://mainnet-rpc.stafi.io"

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
