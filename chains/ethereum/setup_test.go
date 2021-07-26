package ethereum

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stafiprotocol/chainbridge/bindings/Bridge"
	ethutils "github.com/stafiprotocol/chainbridge/shared/ethereum"
	"github.com/stafiprotocol/chainbridge/utils/crypto/secp256k1"
	"github.com/stafiprotocol/chainbridge/utils/keystore"
	"github.com/stafiprotocol/chainbridge/utils/msg"
	"math/big"
	"os"
	"testing"
)

var (
	RfisContract   = common.HexToAddress("0xB89F77F28c2D84D32ED8bb2311AFA29A6089559C")
	WRAContract    = common.HexToAddress("0x54cdC5A3fc5C803AbB2EFafd93CDa01b62347b69")
	Erc20Contract  = common.HexToAddress("0x63E8273B4C4E4d47CDF91047907c61665fC606f3")
	Erc20Handler   = common.HexToAddress("0x290f95d650103Bfd765459cDe4B7277749374335")
	BridgeContract = common.HexToAddress("0xe169c92136f45617266789251Ef4EA4604Abc96c")

	mintAmount, _ = big.NewInt(0).SetString("100000000000000000000000000000", 10)

	//resourceId should be found through substrate constants
	nativeResouceId = msg.ResourceIdFromSlice(hexutil.MustDecode("0x000000000000000000000000000000a9e0095b8965c01e6a09c97938f3860901"))
	rfisResourceId  = msg.ResourceIdFromSlice(hexutil.MustDecode("0x000000000000000000000000000000b9e0095b8965c01e6a09c97938f3860901"))
	wraResourceId   = msg.ResourceIdFromSlice(hexutil.MustDecode("0x000000000000000000000000000000cdd5da9607788c8f4a14c30fcbef84d101"))
	adminAddr       = "0xBd39f5936969828eD9315220659cD11129071814"
	oldAdminAddr    = "0x3aab5AE578FA45744aFe8224DdA506cFE67c508b"
	ropstenEndpoint = "wss://ropsten.infura.io/ws/v3/a325d28f7dda49ec9190c8cb4b7f90b2" //ws
	goerliEndpoint  = "wss://goerli.infura.io/ws/v3/a325d28f7dda49ec9190c8cb4b7f90b2"

	//fill this with keystore path, for eg: /home/stafi/chainbridge/keys
	keystorePath = "/Users/fwj/Go/stafi/chainbridge/keys"

	bscBridgeContract = common.HexToAddress("0xb5Dc44e4e680d9a485fCbe1f465eA0F65287FBc0")
	bep20FisContract  = common.HexToAddress("0x2E7EC0D5E15f044C1f7C138c1045D424380540Fe")
	bep20RFisContract = common.HexToAddress("0xEBa705c21d0F7a32b54C381E1a8391d6656844de")
	bscErc20Handler   = common.HexToAddress("0x88D106ffd21e01ccd4a0011ecaBc06DfD7f60319")
	owner             = "0xBca9567A9e8D5F6F58C419d32aF6190F74C880e6"
	BscTestEndpoint   = "https://data-seed-prebsc-2-s3.binance.org:8545/"
)

