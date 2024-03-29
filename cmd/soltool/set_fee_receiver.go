package main

import (
	"context"
	"fmt"
	"time"

	"github.com/stafiprotocol/chainbridge/shared/solana/vault"
	"github.com/stafiprotocol/solana-go-sdk/bridgeprog"
	solClient "github.com/stafiprotocol/solana-go-sdk/client"
	solCommon "github.com/stafiprotocol/solana-go-sdk/common"
	solTypes "github.com/stafiprotocol/solana-go-sdk/types"
	"github.com/urfave/cli/v2"
)

func setFeeReceiverAction(ctx *cli.Context) error {
	path := ctx.String(configFlag.Name)
	pc := PoolAccounts{}
	err := loadConfig(path, &pc)
	if err != nil {
		return err
	}
	if len(pc.FeeReceiverAccount) == 0 {
		return fmt.Errorf("feeReceiverAccount empty")
	}
	fmt.Printf("\naccounts info:\n %+v\n", pc)
	v, err := vault.NewVaultFromWalletFile(pc.KeystorePath)
	if err != nil {
		return err
	}
	boxer, err := vault.SecretBoxerForType(v.SecretBoxWrap)
	if err != nil {
		return fmt.Errorf("secret boxer: %w", err)
	}

	if err := v.Open(boxer); err != nil {
		return fmt.Errorf("opening: %w", err)
	}

	privKeyMap := make(map[string]vault.PrivateKey)
	for _, privKey := range v.KeyBag {
		privKeyMap[privKey.PublicKey().String()] = privKey
	}

	FeeAccount := solTypes.AccountFromPrivateKeyBytes(privKeyMap[pc.FeeAccount])
	BridgeAccount := solTypes.AccountFromPrivateKeyBytes(privKeyMap[pc.BridgeAccount])
	AdminAccount := solTypes.AccountFromPrivateKeyBytes(privKeyMap[pc.AdminAccountPubkey])
	BridgeProgramId := solCommon.PublicKeyFromString(pc.BridgeProgramId)

	owners := make([]solCommon.PublicKey, 0)
	owners = append(owners, FeeAccount.PublicKey)
	for _, account := range pc.OtherFeeAccountPubkey {
		a := solTypes.AccountFromPrivateKeyBytes(privKeyMap[account])
		owners = append(owners, a.PublicKey)
	}
	if len(owners) < int(pc.Threshold) {
		return fmt.Errorf("owner len < threshold")
	}

	//start inter with solana chain
	c := solClient.NewClient([]string{pc.Endpoint})
	//check if exist
	bridgeInfo, err := c.GetBridgeAccountInfo(context.Background(), BridgeAccount.PublicKey.ToBase58())
	if err != nil {
		return fmt.Errorf("\nbridge account not exist:\n %+v", bridgeInfo)
	}

	res, err := c.GetLatestBlockhash(context.Background(), solClient.GetLatestBlockhashConfig{
		Commitment: solClient.CommitmentConfirmed,
	})
	if err != nil {
		return err
	}

	fmt.Printf("bridgeAccount %s\n", BridgeAccount.PublicKey.ToBase58())
	fmt.Printf("owners %+v\n", owners)
	fmt.Printf("supportChainIds %+v\n", pc.SupportChainIds)
	fmt.Printf("feeReceiver %+v\n", pc.FeeReceiverAccount)

	//create bridge account
	rawTx, err := solTypes.CreateRawTransaction(solTypes.CreateRawTransactionParam{
		Instructions: []solTypes.Instruction{
			bridgeprog.SetFeeReceiver(
				BridgeProgramId,
				BridgeAccount.PublicKey,
				AdminAccount.PublicKey,
				solCommon.PublicKeyFromString(pc.FeeReceiverAccount),
			),
		},
		Signers:         []solTypes.Account{FeeAccount, AdminAccount},
		FeePayer:        FeeAccount.PublicKey,
		RecentBlockHash: res.Blockhash,
	})
	if err != nil {
		return fmt.Errorf("generate tx error, err: %v", err)
	}
	txHash, err := c.SendRawTransaction(context.Background(), rawTx)
	if err != nil {
		return fmt.Errorf("send tx error, err: %v", err)
	}
	fmt.Println("SetFeeReceiver txHash:", txHash)
	time.Sleep(time.Second * 2)

	for i := 0; i < 40; i++ {
		time.Sleep(3 * time.Second)
		bridgeInfo, err := c.GetBridgeAccountInfo(context.Background(), BridgeAccount.PublicKey.ToBase58())
		if err != nil {
			fmt.Println("GetBridgeAccountInfo failed will retry ...", err)
			continue
		}

		fmt.Println("SetFeeReceiver success")
		fmt.Printf("bridge account:\n %+v", bridgeInfo)
		return nil
	}
	fmt.Println("sorry SetFeeReceiver failed")
	return nil
}
