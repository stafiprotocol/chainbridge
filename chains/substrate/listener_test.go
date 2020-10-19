// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"github.com/stafiprotocol/chainbridge-utils/blockstore"
	"github.com/stafiprotocol/chainbridge-utils/msg"
	utils "github.com/stafiprotocol/chainbridge/shared/substrate"
)

type mockRouter struct {
	msgs chan msg.Message
}

func (r *mockRouter) Send(message msg.Message) error {
	r.msgs <- message
	return nil
}

func newTestListener(client *utils.Client, conn *Connection) (*listener, chan error, *mockRouter, error) {
	r := &mockRouter{msgs: make(chan msg.Message)}

	startBlock, err := client.LatestBlock()
	if err != nil {
		return nil, nil, nil, err
	}

	errs := make(chan error)
	l := NewListener(conn, "Alice", 1, startBlock, AliceTestLogger, &blockstore.EmptyStore{}, make(chan int), errs, nil)
	l.setRouter(r)
	err = l.start()
	if err != nil {
		return nil, nil, nil, err
	}

	return l, errs, r, nil
}
