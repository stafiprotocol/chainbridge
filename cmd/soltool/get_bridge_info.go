package main

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	solClient "github.com/stafiprotocol/solana-go-sdk/client"
	"github.com/stafiprotocol/solana-go-sdk/common"
	"github.com/urfave/cli/v2"
)

func getBridgeInfoAction(ctx *cli.Context) error {
	path := ctx.String(configFlag.Name)
	pc := PoolAccounts{}
	err := loadConfig(path, &pc)
	if err != nil {
		return err
	}
	fmt.Printf("\nconfig info:\n %+v\n", pc)

	//start inter with solana chain
	c := solClient.NewClient([]string{pc.Endpoint})
	//check if exist
	bridgeInfo, err := c.GetBridgeAccountInfo(context.Background(), pc.BridgeAccount)
	if err != nil {
		return fmt.Errorf("\nbridge account not exist:\n %+v", bridgeInfo)
	}

	fmt.Printf("\nbridgeAccountInfo: \n")
	fmt.Printf("admin: %s\n", common.PublicKeyFromBytes(bridgeInfo.Admin[:]).ToBase58())
	fmt.Printf("feeReceiver: %s\n", common.PublicKeyFromBytes(bridgeInfo.FeeReceiver[:]).ToBase58())
	fmt.Printf("nonce: %d\n", bridgeInfo.Nonce)
	fmt.Printf("ownerSetSequence: %d\n", bridgeInfo.OwnerSetSeqno)
	fmt.Printf("threshold: %d\n", bridgeInfo.Threshold)
	fmt.Printf("owners: %s\n", strOwners(bridgeInfo.Owners))
	fmt.Printf("supportChainIds: %v\n", bridgeInfo.SupportChainIds)
	fmt.Printf("depositCount: %s\n", strDepositCount(bridgeInfo.DepositCount))
	fmt.Printf("feeAmounts: %s\n", strDepositCount(bridgeInfo.FeeAmount))
	fmt.Printf("resourceIdToMint: %s\n", strResourceToMint(bridgeInfo.ResourceIdToMint))

	return nil
}

func strOwners(owners [][32]uint8) string {
	ret := "\n"
	for _, v := range owners {
		ret += common.PublicKeyFromBytes(v[:]).ToBase58()
		ret += "\n"
	}
	return ret
}

func strDepositCount(m map[uint8]uint64) string {
	ret := "\n"
	for key, v := range m {
		ret += fmt.Sprintf("%d", key)
		ret += " : "
		ret += fmt.Sprintf("%d", v)
		ret += "\n"
	}
	return ret
}

func strResourceToMint(m map[[32]uint8]common.PublicKey) string {
	ret := "\n"
	for k, v := range m {
		ret += hexutil.Encode(k[:])
		ret += " : "
		ret += v.ToBase58()
		ret += "\n"
	}
	return ret
}
