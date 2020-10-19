// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge-utils/keystore"
	"github.com/stafiprotocol/chainbridge-utils/msg"
	"github.com/stafiprotocol/chainbridge/config"
	utils "github.com/stafiprotocol/chainbridge/shared/substrate"
	"os"
	"testing"
)

const TestEndpoint = "ws://127.0.0.1:9944"

var AliceKey = keystore.TestKeyRing.SubstrateKeys[keystore.AliceKey].AsKeyringPair()
var BobKey = keystore.TestKeyRing.SubstrateKeys[keystore.BobKey].AsKeyringPair()

var TestLogLevel = log15.LvlTrace
var AliceTestLogger = newTestLogger("Alice")
var BobTestLogger = newTestLogger("Bob")

var ThisChain msg.ChainId = 1

type testContext struct {
	client         *utils.Client
	listener       *listener
	router         *mockRouter
	writerAlice    *writer
	writerBob      *writer
	latestOutNonce msg.Nonce
	latestInNonce  msg.Nonce
	lSysErr        chan error
	wSysErr        chan error
}

var context testContext
func TestMain(m *testing.M) {
	client, err := utils.CreateClient(AliceKey, TestEndpoint)
	if err != nil {
		panic(err)
	}

	var nativeTokenId [32]byte

	err = client.GetConst(config.BridgeSwap, config.NativeTokenId, &nativeTokenId)
	if err != nil {
		panic(err)
	}

	aliceConn, bobConn, wSysErr, err := createAliceAndBobConnections()
	if err != nil {
		panic(err)
	}
	l, lSysErr, r, err := newTestListener(client, aliceConn)
	if err != nil {
		panic(err)
	}
	alice := NewWriter(aliceConn, AliceTestLogger, wSysErr, nil)
	bob := NewWriter(bobConn, BobTestLogger, wSysErr, nil)
	context = testContext{
		client:         client,
		listener:       l,
		router:         r,
		writerAlice:    alice,
		writerBob:      bob,
		latestInNonce:  0,
		latestOutNonce: 0,
		lSysErr:        lSysErr,
		wSysErr:        wSysErr,
	}

	os.Exit(m.Run())
}

func newTestLogger(name string) log15.Logger {
	tLog := log15.Root().New("chain", name)
	tLog.SetHandler(log15.LvlFilterHandler(TestLogLevel, tLog.GetHandler()))
	return tLog
}

// createAliceConnection creates and starts a connection with the Alice keypair
func createAliceConnection() (*Connection, chan error, error) {
	sysErr := make(chan error)
	alice := NewConnection(TestEndpoint, "Alice", AliceKey, AliceTestLogger, make(chan int), sysErr)
	err := alice.Connect()
	if err != nil {
		return nil, nil, err
	}
	return alice, sysErr, err
}

// createAliceAndBobConnections creates and calls `Connect()` on two Connections using the Alice and Bob keypairs
func createAliceAndBobConnections() (*Connection, *Connection, chan error, error) {
	alice, sysErr, err := createAliceConnection()
	if err != nil {
		return nil, nil, nil, err
	}

	bob := NewConnection(TestEndpoint, "Bob", BobKey, AliceTestLogger, make(chan int), sysErr)
	err = bob.Connect()
	if err != nil {
		return nil, nil, nil, err
	}

	return alice, bob, sysErr, nil
}
