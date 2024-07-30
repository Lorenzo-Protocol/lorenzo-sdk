package query

import (
	"context"

	"github.com/Lorenzo-Protocol/lorenzo/v2/x/agent/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/query"
)

func (c *QueryClient) QueryAgent(f func(ctx context.Context, client types.QueryClient) error) error {
	ctx, cancel := c.getQueryContext()
	defer cancel()

	clientCtx := client.Context{Client: c.RPCClient}
	queryClient := types.NewQueryClient(clientCtx)

	return f(ctx, queryClient)
}

func (c *QueryClient) Agents(pageRequest *query.PageRequest) (*types.QueryAgentsResponse, error) {
	var resp *types.QueryAgentsResponse
	err := c.QueryAgent(func(ctx context.Context, queryClient types.QueryClient) error {
		var err error
		resp, err = queryClient.Agents(ctx, &types.QueryAgentsRequest{
			Pagination: pageRequest,
		})
		return err
	})

	return resp, err
}

func (c *QueryClient) Agent(agentId uint64) (*types.QueryAgentResponse, error) {
	var resp *types.QueryAgentResponse
	err := c.QueryAgent(func(ctx context.Context, queryClient types.QueryClient) error {
		var err error
		resp, err = queryClient.Agent(ctx, &types.QueryAgentRequest{
			Id: agentId,
		})
		return err
	})

	return resp, err
}
