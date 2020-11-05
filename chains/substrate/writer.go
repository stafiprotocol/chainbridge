// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"github.com/stafiprotocol/chainbridge/utils/msg"
)

type writer struct {
}

func NewWriter() *writer {
	return &writer{}
}

func (w *writer) start() error {
	return nil
}

func (w *writer) ResolveMessage(m msg.Message) bool {
	return true
}
