// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package utils

import (
	events "github.com/stafiprotocol/chainbridge/substrate-events"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
)

type Events struct {
	types.EventRecords
	events.Events
}
