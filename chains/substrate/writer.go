// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package substrate

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ChainSafe/log15"
	"github.com/shopspring/decimal"
	"github.com/stafiprotocol/chainbridge/config"
	"github.com/stafiprotocol/chainbridge/utils/msg"
	"github.com/stafiprotocol/go-substrate-rpc-client/types"
)

const (
	msgLimit = 4096
)

var TerminatedError = errors.New("terminated")

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
					panic(result)
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

	var prop *proposal
	var err error

	// Construct the proposal
	switch m.Type {
	case msg.FungibleTransfer:
		prop, err = w.createFungibleProposal(m)
	default:
		w.log.Warn("unrecognized message type received")
		return false
	}

	if err != nil {
		w.sysErr <- fmt.Errorf("construct proposal Error: %s", err)
		return false
	}

	w.log.Info("ResolveMessage prop", "nonce", prop.DepositNonce, "source",
		prop.SourceId, "resource", fmt.Sprintf("%x", prop.ResourceId), "method", prop.Method)

	for i := 0; i < BlockRetryLimit; i++ {
		// Ensure we only submit a vote if status of the proposal is Active
		valid, reason, err := w.proposalValid(prop)
		w.log.Info("ResolveMessage proposalValid", "valid", valid, "reason", reason)
		if err != nil {
			w.log.Error("Failed to assert proposal state", "err", err)
			time.Sleep(BlockRetryInterval)
			continue
		}

		if !valid {
			w.log.Debug("Ignoring proposal", "reason", reason)
			return true
		}

		w.log.Info("Acknowledging proposal on chain")
		ext, err := w.conn.gc.NewUnsignedExtrinsic(config.AcknowledgeProposal, prop.DepositNonce, prop.SourceId, prop.ResourceId, prop.Call)
		if err != nil {
			w.log.Error("Acknowledging NewUnsignedExtrinsic met err")
			return false
		}
		err = w.conn.gc.SignAndSubmitTx(ext)
		if err != nil {
			if err.Error() == TerminatedError.Error() {
				w.log.Error("Acknowledging proposal met TerminatedError")
				return false
			}
			w.log.Error("Acknowledging proposal error", "err", err)
			time.Sleep(BlockRetryInterval)
			continue
		}
		return true
	}
	return false
}

func (w *writer) createFungibleProposal(m msg.Message) (*proposal, error) {
	bigAmt := big.NewInt(0).SetBytes(m.Payload[0].([]byte))
	//should not have 0x prefix and length must 64
	resourceIdStr := strings.ToLower(m.ResourceId.Hex())
	if len(resourceIdStr) != 64 {
		return nil, fmt.Errorf("resourceId  length  must be 64")
	}

	d, ok := w.decimals[resourceIdStr]
	if !ok {
		d, ok = w.decimals[decimalDefault]
		if !ok {
			return nil, fmt.Errorf("failed to get decimal")
		}
	}

	amount := types.NewU128(*decimal.NewFromBigInt(bigAmt, 0).Div(d).BigInt())
	recipient := types.NewAccountID(m.Payload[1].([]byte))
	depositNonce := types.U64(m.DepositNonce)
	method, err := w.resolveResourceId(m.ResourceId)
	if err != nil {
		return nil, err
	}

	meta, err := w.conn.gc.GetLatestMetadata()
	if err != nil {
		return nil, err
	}

	call, err := types.NewCall(
		meta,
		method,
		recipient,
		amount,
		types.NewBytes32(m.ResourceId),
	)
	if err != nil {
		return nil, err
	}

	return &proposal{
		DepositNonce: depositNonce,
		Call:         call,
		SourceId:     types.U8(m.Source),
		ResourceId:   types.NewBytes32(m.ResourceId),
		Method:       method,
	}, nil
}

func (w *writer) resolveResourceId(id [32]byte) (string, error) {
	var res []byte
	exist, err := w.conn.QueryStorage(config.BridgeCommon, "Resources", id[:], nil, &res)
	if err != nil {
		return "", err
	}

	if !exist {
		return "", fmt.Errorf("resource %x not found on chain", id)
	}

	return string(res), nil
}

func (w *writer) proposalValid(prop *proposal) (bool, string, error) {
	var voteRes voteState
	srcId, err := types.EncodeToBytes(prop.SourceId)
	if err != nil {
		return false, "", err
	}

	propBz, err := prop.encode()
	if err != nil {
		return false, "", err
	}

	exists, err := w.conn.QueryStorage(config.BridgeCommon, "Votes", srcId, propBz, &voteRes)
	if err != nil {
		return false, "", err
	}

	if !exists {
		return true, "", nil
	}

	if voteRes.Status != VoteStatusActive {
		return false, fmt.Sprintf("CurrentVoteStatus: %s", voteRes.Status), nil
	}

	if containsVote(voteRes.Voted, types.NewAccountID(w.conn.key.PublicKey)) {
		return false, "already voted", nil
	}

	return true, "", nil
}

func containsVote(votes []types.AccountID, voter types.AccountID) bool {
	for _, v := range votes {
		if bytes.Equal(v[:], voter[:]) {
			return true
		}
	}
	return false
}
