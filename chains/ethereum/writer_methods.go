// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package ethereum

import (
	"context"
	"errors"
	"math/big"
	"time"

	eth "github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	utils "github.com/stafiprotocol/chainbridge/shared/ethereum"
	"github.com/stafiprotocol/chainbridge/utils/msg"
)

// Number of blocks to wait for an finalization event
const ExecuteBlockWatchLimit = 200
const ExecuteBlockQueryLimit = 3

// Time between retrying a failed tx
const TxRetryInterval = time.Second * 2

// Maximum number of tx retries before exiting
const TxRetryLimit = 10

var ErrNonceTooLow = errors.New("nonce too low")
var ErrTxUnderpriced = errors.New("replacement transaction underpriced")
var ErrFatalTx = errors.New("submission of transaction failed")
var ErrFatalQuery = errors.New("query of chain state failed")

type ExecuteErrEnum string

const (
	ExecuteErrStop          = ExecuteErrEnum("stop")
	ExecuteErrRetryExceeded = ExecuteErrEnum("retry_exceeded")
	ExecuteErrFetchLog      = ExecuteErrEnum("fetch_log")
	ExecuteErrBlockLimit    = ExecuteErrEnum("block_limit")
	ExecuteErrEventNotFound = ExecuteErrEnum("event_not_found")
	ExecuteErrNormal        = ExecuteErrEnum("execued")
)

// proposalIsComplete returns true if the proposal state is either Passed, Transferred or Cancelled
func (w *writer) proposalIsComplete(srcId msg.ChainId, nonce msg.Nonce, dataHash [32]byte) bool {
	prop, err := w.bridgeContract.GetProposal(w.conn.CallOpts(), uint8(srcId), uint64(nonce), dataHash)
	if err != nil {
		w.log.Error("Failed to check proposal existence", "err", err)
		return false
	}
	return prop.Status == PassedStatus || prop.Status == TransferredStatus || prop.Status == CancelledStatus
}

// proposalIsComplete returns true if the proposal state is Transferred or Cancelled
func (w *writer) proposalIsFinalized(srcId msg.ChainId, nonce msg.Nonce, dataHash [32]byte) bool {
	prop, err := w.bridgeContract.GetProposal(w.conn.CallOpts(), uint8(srcId), uint64(nonce), dataHash)
	if err != nil {
		w.log.Error("Failed to check proposal existence", "err", err)
		return false
	}
	return prop.Status == TransferredStatus || prop.Status == CancelledStatus // Transferred (3)
}

// hasVoted checks if this relayer has already voted
func (w *writer) hasVoted(srcId msg.ChainId, nonce msg.Nonce, dataHash [32]byte) bool {
	hasVoted, err := w.bridgeContract.HasVotedOnProposal(w.conn.CallOpts(), utils.IDAndNonce(srcId, nonce), dataHash, w.conn.Opts().From)
	if err != nil {
		w.log.Error("Failed to check proposal existence", "err", err)
		return false
	}

	return hasVoted
}

func (w *writer) onlyVoted(m msg.Message, dataHash [32]byte) bool {
	prop, err := w.bridgeContract.GetProposal(w.conn.CallOpts(), uint8(m.Source), uint64(m.DepositNonce), dataHash)
	if err != nil {
		w.log.Error("Failed to check proposal existence", "err", err)
		return false
	}

	if prop.Status != PassedStatus {
		return false
	}

	return w.hasVoted(m.Source, m.DepositNonce, dataHash)
}

func (w *writer) shouldVote(m msg.Message, dataHash [32]byte) bool {
	// Check if proposal has passed and skip if Passed or Transferred
	if w.proposalIsComplete(m.Source, m.DepositNonce, dataHash) {
		w.log.Info("Proposal complete, not voting", "src", m.Source, "nonce", m.DepositNonce)
		return false
	}

	// Check if relayer has previously voted
	if w.hasVoted(m.Source, m.DepositNonce, dataHash) {
		w.log.Info("Relayer has already voted, not voting", "src", m.Source, "nonce", m.DepositNonce)
		return false
	}

	return true
}

