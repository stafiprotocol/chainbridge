package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/chainbridge/shared/solana/vault"
	"github.com/stafiprotocol/solana-go-sdk/bridgeprog"
	solClient "github.com/stafiprotocol/solana-go-sdk/client"
	solCommon "github.com/stafiprotocol/solana-go-sdk/common"
	"github.com/stafiprotocol/solana-go-sdk/sysprog"
	solTypes "github.com/stafiprotocol/solana-go-sdk/types"
	"github.com/urfave/cli/v2"
)

func createBridgeAccountAction(ctx *cli.Context) error {
	path := ctx.String(configFlag.Name)
	pc := PoolAccounts{}
	err := loadConfig(path, &pc)
	if err != nil {
		return err
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
	FeeReceiverAccount := solTypes.AccountFromPrivateKeyBytes(privKeyMap[pc.FeeReceiverAccount])

	BridgeProgramId := solCommon.PublicKeyFromString(pc.BridgeProgramId)

	PdaPubkey, nonce, err := solCommon.FindProgramAddress([][]byte{BridgeAccount.PublicKey.Bytes()}, BridgeProgramId)
	if err != nil {
		return err
	}
	fmt.Println("\nbridgePdaPubkey: ", PdaPubkey.ToBase58())

	owners := make([]solCommon.PublicKey, 0)
	owners = append(owners, FeeAccount.PublicKey)
	for _, account := range pc.OtherFeeAccountPubkey {
		a := solTypes.AccountFromPrivateKeyBytes(privKeyMap[account])
		owners = append(owners, a.PublicKey)
	}
	if len(owners) < int(pc.Threshold) {
		return fmt.Errorf("owner len < threshold")
	}

	resourceIdToMint := make(map[[32]byte]solCommon.PublicKey)
	for key, value := range pc.ResourceIdToMint {
		id, err := hexutil.Decode(key)
		if err != nil {
			return err
		}
		var useId [32]byte
		copy(useId[:], id)
		mintAccountPubkey := solCommon.PublicKeyFromString(value)
		resourceIdToMint[useId] = mintAccountPubkey
	}

	feeAmounts := make(map[uint8]uint64)
	for k, v := range pc.FeeAmounts {
		id, err := strconv.ParseUint(k, 10, 8)
		if err != nil {
			return err
		}
		feeAmounts[uint8(id)] = v
	}

	//start inter with solana chain
	c := solClient.NewClient([]string{pc.Endpoint})
	//check bridge exist
	bridgeInfo, err := c.GetBridgeAccountInfo(context.Background(), BridgeAccount.PublicKey.ToBase58())
	if err == nil {
		return fmt.Errorf("\nbridge account already exist:\n %+v", bridgeInfo)
	}

	res, err := c.GetRecentBlockhash(context.Background())
	if err != nil {
		return err
	}
	bridgeAccountMiniMum, err := c.GetMinimumBalanceForRentExemption(context.Background(), solClient.BridgeAccountLengthDefault)
	if err != nil {
		return err
	}

	fmt.Printf("resourceIdToMint %+v\n", resourceIdToMint)
	fmt.Printf("bridgeAccount %s\n", BridgeAccount.PublicKey.ToBase58())
	fmt.Printf("owners %+v\n", owners)
	fmt.Printf("supportChainIds %+v\n", pc.SupportChainIds)

	//create bridge account
	rawTx, err := solTypes.CreateRawTransaction(solTypes.CreateRawTransactionParam{
		Instructions: []solTypes.Instruction{
			sysprog.CreateAccount(
				FeeAccount.PublicKey,
				BridgeAccount.PublicKey,
				BridgeProgramId,
				bridgeAccountMiniMum*2,
				solClient.BridgeAccountLengthDefault,
			),
			bridgeprog.CreateBridge(
				BridgeProgramId,
				BridgeAccount.PublicKey,
				owners,
				pc.Threshold,
				uint8(nonce),
				pc.SupportChainIds,
				resourceIdToMint,
				AdminAccount.PublicKey,
				FeeReceiverAccount.PublicKey,
				feeAmounts,
			),
		},
		Signers:         []solTypes.Account{FeeAccount, BridgeAccount},
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
	fmt.Println("create bridge account txHash:", txHash)
	time.Sleep(time.Second * 2)

	for i := 0; i < 40; i++ {
		time.Sleep(3 * time.Second)
		bridgeInfo, err := c.GetBridgeAccountInfo(context.Background(), BridgeAccount.PublicKey.ToBase58())
		if err != nil {
			fmt.Println("GetBridgeAccountInfo failed will retry ...", err)
			continue
		}

		fmt.Println("create bridge success")
		fmt.Printf("bridge account:\n %+v", bridgeInfo)
		return nil
	}
	fmt.Println("sorry create bridge failed")
	return nil
}
