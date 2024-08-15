package query

import (
	"context"

	"github.com/Lorenzo-Protocol/lorenzo/v2/x/token/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/query"
)

func (c *QueryClient) QueryToken(f func(ctx context.Context, client types.QueryClient) error) error {
	ctx, cancel := c.getQueryContext()
	defer cancel()

	clientCtx := client.Context{Client: c.RPCClient}
	queryClient := types.NewQueryClient(clientCtx)

	return f(ctx, queryClient)
}

func (c *QueryClient) TokenPairs(pageRequest *query.PageRequest) (*types.QueryTokenPairsResponse, error) {
	var resp *types.QueryTokenPairsResponse
	err := c.QueryToken(func(ctx context.Context, queryClient types.QueryClient) error {
		var err error
		resp, err = queryClient.TokenPairs(ctx, &types.QueryTokenPairsRequest{
			Pagination: pageRequest,
		})
		return err
	})

	return resp, err
}

func (c *QueryClient) TokenPair(tokenAddressOrDenom string) (*types.QueryTokenPairResponse, error) {
	var resp *types.QueryTokenPairResponse
	err := c.QueryToken(func(ctx context.Context, queryClient types.QueryClient) error {
		var err error
		resp, err = queryClient.TokenPair(ctx, &types.QueryTokenPairRequest{
			Token: tokenAddressOrDenom,
		})
		return err
	})

	return resp, err
}

func (c *QueryClient) Balance(accountAddress string, tokenAddressOrDenom string) (*types.QueryBalanceResponse, error) {
	var resp *types.QueryBalanceResponse
	err := c.QueryToken(func(ctx context.Context, queryClient types.QueryClient) error {
		var err error
		resp, err = queryClient.Balance(ctx, &types.QueryBalanceRequest{
			AccountAddress: accountAddress,
			Token:          tokenAddressOrDenom,
		})
		return err
	})

	return resp, err
}

func (c *QueryClient) TokenParams() (*types.QueryParamsResponse, error) {
	var resp *types.QueryParamsResponse
	err := c.QueryToken(func(ctx context.Context, queryClient types.QueryClient) error {
		var err error
		resp, err = queryClient.Params(ctx, &types.QueryParamsRequest{})
		return err
	})

	return resp, err
}
