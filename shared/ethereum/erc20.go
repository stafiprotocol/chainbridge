// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package utils

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stafiprotocol/chainbridge/bindings/Bridge"
	"github.com/stafiprotocol/chainbridge/bindings/ERC20Handler"
	ERC20 "github.com/stafiprotocol/chainbridge/bindings/ERC20PresetMinterPauser"
	"github.com/stafiprotocol/chainbridge/utils/msg"
)

func Erc20Approve(client *Client, erc20Contract, recipient common.Address, amount *big.Int) error {
	err := client.LockNonceAndUpdate()
	if err != nil {
		return err
	}

	instance, err := ERC20.NewERC20PresetMinterPauser(erc20Contract, client.Client)
	if err != nil {
		return err
	}

	tx, err := instance.Approve(client.Opts, recipient, amount)
	if err != nil {
		return err
	}

	err = WaitForTx(client, tx)
	if err != nil {
		return err
	}

	fmt.Println("Erc20Approve txhash", tx.Hash())

	client.UnlockNonce()

	return nil
}

func Erc20GetBalance(client *Client, erc20Contract, account common.Address) (*big.Int, error) {
	instance, err := ERC20.NewERC20PresetMinterPauser(erc20Contract, client.Client)
	if err != nil {
		return nil, err
	}

	bal, err := instance.BalanceOf(client.CallOpts, account)
	if err != nil {
		return nil, err

	}
	return bal, nil

}

func Erc20AddMinter(client *Client, erc20Contract, handler common.Address) error {
	err := client.LockNonceAndUpdate()
	if err != nil {
		return err
	}

	instance, err := ERC20.NewERC20PresetMinterPauser(erc20Contract, client.Client)
	if err != nil {
		return err
	}

	role, err := instance.MINTERROLE(client.CallOpts)
	if err != nil {
		return err
	}

	tx, err := instance.GrantRole(client.Opts, role, handler)
	if err != nil {
		return err
	}

	err = WaitForTx(client, tx)
	if err != nil {
		return err
	}

	fmt.Println("Erc20AddMinter txhash", tx.Hash())

	client.UnlockNonce()

	return nil
}

func SetResourceAndBurnable(client *Client, bridge, handler common.Address, rIds [][32]byte, tokenContractAddresses, burnablTokenContractAddresses []common.Address) error {
	instance, err := Bridge.NewBridge(bridge, client.Client)
	if err != nil {
		return err
	}

	err = client.LockNonceAndUpdate()
	if err != nil {
		return err
	}

	tx, err := instance.AdminSetResourceAndBurnable(client.Opts, handler, rIds, tokenContractAddresses, burnablTokenContractAddresses)
	if err != nil {
		return err
	}

	err = WaitForTx(client, tx)
	if err != nil {
		return err
	}

	client.UnlockNonce()

	fmt.Println("SetResourceAndBurnable txhash", tx.Hash())

	return nil
}

func UpdateNonce(client *Client) error {
	newNonce, err := client.Client.PendingNonceAt(context.Background(), client.CallOpts.From)
	if err != nil {
		return err
	}

	client.Opts.Nonce = big.NewInt(int64(newNonce))

	return nil
}

func GetDepositNonce(client *Client, bridge common.Address, chain msg.ChainId) (uint64, error) {
	instance, err := Bridge.NewBridge(bridge, client.Client)
	if err != nil {
		return 0, err
	}

	count, err := instance.DepositCounts(client.CallOpts, uint8(chain))
	if err != nil {
		return 0, err
	}

	return count, nil
}

func CreateErc20Deposit(client *Client, bridge common.Address, destId msg.ChainId, recipient []byte, amount *big.Int, rId msg.ResourceId) error {
	data := ConstructErc20DepositData(recipient, amount)

	bridgeInstance, err := Bridge.NewBridge(bridge, client.Client)
	if err != nil {
		return err
	}

	tx, err := bridgeInstance.Deposit(
		client.Opts,
		uint8(destId),
		rId,
		data,
	)

	if err != nil {
		return err
	}

	fmt.Println("CreateErc20Deposit txhash", tx.Hash())

	return WaitForTx(client, tx)
}

func Erc20GetAllowance(client *Client, erc20Contract, owner, spender common.Address) (*big.Int, error) {
	instance, err := ERC20.NewERC20PresetMinterPauser(erc20Contract, client.Client)
	if err != nil {
		return nil, err
	}

	amount, err := instance.Allowance(client.CallOpts, owner, spender)
	if err != nil {
		return nil, err
	}

	return amount, nil
}

func Erc20GetResourceId(client *Client, handler common.Address, rId msg.ResourceId) (common.Address, error) {
	instance, err := ERC20Handler.NewERC20Handler(handler, client.Client)
	if err != nil {
		return ZeroAddress, err
	}

	addr, err := instance.ResourceIDToTokenContractAddress(client.CallOpts, rId)
	if err != nil {
		return ZeroAddress, err
	}

	return addr, nil
}

func ConstructErc20DepositData(destRecipient []byte, amount *big.Int) []byte {
	var data []byte
	data = append(data, math.PaddedBigBytes(amount, 32)...)
	data = append(data, math.PaddedBigBytes(big.NewInt(int64(len(destRecipient))), 32)...)
	data = append(data, destRecipient...)
	return data
}
