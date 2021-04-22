package ethereum

import (
	"context"
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
	rfisContract        = common.HexToAddress("0xc372e985fda306cfe0e903657de808cf757f536f")
	erc20Contract       = common.HexToAddress("0x64591e3f2dbf46cdfb398a0d9ba81f41b7cbd449")
	erc20Handler        = common.HexToAddress("0xeab816f88fe0ebae4971e33f3a21e34ff695791b")
	bridgeContract      = common.HexToAddress("0x57e7c280a3828bf9a5356d7c926fcd555cf0bdc8")
	userAddr            = "0xBd39f5936969828eD9315220659cD11129071814"
	userCommonAddr      = common.HexToAddress(userAddr)
	endPoint            = "wss://ropsten.infura.io/ws/v3/a325d28f7dda49ec9190c8cb4b7f90b2"
	currentKeystorePath = "/Users/fwj/Go/stafi/chainbridge/keys"
	rId                 = msg.ResourceIdFromSlice(hexutil.MustDecode("0x000000000000000000000000000000a9e0095b8965c01e6a09c97938f3860901"))
)

func TestEvent(t *testing.T) {
	//the password which used to encrypt keystore file, remove it after this test pass
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kpI, err := keystore.KeypairFromAddress(userAddr, keystore.EthChain, currentKeystorePath, false)
	if err != nil {
		panic(err)
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	client, err := ethutils.NewClient(endPoint, kp)
	if err != nil {
		panic(err)
	}

	latestBlock := big.NewInt(11373068)
	to := big.NewInt(200)
	to = to.Add(to, latestBlock)
	query := buildQuery(bridgeContract, ethutils.ProposalEvent, latestBlock, to)
	evts, err := client.Client.FilterLogs(context.Background(), query)
	if err != nil {
		t.Fatal(err)
	}
	//fmt.Printf("%+v\n", evts)

	for _, evt := range evts {
		sourceId := evt.Topics[1].Big().Uint64()
		depositNonce := evt.Topics[2].Big().Uint64()
		status := evt.Topics[3].Big().Uint64()

		fmt.Printf("sourceId = %d, depositNonce=%d, status=%d\n", sourceId, depositNonce, status)
	}
}

func TestDepositEvent(t *testing.T) {
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kpI, err := keystore.KeypairFromAddress(userAddr, keystore.EthChain, currentKeystorePath, false)
	if err != nil {
		panic(err)
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	client, err := ethutils.NewClient(endPoint, kp)
	if err != nil {
		panic(err)
	}

	destId := msg.ChainId(1)

	nonce, err := ethutils.GetDepositNonce(client, bridgeContract, destId)
	if err != nil {
		panic(err)
	}
	fmt.Println(nonce)

	token, err := ethutils.Erc20GetResourceId(client, erc20Handler, wraResourceId)
	if err != nil {
		panic(err)
	}
	fmt.Println(token.Hex())

	allowance, err := ethutils.Erc20GetAllowance(client, WRAContract, userCommonAddr, erc20Handler)
	if err != nil {
		panic(err)
	}
	fmt.Println(allowance)

	nonce++
	err = ethutils.UpdateNonce(client)
	if err != nil {
		panic(err)
	}

	amount := big.NewInt(0).Mul(big.NewInt(15), big.NewInt(1000000000000000000))
	recepient := hexutil.MustDecode("0x26db25c52b007221331a844e5335e59874e45b03e81c3d76ff007377c2c17965")
	err = ethutils.CreateErc20Deposit(client, bridgeContract, destId, recepient, amount, wraResourceId)
	if err != nil {
		panic(err)
	}
	//bal, err := ethutils.Erc20GetBalance(client, WRAContract, userCommonAddr)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("currentBal=%+v\n", bal)
	//amount := big.NewInt(0).Mul(big.NewInt(50), big.NewInt(1000000000000000000))
	//err = ethutils.Erc20Approve(client, WRAContract, erc20Handler, amount)
	//if err != nil {
	//	panic(err)
	//}
}

func TestRfisDepositEvent(t *testing.T) {
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kpI, err := keystore.KeypairFromAddress(userAddr, keystore.EthChain, currentKeystorePath, false)
	if err != nil {
		panic(err)
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	client, err := ethutils.NewClient(endPoint, kp)
	if err != nil {
		panic(err)
	}

	destId := msg.ChainId(1)

	nonce, err := ethutils.GetDepositNonce(client, bridgeContract, destId)
	if err != nil {
		panic(err)
	}
	fmt.Println(nonce)

	token, err := ethutils.Erc20GetResourceId(client, erc20Handler, rfisResourceId)
	if err != nil {
		panic(err)
	}
	fmt.Println(token.Hex())

	allowance, err := ethutils.Erc20GetAllowance(client, rfisContract, userCommonAddr, erc20Handler)
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
