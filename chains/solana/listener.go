package solana

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ChainSafe/log15"
	"github.com/decred/base58"
	borsh "github.com/near/borsh-go"
	"github.com/stafiprotocol/chainbridge/chains"
	"github.com/stafiprotocol/chainbridge/utils/blockstore"
	"github.com/stafiprotocol/chainbridge/utils/msg"
	"github.com/stafiprotocol/solana-go-sdk/bridgeprog"
	solClient "github.com/stafiprotocol/solana-go-sdk/client"
)

var (
	ErrFatalPolling     = errors.New("listener block polling failed")
	eventTickerInterval = time.Second * 15
)

// listen event or block update from solana
type listener struct {
	name           string
	chainId        msg.ChainId
	conn           *Connection
	router         chains.Router
	startSignature string
	blockstore     blockstore.Blockstorer
	log            log15.Logger
	stop           <-chan int
	sysErr         chan<- error
}

func NewListener(name string, conn *Connection, chainId msg.ChainId, startSignature string, bs blockstore.Blockstorer, log log15.Logger, stop <-chan int, sysErr chan<- error) *listener {
	return &listener{
		name:           name,
		chainId:        chainId,
		conn:           conn,
		log:            log,
		stop:           stop,
		startSignature: startSignature,
		blockstore:     bs,
		sysErr:         sysErr,
	}
}

func (l *listener) setRouter(r chains.Router) {
	l.router = r
}

func (l *listener) start() error {

	go func() {
		err := l.pollBlocks()
		if err != nil {
			l.log.Error("Polling blocks failed", "err", err)
			panic(err)
		}
	}()

	return nil
}

func (l *listener) pollBlocks() error {
	l.log.Info("Polling Blocks...")
	eventTicker := time.NewTicker(eventTickerInterval)
	defer eventTicker.Stop()

	for {
		select {
		case <-l.stop:
			return errors.New("polling terminated")
		case <-eventTicker.C:
			// Parse out events
			err := l.getDepositEventsForBlock(l.startSignature)
			if err != nil {
				l.log.Error("Failed to getDepositEventsForBlock", "sig", l.startSignature, "err", err)
				l.sysErr <- ErrFatalPolling
				return err
			}

		}
	}
}

func (l *listener) getDepositEventsForBlock(untilSignature string) error {
	rpcClient := l.conn.queryClient
	bridgeProgramId := l.conn.poolClient.BridgeProgramId.ToBase58()
	bridgeAccount := l.conn.poolClient.BridgeAccountPubkey.ToBase58()

	signatures, err := rpcClient.GetSignaturesForAddress(
		context.Background(),
		bridgeProgramId,
		solClient.GetSignaturesForAddressConfig{
			Until:      untilSignature,
			Commitment: solClient.CommitmentFinalized,
		})

	if err != nil {
		return fmt.Errorf("rpcClient.GetConfirmedSignaturesForAddress err: %s", err.Error())
	}

	for i := len(signatures) - 1; i >= 0; i-- {
		usesig := signatures[i].Signature
		tx, err := rpcClient.GetTransaction(context.Background(), usesig, solClient.GetTransactionWithLimitConfig{
			Commitment:                     solClient.CommitmentFinalized,
			MaxSupportedTransactionVersion: &solClient.DefaultMaxSupportedTransactionVersion,
		})
		if err != nil {
			return fmt.Errorf("rpcClient.GetConfirmedTransaction err: %s", err.Error())
		}
		//skip failed tx
		if tx.Meta.Err != nil {
			err := l.storeDealedSig(usesig)
			if err != nil {
				return err
			}
			continue
		}
		//skip zero instruction
		if len(tx.Transaction.Message.Instructions) == 0 {
			err := l.storeDealedSig(usesig)
			if err != nil {
				return err
			}
			continue
		}
		for _, instruct := range tx.Transaction.Message.Instructions {

			accountKeys := tx.Transaction.Message.AccountKeys
			programIdIndex := instruct.ProgramIDIndex
			if len(accountKeys) <= int(programIdIndex) {
				return fmt.Errorf("accounts or programIdIndex err, %v", tx)
			}

			//skip if it doesn't call  bridge program
			if !strings.EqualFold(accountKeys[programIdIndex], bridgeProgramId) {
				continue
			}

			// check instruction data
			if len(instruct.Data) == 0 {
				continue
			}

			dataBts := base58.Decode(instruct.Data)
			if len(dataBts) < 8 {
				continue
			}
			// skip if it doesn't call transferOut func
			if !bytes.Equal(dataBts[:8], bridgeprog.InstructionTransferOut[:]) {
				l.log.Warn("call func is not transferOut", "tx", tx)
				continue
			}
			// check bridge account
			if len(instruct.Accounts) == 0 {
				continue
			}

			if !strings.EqualFold(accountKeys[instruct.Accounts[0]], bridgeAccount) {
				l.log.Warn("bridge account not equal", "tx", tx)
				continue
			}

			for _, logMessage := range tx.Meta.LogMessages {
				if strings.HasPrefix(logMessage, bridgeprog.EventTransferOutPrefix) {
					l.log.Info("find log", "log", logMessage, "signature", usesig)
					use_log := strings.TrimPrefix(logMessage, bridgeprog.ProgramLogPrefix)
					logBts, err := base64.StdEncoding.DecodeString(use_log)
					if err != nil {
						return err
					}
					if len(logBts) <= 8 {
						return fmt.Errorf("event pase length err")
					}

					eventTransferOut := EventTransferOut{}
					err = borsh.Deserialize(&eventTransferOut, logBts[8:])
					if err != nil {
						return err
					}
					m := msg.NewFungibleTransfer(
						l.chainId,
						msg.ChainId(eventTransferOut.DestChainId),
						msg.Nonce(eventTransferOut.DepositNonce),
						new(big.Int).SetUint64(eventTransferOut.Amount),
						eventTransferOut.ResourceId,
						eventTransferOut.Receiver,
					)
					l.log.Info("send fungibletransfer msg", "msg", m)
					err = l.router.Send(m)
					if err != nil {
						l.log.Error("router send error: failed to route message", "err", err)
						return err
					}
				}

			}
		}
		err = l.storeDealedSig(usesig)
		if err != nil {
			return err
		}
	}
	return nil
}

// save new signature to storage
// Write to block store. Not a critical operation, no need to retry
func (l *listener) storeDealedSig(sig string) error {
	err := l.blockstore.StoreSignature(sig)
	if err != nil {
		l.log.Error("Failed to write latest signature to blockstore", "sig", sig, "err", err)
		return err
	}
	l.startSignature = sig
	return nil
}
