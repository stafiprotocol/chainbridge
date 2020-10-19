// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package ethereum

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// constructErc20ProposalData returns the bytes to construct a proposal suitable for Erc20
func ConstructErc20ProposalData(amount []byte, recipient []byte) []byte {
	var data []byte
	data = append(data, common.LeftPadBytes(amount, 32)...) // amount (uint256)

	recipientLen := big.NewInt(int64(len(recipient))).Bytes()
	data = append(data, common.LeftPadBytes(recipientLen, 32)...) // length of recipient (uint256)
	data = append(data, recipient...)                             // recipient ([]byte)
	return data
}
