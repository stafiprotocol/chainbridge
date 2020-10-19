// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package utils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stafiprotocol/chainbridge-utils/msg"
	"github.com/stafiprotocol/chainbridge/bindings/Bridge"
)

func RegisterResource(client *Client, bridge, handler common.Address, rId msg.ResourceId, addr common.Address) error {
	instance, err := Bridge.NewBridge(bridge, client.Client)
	if err != nil {
		return err
	}

	err = client.LockNonceAndUpdate()
	if err != nil {
		return err
	}

	tx, err := instance.AdminSetResource(client.Opts, handler, rId, addr)
	if err != nil {
		return err
	}

	err = WaitForTx(client, tx)
	if err != nil {
		return err
	}

	client.UnlockNonce()

	return nil
}

func SetBurnable(client *Client, bridge, handler, contract common.Address) error {
	instance, err := Bridge.NewBridge(bridge, client.Client)
	if err != nil {
		return err
	}

	err = client.LockNonceAndUpdate()
	if err != nil {
		return err
	}

	tx, err := instance.AdminSetBurnable(client.Opts, handler, contract)
	if err != nil {
		return err
	}

	err = WaitForTx(client, tx)
	if err != nil {
		return err
	}

	client.UnlockNonce()

	return nil
}

func IDAndNonce(srcId msg.ChainId, nonce msg.Nonce) *big.Int {
	var data []byte
	data = append(data, nonce.Big().Bytes()...)
	data = append(data, uint8(srcId))
	return big.NewInt(0).SetBytes(data)
}
