package ethereum

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestIsAddressValid(t *testing.T) {
	addr := "0xBd39f5936969828eD9315220659cD11129071814"
	assert.Equal(t, true, IsAddressValid(addr))
	assert.Equal(t, true, IsAddressValid(strings.ToLower(addr)))

	a := "0xabcd"
	assert.Equal(t, false, IsAddressValid(a))
	a1 := "0xBd39f5936969g28eD9315220659cD11129071814"
	assert.Equal(t, false, IsAddressValid(a1))
	a2 := "Bd39f5936969828eD9315220659cD11129071814"
	assert.Equal(t, false, IsAddressValid(a2))
}
