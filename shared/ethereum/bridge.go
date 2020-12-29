// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package utils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stafiprotocol/chainbridge/bindings/Bridge"
	"github.com/stafiprotocol/chainbridge/utils/msg"
)

type EventSig string

const (
	Deposit       EventSig = "Deposit(uint8,bytes32,uint64)"
	ProposalEvent EventSig = "ProposalEvent(uint8,uint64,uint8,bytes32,bytes32)"
	ProposalVote  EventSig = "ProposalVote(uint8,uint64,uint8,bytes32)"
)

type ProposalStatus int

const (
	Inactive ProposalStatus = iota
	Active
	Passed
	Executed
	Cancelled
)

func (es EventSig) GetTopic() common.Hash {
	return crypto.Keccak256Hash([]byte(es))
}

//func IsActive(status uint8) bool {
//	return ProposalStatus(status) == Active
//}

func IsFinalized(status uint8) bool {
	return ProposalStatus(status) == Passed
}

func IsExecuted(status uint8) bool {
	return ProposalStatus(status) == Executed
}

func IDAndNonce(srcId msg.ChainId, nonce msg.Nonce) *big.Int {
	var data []byte
	data = append(data, nonce.Big().Bytes()...)
	data = append(data, uint8(srcId))
	return big.NewInt(0).SetBytes(data)
}

func Hash(data []byte) [32]byte {
	return crypto.Keccak256Hash(data)
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
