package ethereum

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/chainbridge/utils/msg"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	ethutils "github.com/stafiprotocol/chainbridge/shared/ethereum"
)

var (
	Erc20Contract  = common.HexToAddress("0x64591e3f2dbf46cdfb398a0d9ba81f41b7cbd449")
	Erc20Handler   = common.HexToAddress("0xeab816f88fe0ebae4971e33f3a21e34ff695791b")
	BridgeContract = common.HexToAddress("0x57e7c280a3828bf9a5356d7c926fcd555cf0bdc8")

	mintAmount = big.NewInt(10000000)

	//resourceId should be found through substrate constants
	resouceId = msg.ResourceIdFromSlice(hexutil.MustDecode("0x000000000000000000000000000000a9e0095b8965c01e6a09c97938f3860901"))
	adminAddr = "0x3aab5AE578FA45744aFe8224DdA506cFE67c508b"
	endpoint  = "wss://ropsten.infura.io/ws/v3/a325d28f7dda49ec9190c8cb4b7f90b2" //ws

	//fill this with keystore path, for eg: /home/stafi/chainbridge/keys
	keystorePath = "/home/stafi/chainbridge/keys"
)

func TestSetUp(t *testing.T) {
	//the password which used to encrypt keystore file, remove it after this test pass
	password := ""
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

	ethutils.Erc20Approve(client, Erc20Contract, Erc20Handler, mintAmount)
	ethutils.Erc20AddMinter(client, Erc20Contract, Erc20Handler)
	ethutils.RegisterResource(client, BridgeContract, Erc20Handler, resouceId, Erc20Contract)
	ethutils.SetBurnable(client, BridgeContract, Erc20Handler, Erc20Contract)
}
