package utils

import (
	"encoding/json"
	"fmt"
	"github.com/stafihub/neutron-relay-sdk/client"
	"strconv"
	"strings"
)

func VersionCompare(version1, version2 string) (ret int) {
	defer func() {
		switch {
		case ret < 0:
			ret = -1
		case ret > 0:
			ret = 1
		}
	}()

	version1Slice := strings.Split(version1, ".")
	version2Slice := strings.Split(version2, ".")
	if len(version1Slice) != len(version2Slice) || len(version1Slice) != 3 {
		panic(fmt.Sprintf("version format err: 1 %s 2 %s", version1, version2))
	}

	for i := range version1Slice {
		v1, err := strconv.Atoi(version1Slice[i])
		if err != nil {
			panic(err)
		}
		v2, err := strconv.Atoi(version2Slice[i])
		if err != nil {
			panic(err)
		}
		if v1 == v2 {
			continue
		} else {
			return v1 - v2
		}

	}
	return 0

}

type Params struct {
	ChainId      uint64 `json:"chain_id"`
	DepositNonce uint64 `json:"deposit_nonce"`
	Recipient    string `json:"recipient"`
	Amount       string `json:"amount"`
}

type QueryProposalReq struct {
	Proposal Params `json:"proposal"`
}

type QueryProposalRes struct {
	ChainId      uint64   `json:"chain_id"`
	DepositNonce uint64   `json:"deposit_nonce"`
	Recipient    string   `json:"recipient"`
	Amount       string   `json:"amount"`
	Executed     bool     `json:"executed"`
	Voters       []string `json:"voters"`
}

func getQueryProposalReq(params Params) []byte {
	poolReq := QueryProposalReq{
		Proposal: params,
	}
	marshal, _ := json.Marshal(poolReq)
	return marshal
}
func QueryProposal(neutronClient *client.Client, contract string, params Params) (*QueryProposalRes, error) {
	poolInfoRes, err := neutronClient.QuerySmartContractState(contract, getQueryProposalReq(params))
	if err != nil {
		return nil, err
	}
	var res QueryProposalRes
	err = json.Unmarshal(poolInfoRes.Data.Bytes(), &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

type VoteProposalParams struct {
	ChainId      uint64 `json:"chain_id"`
	DepositNonce uint64 `json:"deposit_nonce"`
	Recipient    string `json:"recipient"`
	Amount       string `json:"amount"`
}

func VoteProposalMsg(params VoteProposalParams) []byte {
	msg := struct {
		VoteProposalParams `json:"vote_proposal"`
	}{
		VoteProposalParams: params,
	}
	marshal, _ := json.Marshal(msg)
	return marshal
}
