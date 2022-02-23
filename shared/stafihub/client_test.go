package stafihub_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	// "strings"
	"sync"
	"testing"
	"time"

	"github.com/JFJun/go-substrate-crypto/ss58"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	hubClient "github.com/stafihub/stafi-hub-relay-sdk/client"
	"github.com/stretchr/testify/assert"
)

var client *hubClient.Client

//eda331e37bf66b2393c4c271e384dfaa2bfcdd35
var addrMultiSig1, _ = types.AccAddressFromBech32("cosmos1em384d8ek3y8nlugapz7p5k5skg58j66je3las")
var addrMultiSig2, _ = types.AccAddressFromBech32("cosmos1wmk9ys49zxgmx77pqs7cjnpamnnxuspqu2r87y")

var addrValidatorTestnet, _ = types.ValAddressFromBech32("cosmosvaloper17tpddyr578avyn95xngkjl8nl2l2tf6auh8kpc")
var addrValidatorTestnetStation, _ = types.ValAddressFromBech32("cosmosvaloper1x5wgh6vwye60wv3dtshs9dmqggwfx2ldk5cvqu")
var addrValidatorTestnetAteam, _ = types.ValAddressFromBech32("cosmosvaloper105gvcjgs6s4j5ws9srckx0drt4x8cwgywplh7p")

var adrValidatorTestnetTecos, _ = types.ValAddressFromBech32("cosmosvaloper1p7e37nztj62mmra8xhgqde7sql3llhhu6hvcx8")
var adrValidatorEverStake, _ = types.ValAddressFromBech32("cosmosvaloper1tflk30mq5vgqjdly92kkhhq3raev2hnz6eete3")
var adrValidatorForbole, _ = types.ValAddressFromBech32("cosmosvaloper1w96rrh9sx0h7n7qak00l90un0kx5wala2prmxt")

var addrReleaseAddress, _ = types.AccAddressFromBech32("cosmos1wmk9ys49zxgmx77pqs7cjnpamnnxuspqu2r87y")
var addrReleaseValWetez, _ = types.ValAddressFromBech32("cosmosvaloper1s05va5d09xlq3et8mapsesqh6r5lqy7mkhwshm")

func TestGetAddrHex(t *testing.T) {
	t.Log("cosmosvaloper17tpddyr578avyn95xngkjl8nl2l2tf6auh8kpc", hexutil.Encode(addrValidatorTestnet.Bytes()))
	t.Log("cosmosvaloper1x5wgh6vwye60wv3dtshs9dmqggwfx2ldk5cvqu", hexutil.Encode(addrValidatorTestnetStation.Bytes()))
	t.Log("cosmosvaloper105gvcjgs6s4j5ws9srckx0drt4x8cwgywplh7p", hexutil.Encode(addrValidatorTestnetAteam.Bytes()))

	t.Log("cosmosvaloper1p7e37nztj62mmra8xhgqde7sql3llhhu6hvcx8", hexutil.Encode(adrValidatorTestnetTecos.Bytes()))
	t.Log("cosmosvaloper1tflk30mq5vgqjdly92kkhhq3raev2hnz6eete3", hexutil.Encode(adrValidatorEverStake.Bytes()))
	t.Log("cosmosvaloper1w96rrh9sx0h7n7qak00l90un0kx5wala2prmxt", hexutil.Encode(adrValidatorForbole.Bytes()))
}

func initClient() {
	// key, err := keyring.New(types.KeyringServiceName(), keyring.BackendFile, "/Users/tpkeeper/.gaia", strings.NewReader("tpkeeper\n"))
	// if err != nil {
	// 	panic(err)
	// }

	var err error
	// client, err = rpc.NewClient(key, "stargate-final", "key0", "0.04umuon", "umuon", "https://testcosmosrpc.wetez.io:443")
	client, err = hubClient.NewClient(nil, "my-account", "", "0.04stake", "stake", "http://127.0.0.1:26657")
	// client, err = hubClient.NewClient(nil, "my-account", "", "0.04stake", "stake", "https://testcosmosrpc.wetez.io:443")
	// client, _ = rpc.NewClient(key, "cosmoshub-4", "self", "0.00001uatom", "uatom", "https://cosmos-rpc1.stafi.io:443")
	// client, err = hubClient.NewClient(nil, "cosmoshub-4", "", "0.00001uatom", "uatom", "https://cosmos-rpc1.stafi.io:443")
	if err != nil {
		panic(err)
	}
}

