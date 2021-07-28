// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package ethereum

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stafiprotocol/chainbridge/utils/core"
	"github.com/stafiprotocol/chainbridge/utils/msg"
)

const (
	DefaultGasLimit = 1000000
	DefaultGasPrice = 300000000000

	EthChainId = msg.ChainId(2)
	//BscChainId = msg.ChainId(3)
)

var ZeroAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")

// Config encapsulates all necessary parameters in ethereum compatible forms
type Config struct {
	name                 string      // Human-readable chain name
	id                   msg.ChainId // ChainID
	endpoint             string      // url for rpc endpoint
	from                 string      // address of key to use
	keystorePath         string      // Location of keyfiles
	blockstorePath       string
	freshStart           bool // Disables loading from blockstore at start
	bridgeContract       common.Address
	erc20HandlerContract common.Address
	gasLimit             *big.Int
	maxGasPrice          *big.Int
	http                 bool // Config for type of connection
	startBlock           *big.Int
	etherscanUrl         string
}

// parseChainConfig uses a core.ChainConfig to construct a corresponding Config
func ParseChainConfig(chainCfg *core.ChainConfig) (*Config, error) {
	config := &Config{
		name:                 chainCfg.Name,
		id:                   chainCfg.Id,
		endpoint:             chainCfg.Endpoint,
		from:                 chainCfg.From,
		keystorePath:         chainCfg.KeystorePath,
		blockstorePath:       chainCfg.BlockstorePath,
		freshStart:           chainCfg.FreshStart,
		bridgeContract:       ZeroAddress,
		erc20HandlerContract: ZeroAddress,
		gasLimit:             big.NewInt(DefaultGasLimit),
		maxGasPrice:          big.NewInt(DefaultGasPrice),
		http:                 false,
		startBlock:           big.NewInt(0),
	}

	if contract, ok := chainCfg.Opts["bridge"]; ok && contract != "" {
		config.bridgeContract = common.HexToAddress(contract)
		delete(chainCfg.Opts, "bridge")
	} else {
		return nil, fmt.Errorf("must provide opts.bridge field for ethereum config")
	}

	config.erc20HandlerContract = common.HexToAddress(chainCfg.Opts["erc20Handler"])
	delete(chainCfg.Opts, "erc20Handler")

	if gasPrice, ok := chainCfg.Opts["maxGasPrice"]; ok {
		price := big.NewInt(0)
		_, pass := price.SetString(gasPrice, 10)
		if pass {
			config.maxGasPrice = price
			delete(chainCfg.Opts, "maxGasPrice")
		} else {
			return nil, errors.New("unable to parse max gas price")
		}
	}

	if gasLimit, ok := chainCfg.Opts["gasLimit"]; ok {
		limit := big.NewInt(0)
		_, pass := limit.SetString(gasLimit, 10)
		if pass {
			config.gasLimit = limit
			delete(chainCfg.Opts, "gasLimit")
		} else {
			return nil, errors.New("unable to parse gas limit")
		}
	}

	if HTTP, ok := chainCfg.Opts["http"]; ok && HTTP == "true" {
		config.http = true
		delete(chainCfg.Opts, "http")
	} else if HTTP, ok := chainCfg.Opts["http"]; ok && HTTP == "false" {
		config.http = false
		delete(chainCfg.Opts, "http")
	}

	if startBlock, ok := chainCfg.Opts["startBlock"]; ok && startBlock != "" {
		block := big.NewInt(0)
		_, pass := block.SetString(startBlock, 10)
		if pass {
			config.startBlock = block
			delete(chainCfg.Opts, "startBlock")
		} else {
			return nil, errors.New("unable to parse start block")
		}
	}

	if len(chainCfg.Opts) != 0 {
		return nil, fmt.Errorf("unknown Opts Encountered: %#v", chainCfg.Opts)
	}

	return config, nil
}

func (c *Config) From() string {
	return c.from
}

func (c *Config) KeystorePath() string {
	return c.keystorePath
}

func (c *Config) BlockstorePath() string {
	return c.blockstorePath
}

func (c *Config) BridgeContract() common.Address {
	return c.bridgeContract
}

func (c *Config) Erc20HandlerContract() common.Address {
	return c.erc20HandlerContract
}

func (c *Config) StartBlock() *big.Int {
	return c.startBlock
}

func (c *Config) SetStartBlock(blk *big.Int) {
	c.startBlock = blk
}

func (c *Config) ChainId() msg.ChainId {
	return c.id
}

func (c *Config) FreshStart() bool {
	return c.freshStart
}