func TestSetUp(t *testing.T) {
	//the password which used to encrypt keystore file, remove it after this test pass
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kpI, err := keystore.KeypairFromAddress(owner, keystore.EthChain, keystorePath, false)
	if err != nil {
		t.Fatal(err)
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	client, err := ethutils.NewClient(goerliEndpoint, kp)
	if err != nil {
		t.Fatal(err)
	}

	err = ethutils.Erc20Approve(client, Erc20Contract, Erc20Handler, mintAmount)
	if err != nil {
		t.Fatal(err)
	}

	err = ethutils.Erc20AddMinter(client, Erc20Contract, Erc20Handler)
	if err != nil {
		t.Fatal(err)
	}

	rIds := [][32]byte{nativeResouceId}
	tokens := []common.Address{Erc20Contract}

	err = ethutils.SetResourceAndBurnable(client, BridgeContract, Erc20Handler, rIds, tokens, tokens)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetUpRfis(t *testing.T) {
	//the password which used to encrypt keystore file, remove it after this test pass
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kpI, err := keystore.KeypairFromAddress(owner, keystore.EthChain, keystorePath, false)
	if err != nil {
		panic(err)
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	client, err := ethutils.NewClient(goerliEndpoint, kp)
	if err != nil {
		panic(err)
	}

	//err = ethutils.Erc20Approve(client, RfisContract, Erc20Handler, mintAmount)
	//if err != nil {
	//	t.Fatal(err)
	//}

	//err = ethutils.Erc20AddMinter(client, RfisContract, Erc20Handler)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//
	rIds := [][32]byte{rfisResourceId}
	tokens := []common.Address{RfisContract}
	burnableTokens := []common.Address{RfisContract}

	err = ethutils.SetResourceAndBurnable(client, BridgeContract, Erc20Handler, rIds, tokens, burnableTokens)
	if err != nil {
		t.Fatal(err)
	}
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

	client, err := ethutils.NewClient(ropstenEndpoint, kp)
	if err != nil {
		panic(err)
	}

	err = ethutils.Erc20Approve(client, RfisContract, Erc20Handler, mintAmount)
	if err != nil {
		t.Fatal(err)
	}

	err = ethutils.Erc20AddMinter(client, RfisContract, Erc20Handler)
	if err != nil {
		t.Fatal(err)
	}

	rIds := [][32]byte{wraResourceId}
	tokens := []common.Address{WRAContract}
	err = ethutils.SetResourceAndBurnable(client, BridgeContract, Erc20Handler, rIds, tokens, tokens)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetUpBscFis(t *testing.T) {
	//the password which used to encrypt keystore file, remove it after this test pass
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kpI, err := keystore.KeypairFromAddress(owner, keystore.EthChain, keystorePath, false)
	if err != nil {
		panic(err)
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	client, err := ethutils.NewClient(BscTestEndpoint, kp)
	if err != nil {
		panic(err)
	}

	//err = ethutils.Erc20Approve(client, bep20FisContract, bscErc20Handler, mintAmount)
	//if err != nil {
	//	t.Fatal(err)
	//}

	//err = ethutils.Erc20AddMinter(client, bep20FisContract, bscErc20Handler)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//

	rIds := [][32]byte{nativeResouceId}
	tokens := []common.Address{bep20FisContract}

	err = ethutils.SetResourceAndBurnable(client, bscBridgeContract, bscErc20Handler, rIds, tokens, tokens)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSetUpBscRFis(t *testing.T) {
	//the password which used to encrypt keystore file, remove it after this test pass
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kpI, err := keystore.KeypairFromAddress(owner, keystore.EthChain, keystorePath, false)
	if err != nil {
		panic(err)
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	client, err := ethutils.NewClient(BscTestEndpoint, kp)
	if err != nil {
		panic(err)
	}

	//err = ethutils.Erc20Approve(client, bep20RFisContract, bscErc20Handler, mintAmount)
	//if err != nil {
	//	t.Fatal(err)
	//}

	//err = ethutils.Erc20AddMinter(client, bep20RFisContract, bscErc20Handler)
	//if err != nil {
	//	t.Fatal(err)
	//}

	rIds := [][32]byte{rfisResourceId}
	tokens := []common.Address{bep20RFisContract}
	burnableTokens := []common.Address{}

	err = ethutils.SetResourceAndBurnable(client, bscBridgeContract, bscErc20Handler, rIds, tokens, burnableTokens)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGoerliFisDepositEvent(t *testing.T) {
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kpI, err := keystore.KeypairFromAddress(owner, keystore.EthChain, keystorePath, false)
	if err != nil {
		t.Fatal(err)
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	client, err := ethutils.NewClient(goerliEndpoint, kp)
	if err != nil {
		t.Fatal(err)
	}

	/// 1 for stafi, 2 for eth, 3 for bsc
	//destId := msg.ChainId(1)
	//destId := msg.ChainId(2)
	destId := msg.ChainId(3)

	nonce, err := ethutils.GetDepositNonce(client, BridgeContract, destId)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("nonce", nonce)

	token, err := ethutils.Erc20GetResourceId(client, Erc20Handler, nativeResouceId)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("token", token.Hex())

	allowance, err := ethutils.Erc20GetAllowance(client, token, common.HexToAddress(owner), Erc20Handler)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("allowance", allowance)

	nonce++
	err = ethutils.UpdateNonce(client)
	if err != nil {
		t.Fatal(err)
	}

	depositAmount := big.NewInt(0).Mul(big.NewInt(5), big.NewInt(100000000000000000))
	/// Stafi recepient
	//recepient := hexutil.MustDecode("0x3673009bdb664a3f3b6d9f69c9dd37fc0473551a249aa48542408b016ec62b2e")

	/// Eth recepient
	recepient := common.HexToAddress(owner).Bytes()
	err = ethutils.CreateErc20Deposit(client, BridgeContract, destId, recepient, depositAmount, nativeResouceId)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGoerliRFisDepositEvent(t *testing.T) {
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kpI, err := keystore.KeypairFromAddress(owner, keystore.EthChain, keystorePath, false)
	if err != nil {
		t.Fatal(err)
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	client, err := ethutils.NewClient(goerliEndpoint, kp)
	if err != nil {
		t.Fatal(err)
	}

	/// 1 for stafi, 2 for eth, 3 for bsc
	//destId := msg.ChainId(1)
	//destId := msg.ChainId(2)
	destId := msg.ChainId(3)

	nonce, err := ethutils.GetDepositNonce(client, BridgeContract, destId)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("nonce", nonce)

	token, err := ethutils.Erc20GetResourceId(client, Erc20Handler, rfisResourceId)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("token", token.Hex())

	allowance, err := ethutils.Erc20GetAllowance(client, token, common.HexToAddress(owner), Erc20Handler)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("allowance", allowance)

	nonce++
	err = ethutils.UpdateNonce(client)
	if err != nil {
		t.Fatal(err)
	}

	depositAmount := big.NewInt(0).Mul(big.NewInt(5), big.NewInt(100000000000000000))
	/// Stafi recepient
	//recepient := hexutil.MustDecode("0x3673009bdb664a3f3b6d9f69c9dd37fc0473551a249aa48542408b016ec62b2e")

	/// Eth recepient
	recepient := common.HexToAddress(owner).Bytes()
	err = ethutils.CreateErc20Deposit(client, BridgeContract, destId, recepient, depositAmount, rfisResourceId)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBscFisDepositEvent(t *testing.T) {
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kpI, err := keystore.KeypairFromAddress(owner, keystore.EthChain, keystorePath, false)
	if err != nil {
		t.Fatal(err)
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	client, err := ethutils.NewClient(BscTestEndpoint, kp)
	if err != nil {
		t.Fatal(err)
	}

	/// 1 for stafi, 2 for eth, 3 for bsc
	//destId := msg.ChainId(1)
	destId := msg.ChainId(1)

	nonce, err := ethutils.GetDepositNonce(client, bscBridgeContract, destId)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("nonce", nonce)

	token, err := ethutils.Erc20GetResourceId(client, bscErc20Handler, nativeResouceId)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("token", token.Hex())

	allowance, err := ethutils.Erc20GetAllowance(client, token, common.HexToAddress(owner), bscErc20Handler)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("allowance", allowance)

	nonce++
	err = ethutils.UpdateNonce(client)
	if err != nil {
		t.Fatal(err)
	}

	depositAmount := big.NewInt(0).Mul(big.NewInt(15), big.NewInt(100000000000000000))
	/// Stafi recepient
	recepient := hexutil.MustDecode("0x3673009bdb664a3f3b6d9f69c9dd37fc0473551a249aa48542408b016ec62b2e")

	/// Eth recepient
	//recepient := common.HexToAddress(owner).Bytes()
	err = ethutils.CreateErc20Deposit(client, bscBridgeContract, destId, recepient, depositAmount, nativeResouceId)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBscRFisDepositEvent(t *testing.T) {
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kpI, err := keystore.KeypairFromAddress(owner, keystore.EthChain, keystorePath, false)
	if err != nil {
		t.Fatal(err)
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	client, err := ethutils.NewClient(BscTestEndpoint, kp)
	if err != nil {
		t.Fatal(err)
	}

	/// 1 for stafi, 2 for eth, 3 for bsc
	//destId := msg.ChainId(1)
	destId := msg.ChainId(2)

	nonce, err := ethutils.GetDepositNonce(client, bscBridgeContract, destId)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("nonce", nonce)

	token, err := ethutils.Erc20GetResourceId(client, bscErc20Handler, rfisResourceId)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("token", token.Hex())

	allowance, err := ethutils.Erc20GetAllowance(client, token, common.HexToAddress(owner), bscErc20Handler)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("allowance", allowance)

	nonce++
	err = ethutils.UpdateNonce(client)
	if err != nil {
		t.Fatal(err)
	}

	depositAmount := big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1000000000000000000))
	/// Stafi recepient
	//recepient := hexutil.MustDecode("0x3673009bdb664a3f3b6d9f69c9dd37fc0473551a249aa48542408b016ec62b2e")

	/// Eth recepient
	recepient := common.HexToAddress(owner).Bytes()
	err = ethutils.CreateErc20Deposit(client, bscBridgeContract, destId, recepient, depositAmount, rfisResourceId)
	if err != nil {
		t.Fatal(err)
	}
}

func TestResourceId(t *testing.T) {
	t.Log(rfisResourceId.Hex())
}

func TestGetBalance(t *testing.T) {
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kpI, err := keystore.KeypairFromAddress(owner, keystore.EthChain, keystorePath, false)
	if err != nil {
		t.Fatal(err)
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	client, err := ethutils.NewClient(BscTestEndpoint, kp)
	if err != nil {
		t.Fatal(err)
	}

	bal, err := ethutils.Erc20GetBalance(client, bep20RFisContract, bscErc20Handler)
	if err != nil {
		panic(err)
	}

	t.Log("bal", bal)
}

func TestAddRelayer(t *testing.T) {
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kpI, err := keystore.KeypairFromAddress(owner, keystore.EthChain, keystorePath, false)
	if err != nil {
		t.Fatal(err)
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	client, err := ethutils.NewClient(BscTestEndpoint, kp)
	if err != nil {
		t.Fatal(err)
	}

	relayer := common.HexToAddress("0xaD0bf51f7fc89e262edBbdF53C260088B024D857")
	err = ethutils.AddRelayer(client, bscBridgeContract, relayer)
	if err != nil {
		t.Fatal(err)
	}
}

func TestResourceAndNonce(t *testing.T) {
	password := "123456"
	os.Setenv(keystore.EnvPassword, password)

	kpI, err := keystore.KeypairFromAddress(owner, keystore.EthChain, keystorePath, false)
	if err != nil {
		t.Fatal(err)
	}
	kp, _ := kpI.(*secp256k1.Keypair)

	client, err := ethutils.NewClient(BscTestEndpoint, kp)
	if err != nil {
		t.Fatal(err)
	}

	bridge, err := Bridge.NewBridge(bscBridgeContract, client.Client)
	if err != nil {
		t.Fatal(err)
	}

	handler, err := bridge.ResourceIDToHandlerAddress(nil, nativeResouceId)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(handler)

	nonce, err := bridge.OldDepositNonce(nil, 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(nonce)
}
