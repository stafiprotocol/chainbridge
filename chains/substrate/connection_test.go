// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"testing"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge-utils/keystore"
	"github.com/stafiprotocol/chainbridge-utils/msg"
)

const TestEndpoint = "ws://127.0.0.1:9944"

var AliceKey = keystore.TestKeyRing.SubstrateKeys[keystore.AliceKey].AsKeyringPair()
var TestLogLevel = log15.LvlTrace
var AliceTestLogger = newTestLogger("Alice")
var ThisChain msg.ChainId = 1

func TestConnect_CheckChainId(t *testing.T) {
	// Create connection with Alice key
	errs := make(chan error)
	conn := NewConnection(TestEndpoint, "Alice", AliceKey, AliceTestLogger, make(chan int), errs)
	err := conn.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	err = conn.checkChainId(ThisChain)
	if err != nil {
		t.Fatal(err)
	}

	// Ensure no errors were propagated
	select {
	case err := <-errs:
		t.Fatal(err)
	default:
		return
	}
}

func newTestLogger(name string) log15.Logger {
	tLog := log15.Root().New("chain", name)
	tLog.SetHandler(log15.LvlFilterHandler(TestLogLevel, tLog.GetHandler()))
	return tLog
}
