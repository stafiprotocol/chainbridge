package substrate

import (
	"github.com/itering/scale.go/source"
	"github.com/itering/scale.go/types"
	"io/ioutil"
	"testing"

	"github.com/stafiprotocol/chainbridge-utils/core"
)

func TestParseStartBlock(t *testing.T) {
	// Valid option included in config
	cfg := &core.ChainConfig{Opts: map[string]string{"startBlock": "1000"}}

	blk := parseStartBlock(cfg)

	if blk != 1000 {
		t.Fatalf("Got: %d Expected: %d", blk, 1000)
	}

	// Not included in config
	cfg = &core.ChainConfig{Opts: map[string]string{}}

	blk = parseStartBlock(cfg)

	if blk != 0 {
		t.Fatalf("Got: %d Expected: %d", blk, 0)
	}
}

func TestSome(t *testing.T) {
	types.RuntimeType{}.Reg()
	content, err := ioutil.ReadFile(DefaultTypeFilePath)
	if err != nil {
		panic(err)
	}
	types.RegCustomTypes(source.LoadTypeRegistry(content))
}
