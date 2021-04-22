// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"encoding/hex"
	"fmt"
	"github.com/stafiprotocol/chainbridge/utils/crypto/sr25519"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
	"os"
	"testing"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge/utils/msg"
)

const TestEndpoint = "ws://127.0.0.1:9944"

var (
	TestLogLevel                = log15.LvlTrace
	AliceTestLogger             = newTestLogger("Alice")
	SourceChain     msg.ChainId = 2
	ThisChain       msg.ChainId = 1
	AliceKey                    = keystore.TestKeyRing.SubstrateKeys[keystore.AliceKey].AsKeyringPair()
	From                        = "31yavGB5CVb8EwpqKQaS9XY7JZcfbK6QpWPn5kkweHVpqcov"
	KeystorePath                = "/Users/fwj/Go/stafi/chainbridge/keys"
)

func TestConnect_CheckChainId(t *testing.T) {
	// Create connection with Alice key
	errs := make(chan error)
	conn := NewConnection(TestEndpoint, "Alice", AliceKey, nil, AliceTestLogger, make(chan int), errs)
	err := conn.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	err = conn.checkChainId(ThisChain)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccountInfo(t *testing.T) {
	//the password which used to encrypt keystore file, remove it after this test pass
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kp, err := keystore.KeypairFromAddress(From, keystore.SubChain, KeystorePath, false)
	if err != nil {
		t.Fatal(err)
	}

	krp := kp.(*sr25519.Keypair).AsKeyringPair()
	errs := make(chan error)
	conn := NewConnection(TestEndpoint, "Alice", krp, nil, AliceTestLogger, make(chan int), errs)
	err = conn.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	ac, err := conn.GetAccountInfo()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ac)
}

func TestConnection_SimpleTransfer(t *testing.T) {
	//the password which used to encrypt keystore file, remove it after this test pass
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kp, err := keystore.KeypairFromAddress(From, keystore.SubChain, KeystorePath, false)
	if err != nil {
		t.Fatal(err)
	}

	krp := kp.(*sr25519.Keypair).AsKeyringPair()
	errs := make(chan error)
	conn := NewConnection(TestEndpoint, "Alice", krp, nil, AliceTestLogger, make(chan int), errs)
	err = conn.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	amount := uint64(1000000000000 * 10)
	//// Create a call, transferring 12345 units to Bob
	bob, err := types.NewAddressFromHexAccountID("0xf6241901b8e0048421427ef6cd3513865c2b6d2ad3ca2c3d95d28dfca2b4f722")
	if err != nil {
		panic(err)
	}

	ext, err := conn.NewUnsignedExtrinsic("Balances.transfer", bob, types.NewUCompactFromUInt(amount))
	err = conn.signExtrinsic(ext)
	if err != nil {
		t.Fatal(err)
	}
	x, _ := ext.MarshalJSON()
	fmt.Println(hex.EncodeToString(x))
	//ext.UnmarshalJSON()
	//ext.MarshalJSON()
	//fmt.Println(ext.IsSigned())
	//re, err := types.EncodeToHexString(ext)
	//x, err := json.Marshal(re)
	//assert.NoError(t, err)
	//fmt.Println(x)
	//err = conn.SignAndSubmitTx(ext)

}

func newTestLogger(name string) log15.Logger {
	tLog := log15.Root().New("chain", name)
	tLog.SetHandler(log15.LvlFilterHandler(TestLogLevel, tLog.GetHandler()))
	return tLog
}
