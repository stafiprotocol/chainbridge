// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package chains

import (
	"github.com/stafiprotocol/chainbridge/utils/msg"
)

type Router interface {
	Send(message msg.Message) error
	SupportChainId(chainId msg.ChainId) bool
}
