package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/ChainSafe/log15"

	"github.com/urfave/cli/v2"
)

var app = cli.NewApp()

var configFlag = &cli.StringFlag{
	Name:  "config",
	Usage: "json configuration file",
	Value: "./sol_create_bridge.json",
}
var cliFlags = []cli.Flag{
	configFlag,
}

// init initializes CLI
func init() {
	app.Action = run
	app.Copyright = "Copyright 2020 Stafi Protocol Authors"
	app.Name = "solTool"
	app.Usage = "solTool"
	app.Authors = []*cli.Author{{Name: "Stafi Protocol 2020"}}
	app.Version = "0.0.1"
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		&createBridgeCommand,
		&setResourceIdCommand,
	}

}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func run(ctx *cli.Context) error {
	return nil
}

var createBridgeCommand = cli.Command{
	Name:        "createBridge",
	Usage:       "create bridge account",
	Description: "The createBridge command is used to create a bridge account",
	Action:      createBridgeAccountAction,
	Flags:       cliFlags,
}

var setResourceIdCommand = cli.Command{
	Name:        "setResourceId",
	Usage:       "setResourceId",
	Description: "The setResourceId command is used to update resourceId to mint account",
	Action:      setResourceIdAction,
	Flags:       cliFlags,
}

type PoolAccounts struct {
	KeystorePath          string            `json:"keystorePath"`
	FeeAccount            string            `json:"feeAccount"`    //need private key
	BridgeAccount         string            `json:"bridgeAccount"` //need private key
	AdminAccountPubkey    string            `json:"adminAccountPubkey"`
	BridgeProgramId       string            `json:"bridgeProgramId"`
	Endpoint              string            `json:"endpoint"`
	OtherFeeAccountPubkey []string          `json:"otherFeeAccountPubkey"`
	Threshold             uint64            `json:"threshold"`
	SupportChainIds       []uint8           `json:"supportChainIds"`
	ResourceIdToMint      map[string]string `json:"resourceIdToMint"`
}

func loadConfig(file string, config *PoolAccounts) (err error) {
	ext := filepath.Ext(file)
	fp, err := filepath.Abs(file)
	if err != nil {
		return err
	}

	log.Debug("Loading configuration", "path", filepath.Clean(fp))

	f, err := os.Open(filepath.Clean(fp))
	if err != nil {
		return err
	}
	defer func() {
		err = f.Close()
	}()

	if ext != ".json" {
		return fmt.Errorf("unrecognized extention: %s", ext)
	}
	return json.NewDecoder(f).Decode(&config)
}
