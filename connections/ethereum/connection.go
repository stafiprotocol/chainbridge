// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package ethereum

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ChainSafe/log15"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcommon "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/chainbridge/utils/msg"
)

var (
	BlockRetryInterval = time.Second * 5
	lowExtraGasPrice   = big.NewInt(5e9)
	highExtraGasPrice  = big.NewInt(10e9)
	standGasPrice      = big.NewInt(20e9)
)

var Gwei5 = big.NewInt(5e9)
var Gwei10 = big.NewInt(10e9)
var Gwei20 = big.NewInt(20e9)
var MinGasTipCap = big.NewInt(3e9)

type Connection struct {
	endpoint    string
	chainId     msg.ChainId
	http        bool
	kp          *secp256k1.Keypair
	gasLimit    *big.Int
	maxGasPrice *big.Int
	conn        *ethclient.Client
	// signer    ethtypes.Signer
	opts     *bind.TransactOpts
	callOpts *bind.CallOpts
	nonce    uint64
	optsLock sync.Mutex
	log      log15.Logger
	stop     chan int // All routines should exit when this channel is closed
}

// NewConnection returns an uninitialized connection, must call Connection.Connect() before using.
func NewConnection(cfg *Config, kp *secp256k1.Keypair, log log15.Logger) *Connection {
	return &Connection{
		endpoint:    cfg.endpoint,
		chainId:     cfg.id,
		http:        cfg.http,
		kp:          kp,
		gasLimit:    cfg.gasLimit,
		maxGasPrice: cfg.maxGasPrice,
		log:         log,
		stop:        make(chan int),
	}
}

// Connect starts the ethereum WS connection
func (c *Connection) Connect() error {
	c.log.Info("Connecting to ethereum chain...", "url", c.endpoint)
	var rpcClient *rpc.Client
	var err error
	// Start http or ws client
	if c.http {
		rpcClient, err = rpc.DialHTTP(c.endpoint)
	} else {
		rpcClient, err = rpc.DialWebsocket(context.Background(), c.endpoint, "/ws")
	}
	if err != nil {
		return err
	}
	c.conn = ethclient.NewClient(rpcClient)

	var chainId *big.Int
	retry := 0
	for {
		if retry > 50 {
			return fmt.Errorf("get chainId err: %s", err)
		}
		chainId, err = c.conn.ChainID(context.Background())
		if err != nil {
			c.log.Warn("get chainId err", "err", err)
			retry++
			time.Sleep(time.Second * 3)
			continue
		}
		break
	}

	// Construct tx opts, call opts, and nonce mechanism
	opts, _, err := c.newTransactOpts(big.NewInt(0), c.gasLimit, c.maxGasPrice, chainId)
	if err != nil {
		return err
	}
	c.opts = opts
	c.nonce = 0
	c.callOpts = &bind.CallOpts{From: c.kp.CommonAddress()}
	return nil
}

