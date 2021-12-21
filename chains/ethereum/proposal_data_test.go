package ethereum

import (
	"encoding/hex"
	// "github.com/stretchr/testify/assert"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func TestConstructErc20ProposalData(t *testing.T) {
	a := big.NewInt(1000000).Bytes()
	b := hexutil.MustDecode("0x2e0ddc00b2fddeb5da10b71f74611036ed2156862c0a4d4ac9a0ab08bb4c9a45")
	c := ConstructErc20ProposalData(a, b)
	t.Log(hex.EncodeToString(c))
	// expected := "00000000000000000000000000000000000000000000000000000000000f4240000000000000000000000000000000000000000000000000000000000000003c307861626364656661626364656661626162616162616161616161626162616261626162616261626162616261626162616261626162616261626162"
	// assert.Equal(t, expected, hex.EncodeToString(c))
}
