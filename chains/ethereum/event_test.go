package ethereum

import (
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethutils "github.com/stafiprotocol/chainbridge/shared/ethereum"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/chainbridge/utils/msg"
)

var (
	rfisContract   = common.HexToAddress("0xc372e985fda306cfe0e903657de808cf757f536f")
	userAddr       = "0xBd39f5936969828eD9315220659cD11129071814"
	userCommonAddr = common.HexToAddress(userAddr)
)

func TestDepositEvent(t *testing.T) {
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kpI, err := keystore.KeypairFromAddress(userAddr, keystore.EthChain, keystorePath, false)
	if err != nil {
		panic(err)
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	client, err := ethutils.NewClient(ropstenEndpoint, kp)
	if err != nil {
		panic(err)
	}

	destId := msg.ChainId(1)

	nonce, err := ethutils.GetDepositNonce(client, BridgeContract, destId)
	if err != nil {
		panic(err)
	}
	fmt.Println(nonce)

	token, err := ethutils.Erc20GetResourceId(client, Erc20Handler, nativeResouceId)
	if err != nil {
		panic(err)
	}
	fmt.Println(token.Hex())

	allowance, err := ethutils.Erc20GetAllowance(client, Erc20Contract, userCommonAddr, Erc20Handler)
	if err != nil {
		panic(err)
	}
	fmt.Println(allowance)

	bal, err := ethutils.Erc20GetBalance(client, Erc20Contract, userCommonAddr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("currentBal=%+v\n", bal)
	amount := big.NewInt(0).Mul(big.NewInt(50), big.NewInt(1000000000000000000))
	err = ethutils.Erc20Approve(client, Erc20Contract, Erc20Handler, amount)
	if err != nil {
		panic(err)
	}

	nonce++
	err = ethutils.UpdateNonce(client)
	if err != nil {
		panic(err)
	}

	depositAmount := big.NewInt(0).Mul(big.NewInt(15), big.NewInt(1000000000000000000))
	recepient := hexutil.MustDecode("0x26db25c52b007221331a844e5335e59874e45b03e81c3d76ff007377c2c17965")
	err = ethutils.CreateErc20Deposit(client, BridgeContract, destId, recepient, depositAmount, nativeResouceId)
	if err != nil {
		panic(err)
	}

}

func TestRfisDepositEvent(t *testing.T) {
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kpI, err := keystore.KeypairFromAddress(userAddr, keystore.EthChain, keystorePath, false)
	if err != nil {
		panic(err)
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	client, err := ethutils.NewClient(ropstenEndpoint, kp)
	if err != nil {
		panic(err)
	}

	destId := msg.ChainId(1)

	nonce, err := ethutils.GetDepositNonce(client, BridgeContract, destId)
	if err != nil {
		panic(err)
	}
	fmt.Println(nonce)

	token, err := ethutils.Erc20GetResourceId(client, Erc20Handler, rfisResourceId)
	if err != nil {
		panic(err)
	}
	fmt.Println(token.Hex())

	allowance, err := ethutils.Erc20GetAllowance(client, rfisContract, userCommonAddr, Erc20Handler)
	if err != nil {
		panic(err)
	}
	fmt.Println(allowance)

	nonce++
	err = ethutils.UpdateNonce(client)
	if err != nil {
		panic(err)
	}

	bal, err := ethutils.Erc20GetBalance(client, rfisContract, userCommonAddr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("currentBal=%+v\n", bal)
	//amount := big.NewInt(0).Mul(big.NewInt(100), big.NewInt(1000000000000000000))
	//err = ethutils.Erc20Approve(client, rfisContract, erc20Handler, amount)
	//if err != nil {
	//	panic(err)
	//}
	//
	//<-time.After(60*time.Second)
	//
	//allowance, err = ethutils.Erc20GetAllowance(client, rfisContract, userCommonAddr, erc20Handler)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(allowance)
	//
	//depositAmount := big.NewInt(0).Mul(big.NewInt(15), big.NewInt(1000000000000000000))
	//recepient := hexutil.MustDecode("0x3673009bdb664a3f3b6d9f69c9dd37fc0473551a249aa48542408b016ec62b2e")
	//err = ethutils.CreateErc20Deposit(client, bridgeContract, destId, recepient, depositAmount, rfisResourceId)
	//if err != nil {
	//	panic(err)
	//}

}

func TestBigInt(t *testing.T) {
	to := big.NewInt(200)
	Inc(to)
	fmt.Println(to)
}

func Inc(a *big.Int) {
	a = a.Add(a, big.NewInt(1))
}
