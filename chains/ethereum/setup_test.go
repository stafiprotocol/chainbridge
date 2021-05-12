package ethereum

import (
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
	RfisContract   = common.HexToAddress("0xc372e985fda306cfe0e903657de808cf757f536f")
	WRAContract    = common.HexToAddress("0x54cdC5A3fc5C803AbB2EFafd93CDa01b62347b69")
	Erc20Contract  = common.HexToAddress("0x64591e3f2dbf46cdfb398a0d9ba81f41b7cbd449")
	Erc20Handler   = common.HexToAddress("0x16BEDc0b5cA0B42826F3eB3259304fBF83D414B6")
	BridgeContract = common.HexToAddress("0x129e704e435e726103E5f96909B936381F1e8FB9")

	mintAmount, _ = big.NewInt(0).SetString("100000000000000000000000000000", 10)

	//resourceId should be found through substrate constants
	nativeResouceId = msg.ResourceIdFromSlice(hexutil.MustDecode("0x000000000000000000000000000000a9e0095b8965c01e6a09c97938f3860901"))
	rfisResourceId  = msg.ResourceIdFromSlice(hexutil.MustDecode("0x000000000000000000000000000000b9e0095b8965c01e6a09c97938f3860901"))
	wraResourceId   = msg.ResourceIdFromSlice(hexutil.MustDecode("0x000000000000000000000000000000cdd5da9607788c8f4a14c30fcbef84d101"))
	adminAddr       = "0xBd39f5936969828eD9315220659cD11129071814"
	oldAdminAddr       = "0x3aab5AE578FA45744aFe8224DdA506cFE67c508b"
	endpoint        = "wss://ropsten.infura.io/ws/v3/a325d28f7dda49ec9190c8cb4b7f90b2" //ws

	//fill this with keystore path, for eg: /home/stafi/chainbridge/keys
	keystorePath = "/Users/fwj/Go/stafi/chainbridge/keys"
)

func TestSetUp(t *testing.T) {
	//the password which used to encrypt keystore file, remove it after this test pass
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kpI, err := keystore.KeypairFromAddress(oldAdminAddr, keystore.EthChain, keystorePath, false)
	if err != nil {
		panic(err)
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	client, err := ethutils.NewClient(endpoint, kp)
	if err != nil {
		panic(err)
	}

	ethutils.Erc20Approve(client, Erc20Contract, Erc20Handler, mintAmount)
	ethutils.Erc20AddMinter(client, Erc20Contract, Erc20Handler)
	//ethutils.RegisterResource(client, BridgeContract, Erc20Handler, nativeResouceId, Erc20Contract)
	//ethutils.SetBurnable(client, BridgeContract, Erc20Handler, Erc20Contract)
}

func TestSetUpRfis(t *testing.T) {
	//the password which used to encrypt keystore file, remove it after this test pass
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kpI, err := keystore.KeypairFromAddress(adminAddr, keystore.EthChain, keystorePath, false)
	if err != nil {
		panic(err)
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	client, err := ethutils.NewClient(endpoint, kp)
	if err != nil {
		panic(err)
	}

	ethutils.Erc20Approve(client, RfisContract, Erc20Handler, mintAmount)
	ethutils.Erc20AddMinter(client, RfisContract, Erc20Handler)
	ethutils.RegisterResource(client, BridgeContract, Erc20Handler, rfisResourceId, RfisContract)
	ethutils.SetBurnable(client, BridgeContract, Erc20Handler, RfisContract)
}

func TestSetUpWRA(t *testing.T) {
	//the password which used to encrypt keystore file, remove it after this test pass
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kpI, err := keystore.KeypairFromAddress(adminAddr, keystore.EthChain, keystorePath, false)
	if err != nil {
		panic(err)
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	client, err := ethutils.NewClient(endpoint, kp)
	if err != nil {
		panic(err)
	}

	ethutils.Erc20Approve(client, WRAContract, Erc20Handler, mintAmount)
	ethutils.Erc20AddMinter(client, WRAContract, Erc20Handler)
	ethutils.RegisterResource(client, BridgeContract, Erc20Handler, wraResourceId, WRAContract)
	ethutils.SetBurnable(client, BridgeContract, Erc20Handler, WRAContract)
}

func TestResourceId(t *testing.T) {
	t.Log(rfisResourceId.Hex())
}