// newTransactOpts builds the TransactOpts for the connection's keypair.
func (c *Connection) newTransactOpts(value, gasLimit, gasPrice, chainId *big.Int) (*bind.TransactOpts, uint64, error) {
	privateKey := c.kp.PrivateKey()
	address := ethcrypto.PubkeyToAddress(privateKey.PublicKey)

	nonce, err := c.conn.PendingNonceAt(context.Background(), address)
	if err != nil {
		return nil, 0, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		return nil, 0, err
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = value
	auth.GasLimit = uint64(gasLimit.Int64())
	auth.Context = context.Background()

	return auth, nonce, nil
}

func (c *Connection) Keypair() *secp256k1.Keypair {
	return c.kp
}

func (c *Connection) Client() *ethclient.Client {
	return c.conn
}

func (c *Connection) Opts() *bind.TransactOpts {
	return c.opts
}

func (c *Connection) CallOpts() *bind.CallOpts {
	return c.callOpts
}

func (c *Connection) SafeEstimateGas(ctx context.Context) (*big.Int, error) {
	gasPrice, err := c.conn.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	if c.chainId == EthChainId {
		if gasPrice.Cmp(standGasPrice) == 1 {
			gasPrice = gasPrice.Add(gasPrice, highExtraGasPrice)
		} else {
			gasPrice = gasPrice.Add(gasPrice, lowExtraGasPrice)
		}
	}

	if gasPrice.Cmp(c.maxGasPrice) > 0 {
		gasPrice = c.maxGasPrice
	}

	return gasPrice, nil
}

// return suggest gastipcap gasfeecap
func (c *Connection) SafeEstimateFee(ctx context.Context) (*big.Int, *big.Int, error) {
	marketGasTipCap, err := c.conn.SuggestGasTipCap(ctx)
	if err != nil {
		return nil, nil, err
	}
	marketGasFeeCap, err := c.conn.SuggestGasPrice(ctx)
	if err != nil {
		return nil, nil, err
	}

	if marketGasFeeCap.Cmp(Gwei20) < 0 {
		marketGasFeeCap = new(big.Int).Add(marketGasFeeCap, Gwei5)
	} else {
		marketGasFeeCap = new(big.Int).Add(marketGasFeeCap, Gwei10)
	}

	gasFeeCap, _ := new(big.Float).Mul(new(big.Float).SetInt(marketGasFeeCap), big.NewFloat(1.2)).Int(nil)
	gasTipCap, _ := new(big.Float).Mul(new(big.Float).SetInt(marketGasTipCap), big.NewFloat(1.2)).Int(nil)

	if gasFeeCap.Cmp(c.maxGasPrice) > 0 {
		gasFeeCap = c.maxGasPrice
	}
	if gasTipCap.Cmp(MinGasTipCap) < 0 {
		gasTipCap = MinGasTipCap
	}

	return gasTipCap, gasFeeCap, nil
}

// LockAndUpdateOpts acquires a lock on the opts before updating the nonce
// and gas price.
func (c *Connection) LockAndUpdateOpts() error {
	c.optsLock.Lock()

	var err error
	defer func() {
		if err != nil {
			c.optsLock.Unlock()
		}
	}()
	gasTipCap, gasFeeCap, err := c.SafeEstimateFee(context.Background())
	if err != nil {
		return err
	}
	c.opts.GasTipCap = gasTipCap
	c.opts.GasFeeCap = gasFeeCap

	nonce, err := c.conn.NonceAt(context.Background(), c.opts.From, nil)
	if err != nil {
		return err
	}
	c.opts.Nonce.SetUint64(nonce)
	return nil
}

type GasRes struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  struct {
		LastBlock       string `json:"LastBlock"`
		SafeGasPrice    string `json:"SafeGasPrice"`
		ProposeGasPrice string `json:"ProposeGasPrice"`
		FastGasPrice    string `json:"FastGasPrice"`
	} `json:"result"`
}

// LockAndUpdateOpts acquires a lock on the opts before updating the nonce
// and gas price.
// func (c *Connection) LockAndUpdateOpts() error {
// 	c.optsLock.Lock()

// 	gasPrice, err := c.SafeEstimateGas(context.TODO())
// 	if err != nil {
// 		return err
// 	}

// 	c.opts.GasPrice = gasPrice

// 	nonce, err := c.conn.PendingNonceAt(context.Background(), c.opts.From)
// 	if err != nil {
// 		c.optsLock.Unlock()
// 		return err
// 	}
// 	c.opts.Nonce.SetUint64(nonce)
// 	return nil
// }

func (c *Connection) UnlockOpts() {
	c.optsLock.Unlock()
}

// LatestBlock returns the latest block from the current chain
func (c *Connection) LatestBlock() (*big.Int, error) {
	header, err := c.conn.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return header.Number, nil
}

// EnsureHasBytecode asserts if contract code exists at the specified address
func (c *Connection) EnsureHasBytecode(addr ethcommon.Address) error {
	code, err := c.conn.CodeAt(context.Background(), addr, nil)
	if err != nil {
		return err
	}

	if len(code) == 0 {
		return fmt.Errorf("no bytecode found at %s", addr.Hex())
	}
	return nil
}

// WaitForBlock will poll for the block number until the current block is equal or greater than
func (c *Connection) WaitForBlock(block *big.Int) error {
	blk := big.NewInt(3)
	blk = blk.Add(blk, block)
	for {
		select {
		case <-c.stop:
			return errors.New("connection terminated")
		default:
			currBlock, err := c.LatestBlock()
			if err != nil {
				return err
			}

			// Greater than target
			if currBlock.Cmp(blk) >= 0 {
				return nil
			}
			c.log.Trace("Block not ready, waiting", "target", block, "current", currBlock)
			time.Sleep(BlockRetryInterval)
			continue
		}
	}
}

// Close terminates the client connection and stops any running routines
func (c *Connection) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
	close(c.stop)
}
