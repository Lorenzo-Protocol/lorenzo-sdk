package query

import (
	"context"
	"encoding/hex"
	"github.com/Lorenzo-Protocol/lorenzo/x/btcstaking/types"
	"github.com/cosmos/cosmos-sdk/client"
)

func (c *QueryClient) QueryBTCStaking(f func(ctx context.Context, queryClient types.QueryClient) error) error {
	ctx, cancel := c.getQueryContext()
	defer cancel()

	clientCtx := client.Context{Client: c.RPCClient}
	queryClient := types.NewQueryClient(clientCtx)

	return f(ctx, queryClient)
}

func (c *QueryClient) QueryBTCStakingParams() (*types.QueryParamsResponse, error) {
	var resp *types.QueryParamsResponse
	err := c.QueryBTCStaking(func(ctx context.Context, queryClient types.QueryClient) error {
		req := &types.QueryParamsRequest{}

		var err error
		resp, err = queryClient.Params(ctx, req)
		return err
	})

	return resp, err
}

func (c *QueryClient) GetBTCStakingRecord(txHash string) (*types.QueryStakingRecordResponse, error) {
	txHashBytes, err := hex.DecodeString(txHash)
	if err != nil {
		return nil, err
	}

	var resp *types.QueryStakingRecordResponse
	err = c.QueryBTCStaking(func(ctx context.Context, queryClient types.QueryClient) error {
		req := &types.QueryStakingRecordRequest{
			TxHash: txHashBytes,
		}

		var err error
		resp, err = queryClient.StakingRecord(ctx, req)
		return err
	})

	return resp, err
}
