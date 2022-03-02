package stafihub

import (
	"context"

	"github.com/stafihub/rtoken-relay-core/common/core"
	stafiHubXBridgeTypes "github.com/stafihub/stafihub/x/bridge/types"
	stafiHubXLedgerTypes "github.com/stafihub/stafihub/x/ledger/types"
)

func (c *Client) QuerySnapshot(shotId string) (*stafiHubXLedgerTypes.QueryGetSnapshotResponse, error) {
	done := core.UseSdkConfigContext(AccountPrefix)
	defer done()

	queryClient := stafiHubXLedgerTypes.NewQueryClient(c.Ctx())
	params := &stafiHubXLedgerTypes.QueryGetSnapshotRequest{
		ShotId: shotId,
	}

	cc, err := Retry(func() (interface{}, error) {
		return queryClient.GetSnapshot(context.Background(), params)
	})
	if err != nil {
		return nil, err
	}
	return cc.(*stafiHubXLedgerTypes.QueryGetSnapshotResponse), nil
}

func (c *Client) QueryPoolUnbond(denom, pool string, era uint32) (*stafiHubXLedgerTypes.QueryGetPoolUnbondResponse, error) {
	done := core.UseSdkConfigContext(AccountPrefix)
	defer done()

	queryClient := stafiHubXLedgerTypes.NewQueryClient(c.Ctx())
	params := &stafiHubXLedgerTypes.QueryGetPoolUnbondRequest{
		Denom: denom,
		Pool:  pool,
		Era:   era,
	}

	cc, err := Retry(func() (interface{}, error) {
		return queryClient.GetPoolUnbond(context.Background(), params)
	})
	if err != nil {
		return nil, err
	}
	return cc.(*stafiHubXLedgerTypes.QueryGetPoolUnbondResponse), nil
}

func (c *Client) QueryPoolDetail(denom, pool string) (*stafiHubXLedgerTypes.QueryGetPoolDetailResponse, error) {
	done := core.UseSdkConfigContext(AccountPrefix)
	defer done()

	queryClient := stafiHubXLedgerTypes.NewQueryClient(c.Ctx())
	params := &stafiHubXLedgerTypes.QueryGetPoolDetailRequest{
		Denom: denom,
		Pool:  pool,
	}

	cc, err := Retry(func() (interface{}, error) {
		return queryClient.GetPoolDetail(context.Background(), params)
	})
	if err != nil {
		return nil, err
	}
	return cc.(*stafiHubXLedgerTypes.QueryGetPoolDetailResponse), nil
}

func (c *Client) QuerySignature(denom, pool string, era uint32, txType stafiHubXLedgerTypes.OriginalTxType, proposalId string) (*stafiHubXLedgerTypes.QueryGetSignatureResponse, error) {
	done := core.UseSdkConfigContext(AccountPrefix)
	defer done()

	queryClient := stafiHubXLedgerTypes.NewQueryClient(c.Ctx())
	params := &stafiHubXLedgerTypes.QueryGetSignatureRequest{
		Denom:  denom,
		Era:    era,
		Pool:   pool,
		TxType: txType,
		PropId: proposalId,
	}

	cc, err := Retry(func() (interface{}, error) {
		return queryClient.GetSignature(context.Background(), params)
	})
	if err != nil {
		return nil, err
	}
	return cc.(*stafiHubXLedgerTypes.QueryGetSignatureResponse), nil
}

func (c *Client) QueryPools(denom string) (*stafiHubXLedgerTypes.QueryBondedPoolsByDenomResponse, error) {
	done := core.UseSdkConfigContext(AccountPrefix)
	defer done()

	queryClient := stafiHubXLedgerTypes.NewQueryClient(c.Ctx())
	params := &stafiHubXLedgerTypes.QueryBondedPoolsByDenomRequest{
		Denom: denom,
	}

	cc, err := Retry(func() (interface{}, error) {
		return queryClient.BondedPoolsByDenom(context.Background(), params)
	})
	if err != nil {
		return nil, err
	}
	return cc.(*stafiHubXLedgerTypes.QueryBondedPoolsByDenomResponse), nil
}

