// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"fmt"
	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge-utils/msg"
	"github.com/stafiprotocol/chainbridge/config"
	events "github.com/stafiprotocol/chainbridge/substrate-events"
)

type eventName string
type eventHandler func(interface{}, log15.Logger) (msg.Message, error)

const FungibleTransfer eventName = "FungibleTransfer"

var Subscriptions = []struct {
	name    eventName
	handler eventHandler
}{{FungibleTransfer, fungibleTransferHandler}}

func fungibleTransferHandler(evtI interface{}, log log15.Logger) (msg.Message, error) {
	evt, ok := evtI.(events.EventFungibleTransfer)
	if !ok {
		return msg.Message{}, fmt.Errorf("failed to cast EventFungibleTransfer type")
	}

	resourceId := msg.ResourceId(evt.ResourceId)
	log.Info("Got fungible transfer event!", "destination", evt.Destination, "resourceId", resourceId.Hex(), "amount", evt.Amount)

	return msg.NewFungibleTransfer(0, msg.ChainId(evt.Destination), msg.Nonce(evt.DepositNonce),
		evt.Amount.Mul(evt.Amount.Int, config.DecimalFactor), resourceId, evt.Recipient), nil
}
