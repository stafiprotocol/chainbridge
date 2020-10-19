// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"testing"

	"github.com/stafiprotocol/go-substrate-rpc-client/types"
)

func TestConnect_QueryStorage(t *testing.T) {
	// Create connection with Alice key
	errs := make(chan error)
	conn := NewConnection(TestEndpoint, "Alice", AliceKey, AliceTestLogger, make(chan int), errs)
	err := conn.Connect()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// Query storage
	var data types.AccountInfo
	_, err = conn.queryStorage("System", "Account", conn.key.PublicKey, nil, &data)
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
