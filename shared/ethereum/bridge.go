// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package utils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stafiprotocol/chainbridge/utils/msg"
)

type EventSig string

const (
	Deposit EventSig = "Deposit(uint8,bytes32,uint64)"
)

type ProposalStatus int

const (
	Inactive ProposalStatus = iota
	Active
	Executed
	Cancelled
)

func (es EventSig) GetTopic() common.Hash {
	return crypto.Keccak256Hash([]byte(es))
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
