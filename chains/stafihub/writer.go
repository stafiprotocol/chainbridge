// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package stafihub

import (
	"encoding/hex"
	"errors"
	"fmt"
	errType "github.com/cosmos/cosmos-sdk/types/errors"
	"math/big"
	"strings"
	"time"

	"github.com/ChainSafe/log15"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/shopspring/decimal"
	stafiHubXBridgeTypes "github.com/stafihub/stafihub/x/bridge/types"
	"github.com/stafiprotocol/chainbridge/utils/msg"
)

const (
	msgLimit = 4096
)

var ErrorTerminated = errors.New("terminated")

type writer struct {
	conn     *Connection
	log      log15.Logger
	sysErr   chan<- error
	msgChan  chan msg.Message
	stop     <-chan int
	decimals map[string]decimal.Decimal
}

func NewWriter(conn *Connection, log log15.Logger, sysErr chan<- error, stop <-chan int, decimals map[string]decimal.Decimal) *writer {
	return &writer{
		conn:     conn,
		log:      log,
		sysErr:   sysErr,
		msgChan:  make(chan msg.Message, msgLimit),
		stop:     stop,
		decimals: decimals,
	}
}

func (w *writer) start() error {
	go func() {
		for {
			select {
			case <-w.stop:
				close(w.msgChan)
				w.log.Info("writer stopped")
				return
			case msg := <-w.msgChan:
				result := w.processMessage(msg)
				w.log.Info("processMessage", "result", result)
				if !result {
					w.sysErr <- fmt.Errorf("processMessage failed")
				}
			}
		}
	}()

	return nil
}

func (w *writer) ResolveMessage(m msg.Message) bool {
	w.msgChan <- m
	return true
}

func (w *writer) processMessage(m msg.Message) bool {
	w.log.Info("ResolveMessage", "Name", w.conn.name, "Destination", m.Destination)
	switch m.Type {
	case msg.FungibleTransfer:
		bigAmt := big.NewInt(0).SetBytes(m.Payload[0].([]byte))
		//should not have 0x prefix and length must 64
		resourceIdStr := strings.ToLower(m.ResourceId.Hex())
		if len(resourceIdStr) != 64 {
			fmt.Errorf("resourceId  length  must be 64")
			return false
		}

		d, ok := w.decimals[resourceIdStr]
		if !ok {
			d, ok = w.decimals[decimalDefault]
			if !ok {
				fmt.Errorf("failed to get decimal")
				return false
			}
		}
		amount := *decimal.NewFromBigInt(bigAmt, 0).Div(d).BigInt()
		depositNonce := m.DepositNonce.Big().Uint64()
		recipient := m.Payload[1].([]byte)
		reciver, err := types.AccAddressFromHex(hex.EncodeToString(recipient))
		if err != nil {
			return false
		}

		voteMsg := stafiHubXBridgeTypes.NewMsgVoteProposal(w.conn.Address(), uint32(m.Source), depositNonce, resourceIdStr, types.NewIntFromBigInt(&amount), reciver.String())
		txBts, err := w.conn.client.ConstructAndSignTx(voteMsg)
		if err != nil {
			return false
		}
		txHash, err := w.conn.client.BroadcastTx(txBts)
		err = w.checkAndReSend(txHash, txBts, "voteproposal", err)
		if err != nil {
			return false
		}
		return true

	default:
		w.log.Warn("message type unsupported", "type", m.Type)
		return false
	}
}

func (h *writer) checkAndReSend(txHashStr string, txBts []byte, typeStr string, err error) error {
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "signature repeated"):
			h.log.Info("no need send, already submit signature", "txHash", txHashStr, "type", typeStr)
			return nil
		}
		return err
	} else {
		retry := BlockRetryLimit
		for {
			if retry <= 0 {
				return fmt.Errorf("checkAndSend broadcast tx reach retry limit, tx hash: %s", txHashStr)
			}
			//check on chain
			res, err := h.conn.client.QueryTxByHash(txHashStr)
			if err != nil || res.Empty() || res.Code != 0 {
				h.log.Warn(fmt.Sprintf(
					"checkAndSend QueryTxByHash failed. will rebroadcast after %f second",
					BlockRetryInterval.Seconds()),
					"tx hash", txHashStr,
					"err or res.empty", err)

				//broadcast if not on chain
				_, err = h.conn.client.BroadcastTx(txBts)
				if err != nil && err != errType.ErrTxInMempoolCache {
					h.log.Warn("checkAndSend BroadcastTx failed  will retry", "failed info", err)
				}
				time.Sleep(BlockRetryInterval)
				retry--
				continue
			}
			break
		}
	}
	h.log.Info("checkAndSend success", "txHash", txHashStr, "type", typeStr)
	return nil
}
