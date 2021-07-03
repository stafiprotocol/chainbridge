// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package ethereum

import (
	"github.com/stafiprotocol/chainbridge/utils/msg"
)

type MockRouter struct {
	msgs chan msg.Message
}

func (r *MockRouter) Send(message msg.Message) error {
	r.msgs <- message
	return nil
}
