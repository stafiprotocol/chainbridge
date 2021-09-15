package solana

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ChainSafe/log15"
	borsh "github.com/near/borsh-go"
	"github.com/stafiprotocol/chainbridge/chains"
	"github.com/stafiprotocol/chainbridge/utils/blockstore"
	"github.com/stafiprotocol/chainbridge/utils/msg"
	"github.com/stafiprotocol/solana-go-sdk/bridgeprog"
	solClient "github.com/stafiprotocol/solana-go-sdk/client"
)

var (
	BlockRetryInterval = time.Second * 3
	BlockRetryLimit    = 50
	ErrFatalPolling    = errors.New("listener block polling failed")
	logInterval        = uint64(100)
)

//listen event or block update from solana
type listener struct {
	name       string
	chainId    msg.ChainId
	conn       *Connection
	router     chains.Router
	startBlock uint64
	blockstore blockstore.Blockstorer
	log        log15.Logger
	stop       <-chan int
	sysErr     chan<- error
}

func NewListener(name string, conn *Connection, chainId msg.ChainId, startBlock uint64, bs blockstore.Blockstorer, log log15.Logger, stop <-chan int, sysErr chan<- error) *listener {
	return &listener{
		name:       name,
		chainId:    chainId,
		conn:       conn,
		log:        log,
		stop:       stop,
		startBlock: startBlock,
		blockstore: bs,
		sysErr:     sysErr,
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
	var currentBlock = l.startBlock
	var retry = BlockRetryLimit

	for {
		select {
		case <-l.stop:
			return errors.New("polling terminated")
		default:
			// No more retries, goto next block
			if retry == 0 {
				l.log.Error("Polling failed, retries exceeded")
				l.sysErr <- ErrFatalPolling
				return nil
			}
			rpcClient := l.conn.GetQueryClient()

			latestBlock, err := rpcClient.GetBlockHeight(context.Background(), solClient.GetBlockHeightConfig{solClient.CommitmentFinalized})
			if err != nil {
				l.log.Error("Unable to get latest block", "block", currentBlock, "err", err)
				retry--
				time.Sleep(BlockRetryInterval)
				continue
			}

			if currentBlock%logInterval == 0 {
				l.log.Debug("pollBlocks", "target", currentBlock, "latest", latestBlock)
			}

			// Sleep if the latestFinishBlcok is less than currentBlock
			if currentBlock > latestBlock {
				time.Sleep(BlockRetryInterval)
				continue
			}

			// Parse out events
			err = l.getDepositEventsForBlock(currentBlock)
			if err != nil {
				l.log.Error("Failed to get events for block", "block", currentBlock, "err", err)
				retry--
				continue
			}

			// Write to block store. Not a critical operation, no need to retry
			err = l.blockstore.StoreBlock(big.NewInt(int64(currentBlock)))
			if err != nil {
				l.log.Error("Failed to write latest block to blockstore", "block", currentBlock, "err", err)
			}

			// Goto next block and reset retry counter
			currentBlock++
			retry = BlockRetryLimit
		}
	}
}

func (l *listener) getDepositEventsForBlock(blockNumber uint64) error {
	rpcClient := l.conn.queryClient
	block, err := rpcClient.GetConfirmedBlock(context.Background(), blockNumber)
	if err != nil {
		return err
	}

	for _, tx := range block.Transactions {
		for _, logMessage := range tx.Meta.LogMessages {
			if strings.HasPrefix(logMessage, bridgeprog.EventTransferOutPrefix) {
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
				err = l.router.Send(m)
				if err != nil {
					l.log.Error("subscription error: failed to route message", "err", err)
				}

			}
		}
	}

	return nil
}
