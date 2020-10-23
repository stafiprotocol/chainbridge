// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package ethtest

import (
	"context"
	"math/big"
	"testing"

	"github.com/stafiprotocol/chainbridge-utils/crypto/secp256k1"
	utils "github.com/stafiprotocol/chainbridge/shared/ethereum"
)

func NewClient(t *testing.T, endpoint string, kp *secp256k1.Keypair) *utils.Client {
	client, err := utils.NewClient(endpoint, kp)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func GetLatestBlock(t *testing.T, client *utils.Client) *big.Int {
	block, err := client.Client.BlockByNumber(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	return block.Number()
}
