// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"github.com/ChainSafe/log15"
	metrics "github.com/stafiprotocol/chainbridge-utils/metrics/types"
	"github.com/stafiprotocol/chainbridge-utils/msg"
)

type writer struct {
	conn    *Connection
	log     log15.Logger
	sysErr  chan<- error
	metrics *metrics.ChainMetrics
}

func NewWriter(conn *Connection, log log15.Logger, sysErr chan<- error, m *metrics.ChainMetrics) *writer {
	return &writer{
		conn:    conn,
		log:     log,
		sysErr:  sysErr,
		metrics: m,
	}
}

func (w *writer) start() error {
	return nil
}

func (w *writer) ResolveMessage(m msg.Message) bool {
	w.log.Info("msg resolved: ", m)
	return true
}
