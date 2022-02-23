package stafihub

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"syscall"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	xAuthTx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	xBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	xDistriTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	xStakeTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stafihub/rtoken-relay-core/common/core"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

const retryLimit = 60
const waitTime = time.Millisecond * 500

//no 0x prefix
func (c *Client) QueryTxByHash(hashHexStr string) (*types.TxResponse, error) {
	done := core.UseSdkConfigContext(AccountPrefix)
	defer done()
	cc, err := retry(func() (interface{}, error) {
		return xAuthTx.QueryTx(c.clientCtx, hashHexStr)
	})
	if err != nil {
		return nil, err
	}
	return cc.(*types.TxResponse), nil
}

func (c *Client) QueryDelegation(delegatorAddr types.AccAddress, validatorAddr types.ValAddress, height int64) (*xStakeTypes.QueryDelegationResponse, error) {
	done := core.UseSdkConfigContext(AccountPrefix)
	defer done()
	client := c.clientCtx.WithHeight(height)
	queryClient := xStakeTypes.NewQueryClient(client)
	params := &xStakeTypes.QueryDelegationRequest{
		DelegatorAddr: delegatorAddr.String(),
		ValidatorAddr: validatorAddr.String(),
	}

	cc, err := retry(func() (interface{}, error) {
		return queryClient.Delegation(context.Background(), params)
	})
	if err != nil {
		return nil, err
	}
	return cc.(*xStakeTypes.QueryDelegationResponse), nil
}

func (c *Client) QueryUnbondingDelegation(delegatorAddr types.AccAddress, validatorAddr types.ValAddress, height int64) (*xStakeTypes.QueryUnbondingDelegationResponse, error) {
	done := core.UseSdkConfigContext(AccountPrefix)
	defer done()
	client := c.clientCtx.WithHeight(height)
	queryClient := xStakeTypes.NewQueryClient(client)
	params := &xStakeTypes.QueryUnbondingDelegationRequest{
		DelegatorAddr: delegatorAddr.String(),
		ValidatorAddr: validatorAddr.String(),
	}

	cc, err := retry(func() (interface{}, error) {
		return queryClient.UnbondingDelegation(context.Background(), params)
	})
	if err != nil {
		return nil, err
	}
	return cc.(*xStakeTypes.QueryUnbondingDelegationResponse), nil
}

func (c *Client) QueryDelegations(delegatorAddr types.AccAddress, height int64) (*xStakeTypes.QueryDelegatorDelegationsResponse, error) {
	done := core.UseSdkConfigContext(AccountPrefix)
	defer done()
	client := c.clientCtx.WithHeight(height)
	queryClient := xStakeTypes.NewQueryClient(client)
	params := &xStakeTypes.QueryDelegatorDelegationsRequest{
		DelegatorAddr: delegatorAddr.String(),
		Pagination:    &query.PageRequest{},
	}
	cc, err := retry(func() (interface{}, error) {
		return queryClient.DelegatorDelegations(context.Background(), params)
	})
	if err != nil {
		return nil, err
	}
	return cc.(*xStakeTypes.QueryDelegatorDelegationsResponse), nil
}

func (c *Client) QueryDelegationRewards(delegatorAddr types.AccAddress, validatorAddr types.ValAddress, height int64) (*xDistriTypes.QueryDelegationRewardsResponse, error) {
	done := core.UseSdkConfigContext(AccountPrefix)
	defer done()
	client := c.clientCtx.WithHeight(height)
	queryClient := xDistriTypes.NewQueryClient(client)
	cc, err := retry(func() (interface{}, error) {
		return queryClient.DelegationRewards(
			context.Background(),
			&xDistriTypes.QueryDelegationRewardsRequest{DelegatorAddress: delegatorAddr.String(), ValidatorAddress: validatorAddr.String()},
		)
	})
	if err != nil {
		return nil, err
	}
	return cc.(*xDistriTypes.QueryDelegationRewardsResponse), nil
}

func (c *Client) QueryDelegationTotalRewards(delegatorAddr types.AccAddress, height int64) (*xDistriTypes.QueryDelegationTotalRewardsResponse, error) {
	done := core.UseSdkConfigContext(AccountPrefix)
	defer done()
	client := c.clientCtx.WithHeight(height)
	queryClient := xDistriTypes.NewQueryClient(client)

	cc, err := retry(func() (interface{}, error) {
		return queryClient.DelegationTotalRewards(
			context.Background(),
			&xDistriTypes.QueryDelegationTotalRewardsRequest{DelegatorAddress: delegatorAddr.String()},
		)
	})
	if err != nil {
		return nil, err
	}
	return cc.(*xDistriTypes.QueryDelegationTotalRewardsResponse), nil
}