func TestClient_SendTo(t *testing.T) {
	initClient()
	err := client.SingleTransferTo(addrMultiSig1, types.NewCoins(types.NewInt64Coin(client.GetDenom(), 50000)))
	assert.NoError(t, err)
}

func TestClient_QueryTxByHash(t *testing.T) {
	initClient()

	for {

		res, err := client.QueryTxByHash("cb21ff0aae81b98f39e1488bb26c098597fab336db45521be7d8aad057a973e3")
		assert.NoError(t, err)
		t.Log(res.Height)

		curBlock, err := client.GetCurrentBlockHeight()
		assert.NoError(t, err)
		t.Log(curBlock)
		time.Sleep(1 * time.Second)
		t.Log("\n")
	}
}

func TestGetTxs(t *testing.T) {
	initClient()
	txs, err := client.GetBlockTxs(5879)
	if err != nil {
		t.Fatal(err)
	}
	for _, tx := range txs {
		// t.Log(fmt.Sprintf("%+v", tx.Logs))
		for _, log := range tx.Logs {
			for _, event := range log.Events {
				t.Log(fmt.Sprintf("%+v", event))
			}
		}

	}
}

func TestGetPubKey(t *testing.T) {
	initClient()
	test, _ := types.AccAddressFromBech32("cosmos1u22lut8qgqg8znxam72pwgqp8c09rnvme00kea")
	account, _ := client.QueryAccount(test)
	t.Log(hex.EncodeToString(account.GetPubKey().Bytes()))

}

func TestClient_Sign(t *testing.T) {
	initClient()
	bts, err := hex.DecodeString("0E4F8F8FF7A3B67121711DA17FBE5AE8CB25DB272DDBF7DC0E02122947266604")
	assert.NoError(t, err)
	sigs, pubkey, err := client.Sign("recipient", bts)
	assert.NoError(t, err)
	t.Log(hex.EncodeToString(sigs))
	//4c6902bda88424923c62f95b3e3ead40769edab4ec794108d1c18994fac90d490087815823bd1a8af3d6a0271538cef4622b4b500a6253d2bd4c80d38e95aa6d
	t.Log(hex.EncodeToString(pubkey.Bytes()))
	//02e7710b4f7147c10ad90da06b69d2d6b8ff46786ef55a3f1e889c33de2bf0b416
}

func TestAddress(t *testing.T) {
	addrKey1, _ := types.AccAddressFromBech32("cosmos1a8mg9rj4nklhmwkf5vva8dvtgx4ucd9yjasret")
	addrKey2, _ := types.AccAddressFromBech32("cosmos1ztquzhpkve7szl99jkugq4l8jtpnhln76aetam")
	addrKey3, _ := types.AccAddressFromBech32("cosmos12zz2hm02sxe9f4pwt7y5q9wjhcu98vnuwmjz4x")
	addrKey4, _ := types.AccAddressFromBech32("cosmos12yprrdprzat35zhqxe2fcnn3u26gwlt6xcq0pj")
	addrKey5, _ := types.AccAddressFromBech32("cosmos1em384d8ek3y8nlugapz7p5k5skg58j66je3las")
	t.Log(hex.EncodeToString(addrKey1.Bytes()))
	t.Log(hex.EncodeToString(addrKey2.Bytes()))
	t.Log(hex.EncodeToString(addrKey3.Bytes()))
	t.Log(hex.EncodeToString(addrKey4.Bytes()))
	t.Log(hex.EncodeToString(addrKey5.Bytes()))
	//client_test.go:347: e9f6828e559dbf7dbac9a319d3b58b41abcc34a4
	//client_test.go:348: 12c1c15c36667d017ca595b88057e792c33bfe7e
	//client_test.go:349: 5084abedea81b254d42e5f894015d2be3853b27c
}

func TestClient_QueryDelegations(t *testing.T) {
	initClient()
	addr, err := types.AccAddressFromBech32("cosmos12yprrdprzat35zhqxe2fcnn3u26gwlt6xcq0pj")
	assert.NoError(t, err)
	res, err := client.QueryDelegations(addr, 2458080)
	assert.NoError(t, err)
	t.Log(res.String())
	for i, d := range res.GetDelegationResponses() {
		t.Log(i, d.Balance.Amount.IsZero())
	}
}

func TestClient_QueryBalance(t *testing.T) {
	initClient()
	res, err := client.QueryBalance(addrMultiSig1, "umuon", 440000)
	assert.NoError(t, err)
	t.Log(res.Balance.Amount)
}

