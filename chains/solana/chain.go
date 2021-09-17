package solana

import (
	"context"
	"errors"
	"strconv"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge/utils/blockstore"
	"github.com/stafiprotocol/chainbridge/utils/core"
	"github.com/stafiprotocol/chainbridge/utils/msg"
)

var TerminatedError = errors.New("terminated")

type Chain struct {
	cfg      *core.ChainConfig // The config of the chain
	conn     *Connection
	listener *listener // The listener of this chain
	writer   *writer   // The writer of the chain
	stop     chan<- int
}

func InitializeChain(cfg *core.ChainConfig, logger log15.Logger, sysErr chan<- error) (*Chain, error) {
	stop := make(chan int)
	conn, err := NewConnection(cfg, logger, stop)
	if err != nil {
		return nil, err
	}

	// Attempt to load latest block
	bs, err := blockstore.NewBlockstore(cfg.BlockstorePath, cfg.Id, conn.poolClient.FeeAccount.PublicKey.ToBase58())
	if err != nil {
		return nil, err
	}

	startSignature := cfg.Opts["startSignature"]
	startSignatureBlock := 0
	bsStartSignature, err := bs.TryLoadLatestSignature()
	if err != nil {
		return nil, err
	}
	bsStartSignatureBlock := 0

	if len(startSignature) != 0 {
		tx, err := conn.GetQueryClient().GetConfirmedTransaction(context.Background(), startSignature)
		if err != nil {
			return nil, err
		}
		startSignatureBlock = int(tx.Slot)
	}
	if len(bsStartSignature) != 0 {
		tx, err := conn.GetQueryClient().GetConfirmedTransaction(context.Background(), bsStartSignature)
		if err != nil {
			return nil, err
		}
		bsStartSignatureBlock = int(tx.Slot)
	}
	useStartSignature := startSignature
	if startSignatureBlock < bsStartSignatureBlock {
		useStartSignature = bsStartSignature
	}

	// Setup listener & writer
	l := NewListener(cfg.Name, conn, cfg.Id, useStartSignature, bs, logger, stop, sysErr)
	w := NewWriter(conn, logger, stop, sysErr)
	return &Chain{cfg: cfg, conn: conn, listener: l, writer: w, stop: stop}, nil
}

func (c *Chain) SetRouter(r *core.Router) {
	r.Listen(c.cfg.Id, c.writer)
	c.listener.setRouter(r)
	c.writer.setRouter(r)
}

func (c *Chain) Start() error {
	err := c.listener.start()
	if err != nil {
		return err
	}

	err = c.writer.start()
	if err != nil {
		return err
	}

	c.writer.log.Debug("Successfully started chain")
	return nil
}

func (c *Chain) Id() msg.ChainId {
	return c.cfg.Id
}

func (c *Chain) Name() string {
	return c.cfg.Name
}

// Stop signals to any running routines to exit
func (c *Chain) Stop() {
	close(c.stop)
}

func parseStartBlock(cfg *core.ChainConfig) uint64 {
	if blk, ok := cfg.Opts["startBlock"]; ok {
		res, err := strconv.ParseUint(blk, 10, 32)
		if err != nil {
			panic(err)
		}
		return res
	}
	return 0
}

// checkBlockstore queries the blockstore for the latest known block. If the latest block is
// greater than startBlock, then the latest block is returned, otherwise startBlock is.
func checkBlockstore(bs *blockstore.Blockstore, startBlock uint64) (uint64, error) {
	latestBlock, err := bs.TryLoadLatestBlock()
	if err != nil {
		return 0, err
	}

	if latestBlock.Uint64() > startBlock {
		return latestBlock.Uint64(), nil
	} else {
		return startBlock, nil
	}
}
