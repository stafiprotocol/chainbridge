package substrate

import (
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"os"
	"testing"
	"time"

	"github.com/stafiprotocol/chainbridge/utils/blockstore"
	"github.com/stretchr/testify/assert"
)

const (
	StafiMainetEndpoint = "wss://mainnet-rpc.stafi.io"
	SeiyaEndpoint       = "wss://stafi-seiya.stafi.io"
)

func TestListener_start(t *testing.T) {
	errs := make(chan error)
	stop := make(chan int)
	conn, err := NewConnection(seiyaCfg, AliceTestLogger, stop)
	assert.NoError(t, err)
	l := NewListener(conn, "stafi", ThisChain, 100000, AliceTestLogger, &blockstore.EmptyStore{}, stop, errs, nil)
	err = l.start()
	assert.NoError(t, err)

	<-time.After(time.Minute)
}

func TestSeiya(t *testing.T) {
	//errs := make(chan error)
	stop := make(chan int)
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)
	conn, err := NewConnection(seiyaCfg, AliceTestLogger, stop)
	if err != nil {
		t.Fatal(err)
	}

	evts, err := conn.GetEvents(2963178)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(len(evts))
	//assert.NoError(t, err)
	//l := NewListener(conn, "stafi", ThisChain, 2963178, AliceTestLogger, &blockstore.EmptyStore{}, stop, errs, nil)
	//err = l.start()
	//assert.NoError(t, err)
	//
	//<-time.After(30 * time.Minute)
}

func TestSeiya1(t *testing.T) {
	errs := make(chan error)
	stop := make(chan int)
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)
	conn, err := NewConnection(seiyaCfg, AliceTestLogger, stop)

	assert.NoError(t, err)
	l := NewListener(conn, "stafi", ThisChain, 2963177, AliceTestLogger, &blockstore.EmptyStore{}, stop, errs, nil)
	err = l.start()
	assert.NoError(t, err)

	<-time.After(30 * time.Minute)
}
