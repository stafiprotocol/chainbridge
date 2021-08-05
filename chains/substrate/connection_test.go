// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"os"
	"testing"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge/utils/core"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
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

	seiyaCfg = &core.ChainConfig{
		Name:         "Stafi",
		Id:           1,
		Endpoint:     SeiyaEndpoint,
		From:         From,
		KeystorePath: KeystorePath,
		Insecure:     false,
	}

	localCfg = &core.ChainConfig{
		Name:         "Stafi",
		Id:           1,
		Endpoint:     TestEndpoint,
		From:         From,
		KeystorePath: KeystorePath,
		Insecure:     false,
	}

	mainCfg = &core.ChainConfig{
		Name:         "Stafi",
		Id:           1,
		Endpoint:     StafiMainetEndpoint,
		From:         From,
		KeystorePath: KeystorePath,
		Insecure:     false,
	}
)

func TestConnect_CheckChainId(t *testing.T) {
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	// Create connection with Alice key
	conn, err := NewConnection(seiyaCfg, AliceTestLogger, make(chan int))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	err = conn.checkChainId(ThisChain)
	if err != nil {
		t.Fatal(err)
	}
}

func newTestLogger(name string) log15.Logger {
	tLog := log15.Root().New("chain", name)
	tLog.SetHandler(log15.LvlFilterHandler(TestLogLevel, tLog.GetHandler()))
	return tLog
}