func (c *Client) QueryChainEra(denom string) (*stafiHubXLedgerTypes.QueryGetChainEraResponse, error) {
	done := core.UseSdkConfigContext(AccountPrefix)
	defer done()

	queryClient := stafiHubXLedgerTypes.NewQueryClient(c.Ctx())
	params := &stafiHubXLedgerTypes.QueryGetChainEraRequest{
		Denom: denom,
	}

	cc, err := Retry(func() (interface{}, error) {
		return queryClient.GetChainEra(context.Background(), params)
	})
	if err != nil {
		return nil, err
	}
	return cc.(*stafiHubXLedgerTypes.QueryGetChainEraResponse), nil
}

func (c *Client) QueryEraSnapShotList(denom string, era uint32) (*stafiHubXLedgerTypes.QueryGetEraSnapshotResponse, error) {
	done := core.UseSdkConfigContext(AccountPrefix)
	defer done()

	queryClient := stafiHubXLedgerTypes.NewQueryClient(c.Ctx())
	params := &stafiHubXLedgerTypes.QueryGetEraSnapshotRequest{
		Denom: denom,
		Era:   era,
	}

	cc, err := Retry(func() (interface{}, error) {
		return queryClient.GetEraSnapshot(context.Background(), params)
	})
	if err != nil {
		return nil, err
	}
	return cc.(*stafiHubXLedgerTypes.QueryGetEraSnapshotResponse), nil
}

func (c *Client) QueryEraContinuable(denom string, era uint32) (bool, error) {
	done := core.UseSdkConfigContext(AccountPrefix)
	defer done()

	queryClient := stafiHubXLedgerTypes.NewQueryClient(c.Ctx())
	params := &stafiHubXLedgerTypes.QueryGetEraSnapshotRequest{
		Denom: denom,
		Era:   era,
	}

	cc, err := Retry(func() (interface{}, error) {
		return queryClient.GetEraSnapshot(context.Background(), params)
	})
	if err != nil {
		return false, err
	}
	res := cc.(*stafiHubXLedgerTypes.QueryGetEraSnapshotResponse)
	if len(res.ShotIds) > 0 {
		return false, nil
	}
	return true, nil
}

func (c *Client) QueryRParams(denom string) (*stafiHubXLedgerTypes.QueryGetRParamsResponse, error) {
	done := core.UseSdkConfigContext(AccountPrefix)
	defer done()
	queryClient := stafiHubXLedgerTypes.NewQueryClient(c.Ctx())
	params := &stafiHubXLedgerTypes.QueryGetRParamsRequest{
		Denom: denom,
	}
	cc, err := Retry(func() (interface{}, error) {
		return queryClient.GetRParams(context.Background(), params)
	})
	if err != nil {
		return nil, err
	}
	return cc.(*stafiHubXLedgerTypes.QueryGetRParamsResponse), nil
}

func (c *Client) QueryBondRecord(denom, txHash string) (*stafiHubXLedgerTypes.QueryGetBondRecordResponse, error) {
	done := core.UseSdkConfigContext(AccountPrefix)
	defer done()

	queryClient := stafiHubXLedgerTypes.NewQueryClient(c.Ctx())
	params := &stafiHubXLedgerTypes.QueryGetBondRecordRequest{
		Denom:  denom,
		Txhash: txHash,
	}
	cc, err := Retry(func() (interface{}, error) {
		return queryClient.GetBondRecord(context.Background(), params)
	})
	if err != nil {
		return nil, err
	}
	return cc.(*stafiHubXLedgerTypes.QueryGetBondRecordResponse), nil
}

func (c *Client) QueryBridgeProposalDetail(chainId uint32, depositNonce uint64, resourceId, amount, receiver string) (*stafiHubXBridgeTypes.QueryProposalDetailResponse, error) {
	done := core.UseSdkConfigContext(AccountPrefix)
	defer done()

	queryClient := stafiHubXBridgeTypes.NewQueryClient(c.Ctx())
	params := &stafiHubXBridgeTypes.QueryProposalDetailRequest{
		ChainId:      chainId,
		DepositNonce: depositNonce,
		ResourceId:   resourceId,
		Amount:       amount,
		Receiver:     receiver,
	}
	cc, err := Retry(func() (interface{}, error) {
		return queryClient.ProposalDetail(context.Background(), params)
	})
	if err != nil {
		return nil, err
	}
	return cc.(*stafiHubXBridgeTypes.QueryProposalDetailResponse), nil
}
