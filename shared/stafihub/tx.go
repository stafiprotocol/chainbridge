package stafihub

import (
	"fmt"

	clientTx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xAuthClient "github.com/cosmos/cosmos-sdk/x/auth/client"
	xBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/spf13/cobra"
	"github.com/stafihub/rtoken-relay-core/common/core"
)

func (c *Client) SingleTransferTo(toAddr types.AccAddress, amount types.Coins) error {
	done := core.UseSdkConfigContext(AccountPrefix)
	defer done()
	msg := xBankTypes.NewMsgSend(c.clientCtx.GetFromAddress(), toAddr, amount)
	cmd := cobra.Command{}
	return clientTx.GenerateOrBroadcastTxCLI(c.clientCtx, cmd.Flags(), msg)
}

func (c *Client) BroadcastTx(tx []byte) (string, error) {
	done := core.UseSdkConfigContext(AccountPrefix)
	defer done()
	res, err := c.clientCtx.BroadcastTx(tx)
	if err != nil {
		return "", err
	}
	if res.Code != 0 {
		return "", fmt.Errorf("broadcast err with res.code: %d", res.Code)
	}
	return res.TxHash, nil
}

func (c *Client) ConstructAndSignTx(msgs ...types.Msg) ([]byte, error) {
	account, err := c.GetAccount()
	if err != nil {
		return nil, err
	}
	done := core.UseSdkConfigContext(AccountPrefix)
	defer done()

	cmd := cobra.Command{}
	txf := clientTx.NewFactoryCLI(c.clientCtx, cmd.Flags())
	txf = txf.WithSequence(account.GetSequence()).
		WithAccountNumber(account.GetAccountNumber()).
		WithSignMode(signing.SignMode_SIGN_MODE_DIRECT). //multi sig need this mod
		WithGasAdjustment(1.5).
		WithGas(0).
		WithGasPrices(c.gasPrice).
		WithSimulateAndExecute(true)

	//auto cal gas
	_, adjusted, err := clientTx.CalculateGas(c.clientCtx, txf, msgs...)
	if err != nil {
		return nil, fmt.Errorf("clientTx.CalculateGas failed: %s", err)
	}
	txf = txf.WithGas(adjusted)

	txBuilderRaw, err := clientTx.BuildUnsignedTx(txf, msgs...)
	if err != nil {
		return nil, fmt.Errorf("clientTx.BuildUnsignedTx faild: %s", err)
	}

	err = xAuthClient.SignTx(txf, c.clientCtx, c.clientCtx.GetFromName(), txBuilderRaw, true, true)
	if err != nil {
		return nil, fmt.Errorf("xAuthClient.SignTx failed: %s", err)
	}

	txBytes, err := c.clientCtx.TxConfig.TxEncoder()(txBuilderRaw.GetTx())
	if err != nil {
		return nil, fmt.Errorf("TxConfig.TxEncoder failed: %s", err)
	}
	return txBytes, nil
}