func (c *Client) QueryBlock(height int64) (*ctypes.ResultBlock, error) {
	done := core.UseSdkConfigContext(AccountPrefix)
	defer done()
	node, err := c.clientCtx.GetNode()
	if err != nil {
		return nil, err
	}

	cc, err := retry(func() (interface{}, error) {
		return node.Block(context.Background(), &height)
	})
	if err != nil {
		return nil, err
	}
	return cc.(*ctypes.ResultBlock), nil
}

func (c *Client) QueryAccount(addr types.AccAddress) (client.Account, error) {
	done := core.UseSdkConfigContext(AccountPrefix)
	defer done()
	return c.getAccount(0, addr)
}

func (c *Client) GetSequence(height int64, addr types.AccAddress) (uint64, error) {
	done := core.UseSdkConfigContext(AccountPrefix)
	defer done()
	account, err := c.getAccount(height, addr)
	if err != nil {
		return 0, err
	}
	return account.GetSequence(), nil
}

func (c *Client) QueryBalance(addr types.AccAddress, denom string, height int64) (*xBankTypes.QueryBalanceResponse, error) {
	done := core.UseSdkConfigContext(AccountPrefix)
	defer done()
	client := c.clientCtx.WithHeight(height)
	queryClient := xBankTypes.NewQueryClient(client)
	params := xBankTypes.NewQueryBalanceRequest(addr, denom)

	cc, err := retry(func() (interface{}, error) {
		return queryClient.Balance(context.Background(), params)
	})
	if err != nil {
		return nil, err
	}
	return cc.(*xBankTypes.QueryBalanceResponse), nil
}

func (c *Client) GetCurrentBlockHeight() (int64, error) {
	done := core.UseSdkConfigContext(AccountPrefix)
	defer done()
	status, err := c.getStatus()
	if err != nil {
		return 0, err
	}
	return status.SyncInfo.LatestBlockHeight, nil
}

func (c *Client) getStatus() (*ctypes.ResultStatus, error) {
	cc, err := retry(func() (interface{}, error) {
		return c.clientCtx.Client.Status(context.Background())
	})
	if err != nil {
		return nil, err
	}
	return cc.(*ctypes.ResultStatus), nil
}

func (c *Client) GetAccount() (client.Account, error) {
	done := core.UseSdkConfigContext(AccountPrefix)
	defer done()
	return c.getAccount(0, c.clientCtx.FromAddress)
}

func (c *Client) getAccount(height int64, addr types.AccAddress) (client.Account, error) {
	cc, err := retry(func() (interface{}, error) {
		client := c.clientCtx.WithHeight(height)
		return client.AccountRetriever.GetAccount(c.clientCtx, addr)
	})
	if err != nil {
		return nil, err
	}
	return cc.(client.Account), nil
}

func (c *Client) GetTxs(events []string, page, limit int, orderBy string) (*types.SearchTxsResult, error) {
	done := core.UseSdkConfigContext(AccountPrefix)
	defer done()
	cc, err := retry(func() (interface{}, error) {
		return xAuthTx.QueryTxsByEvents(c.clientCtx, events, page, limit, orderBy)
	})
	if err != nil {
		return nil, err
	}
	return cc.(*types.SearchTxsResult), nil
}

func (c *Client) GetBlockTxs(height int64) ([]*types.TxResponse, error) {
	searchTxs, err := c.GetTxs([]string{fmt.Sprintf("tx.height=%d", height)}, 1, 1000, "asc")
	if err != nil {
		return nil, err
	}
	if searchTxs.TotalCount != searchTxs.Count {
		return nil, fmt.Errorf("tx total count overflow, total: %d", searchTxs.TotalCount)
	}
	return searchTxs.GetTxs(), nil
}

func Retry(f func() (interface{}, error)) (interface{}, error) {
	return retry(f)
}

//only retry func when return connection err here
func retry(f func() (interface{}, error)) (interface{}, error) {
	var err error
	var result interface{}
	for i := 0; i < retryLimit; i++ {
		result, err = f()
		if err != nil && isConnectionError(err) {
			time.Sleep(waitTime)
			continue
		}
		return result, err
	}
	panic(fmt.Sprintf("reach retry limit. err: %s", err))
}

func isConnectionError(err error) bool {
	switch t := err.(type) {
	case *url.Error:
		if t.Timeout() || t.Temporary() {
			return true
		}
		return isConnectionError(t.Err)
	}

	switch t := err.(type) {
	case *net.OpError:
		if t.Op == "dial" || t.Op == "read" {
			return true
		}
		return isConnectionError(t.Err)

	case syscall.Errno:
		if t == syscall.ECONNREFUSED {
			return true
		}
	}

	switch t := err.(type) {
	case wrapError:
		newErr := t.Unwrap()
		return isConnectionError(newErr)
	}

	return false
}

type wrapError interface {
	Unwrap() error
}