func TestClient_QueryDelegationTotalRewards(t *testing.T) {
	initClient()
	addr, err := types.AccAddressFromBech32("cosmos12yprrdprzat35zhqxe2fcnn3u26gwlt6xcq0pj")
	assert.NoError(t, err)
	t.Log(client.GetDenom())
	res, err := client.QueryDelegationTotalRewards(addr, 2458080)
	assert.NoError(t, err)
	for i, _ := range res.Rewards {
		t.Log(i, res.Rewards[i].Reward.AmountOf(client.GetDenom()))
		t.Log(i, res.Rewards[i].Reward.AmountOf(client.GetDenom()).TruncateInt())

	}
	t.Log("total ", res.GetTotal().AmountOf(client.GetDenom()).TruncateInt())
}

func TestClient_GetSequence(t *testing.T) {
	initClient()
	seq, err := client.GetSequence(0, addrReleaseAddress)
	assert.NoError(t, err)
	t.Log(seq)
	// txRes,err:=client.QueryTxByHash("FBD05BD4B9DB0386B16E679184EAC88D444B38DA992F8AFC35B5A580B3FC6AA4")
	// assert.NoError(t,err)
	// t.Log(txRes.String())

	res, err := client.QueryUnbondingDelegation(addrReleaseAddress, addrReleaseValWetez, 0)
	if err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			t.Log(err)
		}

		t.Fatal(err)
	}
	t.Log(res.String())
	t.Log(len(res.GetUnbond().Entries))
}
func TestMemo(t *testing.T) {
	initClient()
	res, err := client.QueryTxByHash("c7e3f7baf5a5f1d8cbc112080f32070dddd7cca5fe4272e06f8d42c17b25193f")
	assert.NoError(t, err)
	tx, err := client.GetTxConfig().TxDecoder()(res.Tx.GetValue())
	//tx, err := client.GetTxConfig().TxJSONDecoder()(res.Tx.Value)
	assert.NoError(t, err)
	memoTx, ok := tx.(types.TxWithMemo)
	assert.Equal(t, true, ok)
	t.Log(memoTx.GetMemo())
	hb, _ := hexutil.Decode("0xbebd0355ae360c8e6a7ed940a819838c66ca7b8f581f9c0e81dbb5faff346a30")
	//t.Log(string(hb))
	bonderAddr, _ := ss58.Encode(hb, ss58.StafiPrefix)
	t.Log(bonderAddr)
}