// createErc20Proposal creates an Erc20 proposal.
// Returns true if the proposal is successfully created or is complete
func (w *writer) createErc20Proposal(m msg.Message, propResult chan<- bool) {
	w.log.Info("Creating erc20 proposal", "src", m.Source, "nonce", m.DepositNonce)

	data := ConstructErc20ProposalData(m.Payload[0].([]byte), m.Payload[1].([]byte))
	dataHash := utils.Hash(append(w.cfg.erc20HandlerContract.Bytes(), data...))

	if w.onlyVoted(m, dataHash) {
		w.executeProposal(m, data, dataHash)
		result := w.proposalIsFinalized(m.Source, m.DepositNonce, dataHash)
		w.log.Info("executeProposal for those onlyVoted", "source", m.Source, "DepositNonce", m.DepositNonce, "result", result)
		propResult <- result
		return
	}

	if !w.shouldVote(m, dataHash) {
		propResult <- true
		return
	}

	// Capture latest block so when know where to watch from
	latestBlock, err := w.conn.LatestBlock()
	if err != nil {
		w.log.Error("Unable to fetch latest block", "err", err)
		propResult <- false
		return
	}

	// watch for execution event
	go func() {
		// copy latestBlock
		from := big.NewInt(0)
		from = from.Add(from, latestBlock)
		result := w.watchThenExecute(m, data, dataHash, latestBlock)
		if result == ExecuteErrNormal {
			propResult <- true
			return
		}

		for i := 0; i < ExecuteBlockQueryLimit; i++ {
			result = w.queryThenExecute(m, data, dataHash, from, latestBlock)
			switch result {
			case ExecuteErrNormal:
				w.log.Info("queryAndExecuteRetrySuccess", "i", i, "from", from, "to", latestBlock)
				propResult <- true
				return
			default:
				w.log.Info("queryAndExecuteRetryFail", "i", i, "from", from, "to", latestBlock, "result", result)
			}
		}
		propResult <- false
	}()

	w.voteProposal(m, dataHash)
}

// watchThenExecute watches for the latest block and executes once the matching finalized event is found
func (w *writer) watchThenExecute(m msg.Message, data []byte, dataHash [32]byte, latestBlock *big.Int) ExecuteErrEnum {
	w.log.Info("Watching for finalization event", "src", m.Source, "nonce", m.DepositNonce, "latestBlock", latestBlock)

	// watching for the latest block, querying and matching the finalized event will be retried up to ExecuteBlockWatchLimit times
	for i := 0; i < ExecuteBlockWatchLimit; i++ {
		select {
		case <-w.stop:
			return ExecuteErrStop
		default:
			// watch for the lastest block, retry up to BlockRetryLimit times
			for waitRetrys := 0; waitRetrys < BlockRetryLimit; waitRetrys++ {
				err := w.conn.WaitForBlock(latestBlock)
				if err != nil {
					w.log.Error("Waiting for block failed", "err", err)
					// Exit if retries exceeded
					if waitRetrys+1 == BlockRetryLimit {
						w.log.Error("Waiting for block retries exceeded, shutting down")
						w.sysErr <- ErrFatalQuery
						return ExecuteErrRetryExceeded
					}
				} else {
					break
				}
			}
			result := w.queryThenExecute(m, data, dataHash, latestBlock, latestBlock)
			switch result {
			case ExecuteErrEventNotFound:
				latestBlock = latestBlock.Add(latestBlock, big.NewInt(1))
			default:
				return result
			}
		}
	}
	w.log.Warn("Block watch limit exceeded, skipping execution", "source", m.Source, "dest", m.Destination, "nonce", m.DepositNonce)
	return ExecuteErrBlockLimit
}

func (w *writer) queryThenExecute(m msg.Message, data []byte, dataHash [32]byte, from, to *big.Int) ExecuteErrEnum {
	// query for logs
	query := buildQuery(w.cfg.bridgeContract, utils.ProposalEvent, from, to)
	evts, err := w.conn.Client().FilterLogs(context.Background(), query)
	if err != nil {
		w.log.Error("Failed to fetch logs", "err", err)
		return ExecuteErrFetchLog
	} else {
		w.log.Debug("queryThenExecute", "from", from, "to", to, "EventNum", len(evts))
	}

	// execute the proposal once we find the matching finalized event
	for _, evt := range evts {
		sourceId := evt.Topics[1].Big().Uint64()
		depositNonce := evt.Topics[2].Big().Uint64()
		status := evt.Topics[3].Big().Uint64()

		if m.Source == msg.ChainId(sourceId) && m.DepositNonce.Big().Uint64() == depositNonce {
			st := uint8(status)
			switch {
			case utils.IsFinalized(st):
				w.executeProposal(m, data, dataHash)
				return ExecuteErrNormal
			case utils.IsExecuted(st):
				w.log.Info("Proposal already executed", "src", m.Source, "depositNonce", m.DepositNonce)
				return ExecuteErrNormal
			default:
				w.log.Trace("event status not satisfied", "src", sourceId, "nonce", depositNonce, "status", st)
			}
		} else {
			w.log.Trace("Ignoring event", "src", sourceId, "nonce", depositNonce, "status", status)
		}
	}
	w.log.Trace("No finalization event found in current block", "from", from, "to", to, "src", m.Source, "nonce", m.DepositNonce)
	return ExecuteErrEventNotFound
}

