// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package ethereum

import (
	"github.com/stafiprotocol/chainbridge/bindings/Bridge"
	"github.com/stafiprotocol/chainbridge/bindings/ERC20Handler"
	utils "github.com/stafiprotocol/chainbridge/shared/ethereum"
	ethtest "github.com/stafiprotocol/chainbridge/shared/ethereum/testing"
	"github.com/stafiprotocol/chainbridge/utils/blockstore"
	"github.com/stafiprotocol/chainbridge/utils/msg"
	"testing"
)

type MockRouter struct {
	msgs chan msg.Message
}

func (r *MockRouter) Send(message msg.Message) error {
	r.msgs <- message
	return nil
}

func createTestListener(t *testing.T, config *Config, contracts *utils.DeployedContracts, stop <-chan int, sysErr chan<- error) (*listener, *MockRouter) {
	// Create copy and add deployed contract addresses
	newConfig := *config
	newConfig.bridgeContract = contracts.BridgeAddress
	newConfig.erc20HandlerContract = contracts.ERC20HandlerAddress

	conn := newLocalConnection(t, &newConfig)
	latestBlock, err := conn.LatestBlock()
	if err != nil {
		t.Fatal(err)
	}
	newConfig.startBlock = latestBlock

	bridgeContract, err := Bridge.NewBridge(newConfig.bridgeContract, conn.Client())
	if err != nil {
		t.Fatal(err)
	}
	erc20HandlerContract, err := ERC20Handler.NewERC20Handler(newConfig.erc20HandlerContract, conn.Client())
	if err != nil {
		t.Fatal(err)
	}

	router := &MockRouter{msgs: make(chan msg.Message)}
	listener := NewListener(conn, &newConfig, TestLogger, &blockstore.EmptyStore{}, stop, sysErr)
	listener.setContracts(bridgeContract, erc20HandlerContract)
	listener.setRouter(router)
	// Start the listener
	err = listener.start()
	if err != nil {
		t.Fatal(err)
	}

	return listener, router
}

func TestListener_start_stop(t *testing.T) {
	client := ethtest.NewClient(t, TestEndpoint, AliceKp)
	contracts := deployTestContracts(t, client, aliceTestConfig.id)
	stop := make(chan int)
	l, _ := createTestListener(t, aliceTestConfig, contracts, stop, nil)

	err := l.start()
	if err != nil {
		t.Fatal(err)
	}

	// Initiate shutdown
	close(stop)
}
