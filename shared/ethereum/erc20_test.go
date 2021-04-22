package utils

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/math"
	"math/big"
	"testing"
)

func TestConstructErc20DepositData(t *testing.T) {
	depositAmount := big.NewInt(0).Mul(big.NewInt(15), big.NewInt(100))
	fmt.Println(depositAmount.BitLen())
	w := depositAmount.Bits()[0]
	fmt.Println(w)

	fmt.Println(byte(w))
	w >>= 8
	fmt.Println(w)
	fmt.Println(byte(w))
	fmt.Println(math.PaddedBigBytes(depositAmount, 32))
}
