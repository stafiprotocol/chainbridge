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
		&mapResourceIdCommand,
		&getBridgeInfoCommand,
		&setSupportChainIdCommand,
		&setFeeReceiverCommand,
		&setFeeAmountCommand,
		&setMintAuthorityCommand,
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

var mapResourceIdCommand = cli.Command{
	Name:        "mapResourceId",
	Usage:       "map resourceId to mint account",
	Description: "The mapResourceId command is used to update resourceId to mint account",
	Action:      mapResourceIdAction,
	Flags:       cliFlags,
}

var setSupportChainIdCommand = cli.Command{
	Name:        "setChainId",
	Usage:       "set support chain ids",
	Description: "The setChainId command is used to set support chain ids",
	Action:      setSupportChainIdAction,
	Flags:       cliFlags,
}

var setFeeReceiverCommand = cli.Command{
	Name:        "setFeeReceiver",
	Usage:       "set fee receiver",
	Description: "The setFeeReceiver command is used to set fee receiver",
	Action:      setFeeReceiverAction,
	Flags:       cliFlags,
}

var setMintAuthorityCommand = cli.Command{
	Name:        "setMintAuthority",
	Usage:       "set mint authority",
	Description: "The setMintAuthority command is used to set mint authority",
	Action:      setMintAuthority,
	Flags:       cliFlags,
}

var setFeeAmountCommand = cli.Command{
	Name:        "setFeeAmount",
	Usage:       "set fee amount",
	Description: "The setFeeAmount command is used to set fee amounts",
	Action:      setFeeAmountAction,
	Flags:       cliFlags,
}

var getBridgeInfoCommand = cli.Command{
	Name:        "getBridgeInfo",
	Usage:       "get bridge info",
	Description: "The getBridgeInfo command is used to get bridge account info",
	Action:      getBridgeInfoAction,
	Flags:       cliFlags,
}

type PoolAccounts struct {
	KeystorePath          string            `json:"keystorePath"`
	FeeAccount            string            `json:"feeAccount"`         //need private key
	BridgeAccount         string            `json:"bridgeAccount"`      //need private key
	AdminAccountPubkey    string            `json:"adminAccountPubkey"` //need private key
	BridgeProgramId       string            `json:"bridgeProgramId"`
	Endpoint              string            `json:"endpoint"`
	OtherFeeAccountPubkey []string          `json:"otherFeeAccountPubkey"`
	Threshold             uint64            `json:"threshold"`
	SupportChainIds       []uint8           `json:"supportChainIds"`
	ResourceIdToMint      map[string]string `json:"resourceIdToMint"`
	FeeReceiverAccount    string            `json:"feeReceiverAccount"` //need private key
	FeeAmounts            map[string]uint64 `json:"feeAmounts"`
}

type PoolAccountsForSetMint struct {
	KeystorePath       string `json:"keystorePath"`
	FeeAccount         string `json:"feeAccount"`         //need private key
	BridgeAccount      string `json:"bridgeAccount"`      //need private key
	AdminAccountPubkey string `json:"adminAccountPubkey"` //need private key
	BridgeProgramId    string `json:"bridgeProgramId"`
	Endpoint           string `json:"endpoint"`

	NewMintAuthority string `json:"newMintAuthority"`
	BridgeSigner     string `json:"bridgeSigner"`
	RSolMint         string `json:"rsolMint"`
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

func loadConfigForMint(file string, config *PoolAccountsForSetMint) (err error) {
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
