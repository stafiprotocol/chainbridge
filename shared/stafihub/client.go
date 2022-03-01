package stafihub

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptoTypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types"
	xAuthTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	rpcClient "github.com/tendermint/tendermint/rpc/client"
	rpcHttp "github.com/tendermint/tendermint/rpc/client/http"
)

const AccountPrefix = "stafi"

//stafihub client
type Client struct {
	clientCtx     client.Context
	rpcClient     rpcClient.Client
	gasPrice      string
	denom         string
	accountNumber uint64
}

func NewClient(k keyring.Keyring, fromName, gasPrice, endPoint string) (*Client, error) {
	encodingConfig := MakeEncodingConfig()
	var retClient *Client
	if len(fromName) != 0 {
		info, err := k.Key(fromName)
		if err != nil {
			return nil, fmt.Errorf("keyring get address from name:%s err: %s", fromName, err)
		}
		rpcClient, err := rpcHttp.New(endPoint, "/websocket")
		if err != nil {
			return nil, err
		}

		initClientCtx := client.Context{}.
			WithCodec(encodingConfig.Marshaler).
			WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
			WithTxConfig(encodingConfig.TxConfig).
			WithLegacyAmino(encodingConfig.Amino).
			WithInput(os.Stdin).
			WithAccountRetriever(xAuthTypes.AccountRetriever{}).
			WithBroadcastMode(flags.BroadcastBlock).
			WithClient(rpcClient).
			WithSkipConfirmation(true).         //skip password confirm
			WithFromName(fromName).             //keyBase need FromName to find key info
			WithFromAddress(info.GetAddress()). //accountRetriever need FromAddress
			WithKeyring(k)

		retClient = &Client{
			clientCtx: initClientCtx,
			rpcClient: rpcClient,
		}
		chaindId, err := retClient.GetChainId()
		if err != nil {
			return nil, err
		}
		retClient.clientCtx = retClient.clientCtx.WithChainID(chaindId)

		account, err := retClient.GetAccount()
		if err != nil {
			return nil, err
		}

		retClient.accountNumber = account.GetAccountNumber()

		bondedDenom, err := retClient.QueryBondedDenom()
		if err != nil {
			return nil, err
		}

		retClient.setDenom(bondedDenom.Params.BondDenom)
		err = retClient.setGasPrice(gasPrice)
		if err != nil {
			return nil, err
		}
	} else {
		rpcClient, err := rpcHttp.New(endPoint, "/websocket")
		if err != nil {
			return nil, err
		}

		initClientCtx := client.Context{}.
			WithCodec(encodingConfig.Marshaler).
			WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
			WithTxConfig(encodingConfig.TxConfig).
			WithLegacyAmino(encodingConfig.Amino).
			WithInput(os.Stdin).
			WithAccountRetriever(xAuthTypes.AccountRetriever{}).
			WithBroadcastMode(flags.BroadcastBlock).
			WithClient(rpcClient).
			WithSkipConfirmation(true) //skip password confirm

		retClient = &Client{
			clientCtx: initClientCtx,
			rpcClient: rpcClient,
		}
		bondedDenom, err := retClient.QueryBondedDenom()
		if err != nil {
			return nil, err
		}

		retClient.setDenom(bondedDenom.Params.BondDenom)
	}
	return retClient, nil
}

//update clientCtx.FromName and clientCtx.FromAddress
func (c *Client) SetFromName(fromName string) error {
	info, err := c.clientCtx.Keyring.Key(fromName)
	if err != nil {
		return fmt.Errorf("keyring get address from fromKeyname err: %s", err)
	}

	c.clientCtx = c.clientCtx.WithFromName(fromName).WithFromAddress(info.GetAddress())

	account, err := c.GetAccount()
	if err != nil {
		return err
	}
	c.accountNumber = account.GetAccountNumber()
	return nil
}

func (c *Client) GetFromName() string {
	return c.clientCtx.FromName
}

func (c *Client) GetFromAddress() types.AccAddress {
	return c.clientCtx.FromAddress
}

func (c *Client) setGasPrice(gasPrice string) error {
	_, err := types.ParseDecCoins(gasPrice)
	if err != nil {
		return err
	}
	c.gasPrice = gasPrice
	return nil
}

func (c *Client) setDenom(denom string) {
	c.denom = denom
}

func (c *Client) GetDenom() string {
	return c.denom
}

func (c *Client) GetTxConfig() client.TxConfig {
	return c.clientCtx.TxConfig
}

func (c *Client) GetLegacyAmino() *codec.LegacyAmino {
	return c.clientCtx.LegacyAmino
}

func (c *Client) Sign(fromName string, toBeSigned []byte) ([]byte, cryptoTypes.PubKey, error) {
	return c.clientCtx.Keyring.Sign(fromName, toBeSigned)
}

func (c *Client) Ctx() client.Context {
	return c.clientCtx
}
