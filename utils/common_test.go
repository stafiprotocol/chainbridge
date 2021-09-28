package utils_test

import (
	"testing"

	"github.com/stafiprotocol/chainbridge/utils"
)

func TestCompareVersion(t *testing.T) {
	testCase := map[string]string{
		"1.7.0":  "1.7.0",
		"1.7.1":  "1.7.0",
		"1.6.09": "1.7.2",
		"1.16.9": "1.7.0",
		"1.06.9": "1.7.0",
	}

	for k, v := range testCase {
		result := utils.VersionCompare(k, v)
		t.Log(k, v, result)

	}
}
