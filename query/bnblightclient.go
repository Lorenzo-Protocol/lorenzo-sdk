package query

import (
	"context"
	"github.com/ethereum/go-ethereum/common"

	"github.com/Lorenzo-Protocol/lorenzo/v2/x/bnblightclient/types"
	"github.com/cosmos/cosmos-sdk/client"
)

func (c *QueryClient) QueryBNBLightClient(f func(ctx context.Context, client types.QueryClient) error) error {
	ctx, cancel := c.getQueryContext()
	defer cancel()

	clientCtx := client.Context{Client: c.RPCClient}
	queryClient := types.NewQueryClient(clientCtx)

	return f(ctx, queryClient)
}

func (c *QueryClient) BNBHeader(number uint64) (*types.Header, error) {
	var resp *types.QueryHeaderResponse
	err := c.QueryBNBLightClient(func(ctx context.Context, queryClient types.QueryClient) error {
		var err error
		req := &types.QueryHeaderRequest{
			Number: number,
		}
		resp, err = queryClient.Header(ctx, req)
		return err
	})

	if err != nil {
		return nil, err
	}

	return resp.Header, err
}

func (c *QueryClient) BNBHeaderByHash(hash string) (*types.Header, error) {
	var resp *types.QueryHeaderByHashResponse
	err := c.QueryBNBLightClient(func(ctx context.Context, queryClient types.QueryClient) error {
		var err error
		req := &types.QueryHeaderByHashRequest{
			Hash: common.FromHex(hash),
		}
		resp, err = queryClient.HeaderByHash(ctx, req)
		return err
	})

	if err != nil {
		return nil, err
	}

	return resp.Header, err
}

func (c *QueryClient) BNBLatestHeader() (*types.Header, error) {
	var resp *types.QueryLatestHeaderResponse
	err := c.QueryBNBLightClient(func(ctx context.Context, queryClient types.QueryClient) error {
		var err error
		req := &types.QueryLatestHeaderRequest{}
		resp, err = queryClient.LatestHeader(ctx, req)
		return err
	})

	if err != nil {
		return nil, err
	}

	return &resp.Header, err
}

func (c *QueryClient) BNBLightClientParams() (*types.QueryParamsResponse, error) {
	var resp *types.QueryParamsResponse
	err := c.QueryBNBLightClient(func(ctx context.Context, queryClient types.QueryClient) error {
		req := &types.QueryParamsRequest{}

		var err error
		resp, err = queryClient.Params(ctx, req)
		if err != nil {
			return err
		}

		return nil
	})

	return resp, err
}
