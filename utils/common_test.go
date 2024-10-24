package utils_test

import (
	"strings"
	"testing"

	neutronClient "github.com/stafihub/neutron-relay-sdk/client"
	"github.com/stafihub/neutron-relay-sdk/common/log"
	"github.com/stafiprotocol/chainbridge/utils"
)

func TestCompareVersion(t *testing.T) {
	testCase := map[string]string{
		"1.7.0":  "1.7.0",
		"1.7.1":  "1.7.0",
		"1.6.09": "1.7.2",
		"1.16.9": "1.7.0",
		"1.06.9": "1.7.0",
		"1.9.9":  "1.7.0",
		"1.8.14": "1.7.0",
	}

	for k, v := range testCase {
		result := utils.VersionCompare(k, v)
		t.Log(k, v, result)

	}
}

func TestProposal(t *testing.T) {
	c, err := neutronClient.NewClient(nil, "", "", "neutron", []string{"https://rpc-falcron.pion-1.ntrn.tech:443"}, log.NewLog("client"))
	if err != nil {
		t.Log(err)
	}
	res, err := utils.QueryProposal(c, "neutron1eykzzmu8s2a2hy0a66un90p0u62frtz85qmq5jr5hx3z5xewtp0q9vlfkf", utils.QueryProposalParams{Amount: "0"})
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			t.Log(err)
		} else {
			t.Log("not found")
		}
	} else {
		t.Log(res)
	}
}
