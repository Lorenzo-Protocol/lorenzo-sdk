package query

import (
	"context"

	"github.com/Lorenzo-Protocol/lorenzo/v3/x/btcstaking/types"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
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
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *QueryClient) GetBTCStakingRecord(txHash string) (*types.QueryStakingRecordResponse, error) {
	txHashBytes, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return nil, err
	}

	var resp *types.QueryStakingRecordResponse
	err = c.QueryBTCStaking(func(ctx context.Context, queryClient types.QueryClient) error {
		req := &types.QueryStakingRecordRequest{
			TxHash: txHashBytes[:],
		}

		var err error
		resp, err = queryClient.StakingRecord(ctx, req)
		return err
	})

	return resp, err
}