func TestMultiThread(t *testing.T) {
	initClient()
	wg := sync.WaitGroup{}
	wg.Add(50)

	for i := 0; i < 50; i++ {
		go func(i int) {
			t.Log(i)
			time.Sleep(5 * time.Second)
			height, err := client.GetAccount()
			if err != nil {
				t.Log("fail", i, err)
			} else {
				t.Log("success", i, height.GetSequence())
			}
			time.Sleep(15 * time.Second)
			height, err = client.GetAccount()
			if err != nil {
				t.Log("fail", i, err)
			} else {
				t.Log("success", i, height.GetSequence())
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func TestSort(t *testing.T) {
	a := []string{"cosmos1kuyde8vpt8c0ty4pxqgxw3makse7md80umthvg"}
	t.Log(a)
	sort.SliceStable(a, func(i, j int) bool {
		return bytes.Compare([]byte(a[i]), []byte(a[j])) < 0
	})
	t.Log(a)
	// rawTx := "7b22626f6479223a7b226d65737361676573223a5b7b224074797065223a222f636f736d6f732e62616e6b2e763162657461312e4d73674d756c746953656e64222c22696e70757473223a5b7b2261646472657373223a22636f736d6f7331776d6b39797334397a78676d78373770717337636a6e70616d6e6e7875737071753272383779222c22636f696e73223a5b7b2264656e6f6d223a227561746f6d222c22616d6f756e74223a22313134373730227d5d7d5d2c226f757470757473223a5b7b2261646472657373223a22636f736d6f733135366b6b326b71747777776670733836673534377377646c7263326377367163746d36633877222c22636f696e73223a5b7b2264656e6f6d223a227561746f6d222c22616d6f756e74223a2231393936227d5d7d2c7b2261646472657373223a22636f736d6f73316b7579646538767074386330747934707871677877336d616b7365376d643830756d74687667222c22636f696e73223a5b7b2264656e6f6d223a227561746f6d222c22616d6f756e74223a223939383030227d5d7d2c7b2261646472657373223a22636f736d6f73316a6b6b68666c753871656471743463796173643674673730676a7778346a6b6872736536727a222c22636f696e73223a5b7b2264656e6f6d223a227561746f6d222c22616d6f756e74223a223132393734227d5d7d5d7d5d2c226d656d6f223a22222c2274696d656f75745f686569676874223a2230222c22657874656e73696f6e5f6f7074696f6e73223a5b5d2c226e6f6e5f637269746963616c5f657874656e73696f6e5f6f7074696f6e73223a5b5d7d2c22617574685f696e666f223a7b227369676e65725f696e666f73223a5b5d2c22666565223a7b22616d6f756e74223a5b7b2264656e6f6d223a227561746f6d222c22616d6f756e74223a2237353030227d5d2c226761735f6c696d6974223a2231353030303030222c227061796572223a22222c226772616e746572223a22227d7d2c227369676e617475726573223a5b5d7d"
	rawTx := "7b22626f6479223a7b226d65737361676573223a5b7b224074797065223a222f636f736d6f732e62616e6b2e763162657461312e4d73674d756c746953656e64222c22696e70757473223a5b7b2261646472657373223a22636f736d6f7331776d6b39797334397a78676d78373770717337636a6e70616d6e6e7875737071753272383779222c22636f696e73223a5b7b2264656e6f6d223a227561746f6d222c22616d6f756e74223a2231303539363338303631227d5d7d5d2c226f757470757473223a5b7b2261646472657373223a22636f736d6f733135366b6b326b71747777776670733836673534377377646c7263326377367163746d36633877222c22636f696e73223a5b7b2264656e6f6d223a227561746f6d222c22616d6f756e74223a2231303539363338303631227d5d7d5d7d5d2c226d656d6f223a22222c2274696d656f75745f686569676874223a2230222c22657874656e73696f6e5f6f7074696f6e73223a5b5d2c226e6f6e5f637269746963616c5f657874656e73696f6e5f6f7074696f6e73223a5b5d7d2c22617574685f696e666f223a7b227369676e65725f696e666f73223a5b5d2c22666565223a7b22616d6f756e74223a5b7b2264656e6f6d223a227561746f6d222c22616d6f756e74223a2237353030227d5d2c226761735f6c696d6974223a2231353030303030222c227061796572223a22222c226772616e746572223a22227d7d2c227369676e617475726573223a5b5d7d"
	// rawTx:="7b22626f6479223a7b226d65737361676573223a5b7b224074797065223a222f636f736d6f732e62616e6b2e763162657461312e4d73674d756c746953656e64222c22696e70757473223a5b7b2261646472657373223a22636f736d6f7331776d6b39797334397a78676d78373770717337636a6e70616d6e6e7875737071753272383779222c22636f696e73223a5b7b2264656e6f6d223a227561746f6d222c22616d6f756e74223a22313134373730227d5d7d5d2c226f757470757473223a5b7b2261646472657373223a22636f736d6f73316a6b6b68666c753871656471743463796173643674673730676a7778346a6b6872736536727a222c22636f696e73223a5b7b2264656e6f6d223a227561746f6d222c22616d6f756e74223a223132393734227d5d7d2c7b2261646472657373223a22636f736d6f733135366b6b326b71747777776670733836673534377377646c7263326377367163746d36633877222c22636f696e73223a5b7b2264656e6f6d223a227561746f6d222c22616d6f756e74223a2231393936227d5d7d2c7b2261646472657373223a22636f736d6f73316b7579646538767074386330747934707871677877336d616b7365376d643830756d74687667222c22636f696e73223a5b7b2264656e6f6d223a227561746f6d222c22616d6f756e74223a223939383030227d5d7d5d7d5d2c226d656d6f223a22222c2274696d656f75745f686569676874223a2230222c22657874656e73696f6e5f6f7074696f6e73223a5b5d2c226e6f6e5f637269746963616c5f657874656e73696f6e5f6f7074696f6e73223a5b5d7d2c22617574685f696e666f223a7b227369676e65725f696e666f73223a5b5d2c22666565223a7b22616d6f756e74223a5b7b2264656e6f6d223a227561746f6d222c22616d6f756e74223a2237353030227d5d2c226761735f6c696d6974223a2231353030303030222c227061796572223a22222c226772616e746572223a22227d7d2c227369676e617475726573223a5b5d7d"
	txBts, err := hex.DecodeString(rawTx)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(txBts))
}
