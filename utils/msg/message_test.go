package msg

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResourceIdFromSlice(t *testing.T) {
	src := "0x000000000000000000000000000000a9e0095b8965c01e6a09c97938f3860901"
	b, _ := hexutil.Decode(src)

	rId := ResourceIdFromSlice(b)
	assert.Equal(t, src, "0x"+rId.Hex())
}
