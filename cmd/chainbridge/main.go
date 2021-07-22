// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package main

import (
	"errors"
	"os"
	"strconv"

	log "github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge/chains/ethereum"
	"github.com/stafiprotocol/chainbridge/chains/substrate"
	"github.com/stafiprotocol/chainbridge/config"
	"github.com/stafiprotocol/chainbridge/utils/core"
	"github.com/stafiprotocol/chainbridge/utils/msg"
	"github.com/urfave/cli/v2"
)

var app = cli.NewApp()

var cliFlags = []cli.Flag{
	config.ConfigFileFlag,
	config.VerbosityFlag,
	config.KeystorePathFlag,
	config.BlockstorePathFlag,
	config.FreshStartFlag,
	config.LatestBlockFlag,
	config.MetricsFlag,
	config.MetricsPort,
}

var generateFlags = []cli.Flag{
	config.KeystorePathFlag,
}

var accountCommand = cli.Command{
	Name:  "accounts",
	Usage: "manage bridge keystore",
	Description: "The accounts command is used to manage the bridge keystore.\n" +
		"\tMake sure the keystore dir is exist before generating\n" +
		"\tTo generate a substrate keystore: chainbridge accounts gensub\n" +
		"\tTo generate a ethereum keystore: chainbridge accounts geneth\n" +
		"\tTo list keys: chainbridge accounts list",
	Subcommands: []*cli.Command{
		{
			Action: wrapHandler(handleGenerateSubCmd),
			Name:   "gensub",
			Usage:  "generate subsrate keystore",
			Flags:  generateFlags,
			Description: "The generate subcommand is used to generate the substrate keystore.\n",
		},
		{
			Action: wrapHandler(handleGenerateEthCmd),
			Name:   "geneth",
			Usage:  "generate ethereum keystore",
			Flags:  generateFlags,
			Description: "The generate subcommand is used to generate the ethereum keystore.\n",
		},
		{
			Action:      wrapHandler(handleListCmd),
			Name:        "list",
			Usage:       "list bridge keystore",
			Description: "The list subcommand is used to list all of the bridge keystores.\n",
		},
	},
}

// init initializes CLI
func init() {
	app.Action = run
	app.Copyright = "Copyright 2020 Stafi Protocol Authors"
	app.Name = "chainbridge"
	app.Usage = "ChainBridge"
	app.Authors = []*cli.Author{{Name: "Stafi Protocol 2020"}}
	app.Version = "0.0.1"
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		&accountCommand,
	}

	app.Flags = append(app.Flags, cliFlags...)
}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func startLogger(ctx *cli.Context) error {
	logger := log.Root()
	var lvl log.Lvl

	if lvlToInt, err := strconv.Atoi(ctx.String(config.VerbosityFlag.Name)); err == nil {
		lvl = log.Lvl(lvlToInt)
	} else if lvl, err = log.LvlFromString(ctx.String(config.VerbosityFlag.Name)); err != nil {
		return err
	}

	logger.SetHandler(log.MultiHandler(
		log.LvlFilterHandler(
			lvl,
			log.StreamHandler(os.Stdout, log.LogfmtFormat())),
		log.Must.FileHandler("bridge_log.json", log.JsonFormat()),
		log.LvlFilterHandler(
			log.LvlError,
			log.Must.FileHandler("bridge_log_errors.json", log.JsonFormat()))))

	return nil
}

func run(ctx *cli.Context) error {
	err := startLogger(ctx)
	if err != nil {
		return err
	}

	log.Info("Starting ChainBridge...")

	cfg, err := config.GetConfig(ctx)
	if err != nil {
		return err
	}

	// Used to signal core shutdown due to fatal error
	sysErr := make(chan error)
	c := core.NewCore(sysErr)

	for _, chain := range cfg.Chains {
		chainId, err := strconv.Atoi(chain.Id)
		if err != nil {
			return err
		}
		chainConfig := &core.ChainConfig{
			Name:           chain.Name,
			Id:             msg.ChainId(chainId),
			Endpoint:       chain.Endpoint,
			From:           chain.From,
			KeystorePath:   cfg.KeystorePath,
			Insecure:       false,
			BlockstorePath: cfg.BlockStorePath,
			FreshStart:     ctx.Bool(config.FreshStartFlag.Name),
			LatestBlock:    ctx.Bool(config.LatestBlockFlag.Name),
			Opts:           chain.Opts,
			Symbols:        chain.Symbols,
		}
		var newChain core.Chain

		logger := log.Root().New("chain", chainConfig.Name)

		if chain.Type == "ethereum" {
			newChain, err = ethereum.InitializeChain(chainConfig, logger, sysErr)
		} else if chain.Type == "substrate" {
			newChain, err = substrate.InitializeChain(chainConfig, logger, sysErr)
		} else {
			return errors.New("unrecognized Chain Type")
		}

		if err != nil {
			return err
		}
		c.AddChain(newChain)

	}

	c.Start()

	return nil
}