// voteProposal submits a vote proposal
// a vote proposal will try to be submitted up to the TxRetryLimit times
func (w *writer) voteProposal(m msg.Message, dataHash [32]byte) {
	for i := 0; i < TxRetryLimit; i++ {
		select {
		case <-w.stop:
			return
		default:
			err := w.conn.LockAndUpdateOpts()
			if err != nil {
				w.log.Error("Failed to update tx opts", "err", err)
				continue
			}

			tx, err := w.bridgeContract.VoteProposal(
				w.conn.Opts(),
				uint8(m.Source),
				uint64(m.DepositNonce),
				m.ResourceId,
				dataHash,
			)
			w.conn.UnlockOpts()

			if err == nil {
				w.log.Info("Submitted proposal vote", "tx", tx.Hash(), "src", m.Source, "depositNonce", m.DepositNonce)
				return
			} else if err.Error() == ErrNonceTooLow.Error() || err.Error() == ErrTxUnderpriced.Error() {
				w.log.Debug("Nonce too low, will retry")
				time.Sleep(TxRetryInterval)
			} else {
				w.log.Warn("Voting failed", "source", m.Source, "dest", m.Destination, "depositNonce", m.DepositNonce, "err", err)
				time.Sleep(TxRetryInterval)
			}

			// Verify proposal is still open for voting, otherwise no need to retry
			if w.proposalIsComplete(m.Source, m.DepositNonce, dataHash) {
				w.log.Info("Proposal voting complete on chain", "src", m.Source, "dst", m.Destination, "nonce", m.DepositNonce)
				return
			}
		}
	}
	w.log.Error("Submission of Vote transaction failed", "source", m.Source, "dest", m.Destination, "depositNonce", m.DepositNonce)
	w.sysErr <- ErrFatalTx
}

// executeProposal executes the proposal
func (w *writer) executeProposal(m msg.Message, data []byte, dataHash [32]byte) {
	for i := 0; i < TxRetryLimit; i++ {
		select {
		case <-w.stop:
			return
		default:
			err := w.conn.LockAndUpdateOpts()
			if err != nil {
				w.log.Error("Failed to update nonce", "err", err)
				return
			}

			tx, err := w.bridgeContract.ExecuteProposal(
				w.conn.Opts(),
				uint8(m.Source),
				uint64(m.DepositNonce),
				data,
				m.ResourceId,
			)
			w.conn.UnlockOpts()

			if err == nil {
				w.log.Info("Submitted proposal execution", "tx", tx.Hash(), "src", m.Source, "dst", m.Destination, "nonce", m.DepositNonce)
				return
			} else if err.Error() == ErrNonceTooLow.Error() || err.Error() == ErrTxUnderpriced.Error() {
				w.log.Error("Nonce too low, will retry")
				time.Sleep(TxRetryInterval)
			} else {
				w.log.Warn("Execution failed, proposal may already be complete", "err", err)
				time.Sleep(TxRetryInterval)
			}

			// Verify proposal is still open for execution, tx will fail if we aren't the first to execute,
			// but there is no need to retry
			if w.proposalIsFinalized(m.Source, m.DepositNonce, dataHash) {
				w.log.Info("Proposal finalized on chain", "src", m.Source, "dst", m.Destination, "nonce", m.DepositNonce)
				return
			}
		}
	}
	w.log.Error("Submission of Execute transaction failed", "source", m.Source, "dest", m.Destination, "depositNonce", m.DepositNonce)
	w.sysErr <- ErrFatalTx
}

// buildQuery constructs a query for the bridgeContract by hashing sig to get the event topic
func buildQuery(contract ethcommon.Address, sig utils.EventSig, startBlock *big.Int, endBlock *big.Int) eth.FilterQuery {
	query := eth.FilterQuery{
		FromBlock: startBlock,
		ToBlock:   endBlock,
		Addresses: []ethcommon.Address{contract},
		Topics: [][]ethcommon.Hash{
			{sig.GetTopic()},
		},
	}
	return query
}
