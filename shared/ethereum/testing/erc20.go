// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package ethtest

import (
	"math/big"
	"testing"

	"github.com/ChainSafe/log15"
	"github.com/ethereum/go-ethereum/common"
	utils "github.com/stafiprotocol/chainbridge/shared/ethereum"
)

func DeployMintApproveErc20(t *testing.T, client *utils.Client, erc20Handler common.Address, amount *big.Int) common.Address {
	addr, err := utils.DeployMintApproveErc20(client, erc20Handler, amount)
	if err != nil {
		t.Fatal(err)
	}
	return addr
}

func Erc20AssertBalance(t *testing.T, client *utils.Client, amount *big.Int, erc20Contract, account common.Address) { //nolint:unused,deadcode
	actual, err := utils.Erc20GetBalance(client, erc20Contract, account)
	if err != nil {
		t.Fatal(err)
	}

	if actual.Cmp(amount) != 0 {
		t.Fatalf("Balance mismatch. Expected: %s Got: %s", amount.String(), actual.String())
	}
	log15.Info("Asserted balance", "account", account, "balance", actual, "erc20Contract", erc20Contract.Hex())
}
