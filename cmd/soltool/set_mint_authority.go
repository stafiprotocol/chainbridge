package main

import (
	"context"
	"fmt"

	"github.com/stafiprotocol/chainbridge/shared/solana/vault"
	"github.com/stafiprotocol/solana-go-sdk/bridgeprog"
	solClient "github.com/stafiprotocol/solana-go-sdk/client"
	solCommon "github.com/stafiprotocol/solana-go-sdk/common"
	solTypes "github.com/stafiprotocol/solana-go-sdk/types"
	"github.com/urfave/cli/v2"
)

func setMintAuthority(ctx *cli.Context) error {
	path := ctx.String(configFlag.Name)
	pc := PoolAccountsForSetMint{}
	err := loadConfigForMint(path, &pc)
	if err != nil {
		return err
	}
	if len(pc.NewMintAuthority) == 0 {
		return fmt.Errorf("NewMintAuthority empty")
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

	newMintAuthority := solCommon.PublicKeyFromString(pc.NewMintAuthority)
	rsolMint := solCommon.PublicKeyFromString(pc.RSolMint)
	bridgeSigner := solCommon.PublicKeyFromString(pc.BridgeSigner)

	fmt.Printf("bridgeAccount %s\n", BridgeAccount.PublicKey.ToBase58())
	fmt.Printf("newMintAuthority %+v\n", newMintAuthority.ToBase58())
	fmt.Printf("bridgeSigner %+v\n", bridgeSigner.ToBase58())

	fmt.Println("\ncheck newMintAuthority again, then press enter to continue:")
	var input string
	fmt.Scanln(&input)

	//create bridge account
	rawTx, err := solTypes.CreateRawTransaction(solTypes.CreateRawTransactionParam{
		Instructions: []solTypes.Instruction{
			bridgeprog.SetMintAuthority(
				BridgeProgramId,
				BridgeAccount.PublicKey,
				AdminAccount.PublicKey,
				bridgeSigner,
				rsolMint,
				newMintAuthority,
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
	fmt.Println("SetMintAuthority txHash:", txHash)

	return nil
}
