package ethereum

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestConstructErc20ProposalData(t *testing.T) {
	a := big.NewInt(1000000).Bytes()
	b := []byte("0xabcdefabcdefababaabaaaaaababababababababababababababababab")
	c := ConstructErc20ProposalData(a, b)
	expected := "00000000000000000000000000000000000000000000000000000000000f4240000000000000000000000000000000000000000000000000000000000000003c307861626364656661626364656661626162616162616161616161626162616261626162616261626162616261626162616261626162616261626162"
	assert.Equal(t, expected, hex.EncodeToString(c))
}
