package query

import (
	"context"
	"cosmossdk.io/math"

	plantypes "github.com/Lorenzo-Protocol/lorenzo/v3/x/plan/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/query"
)

func (c *QueryClient) QueryPlan(f func(ctx context.Context, queryClient plantypes.QueryClient) error) error {
	ctx, cancel := c.getQueryContext()
	defer cancel()

	clientCtx := client.Context{Client: c.RPCClient}
	queryClient := plantypes.NewQueryClient(clientCtx)

	return f(ctx, queryClient)
}

func (c *QueryClient) PlanParams() (*plantypes.QueryParamsResponse, error) {
	var resp *plantypes.QueryParamsResponse
	err := c.QueryPlan(func(ctx context.Context, queryClient plantypes.QueryClient) error {
		req := &plantypes.QueryParamsRequest{}

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

func (c *QueryClient) Plans(pageRequest *query.PageRequest) (*plantypes.QueryPlansResponse, error) {
	var resp *plantypes.QueryPlansResponse
	err := c.QueryPlan(func(ctx context.Context, queryClient plantypes.QueryClient) error {
		req := &plantypes.QueryPlansRequest{
			Pagination: pageRequest,
		}

		var err error
		resp, err = queryClient.Plans(ctx, req)
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

func (c *QueryClient) Plan(planId uint64) (plantypes.Plan, error) {
	var resp *plantypes.QueryPlanResponse
	err := c.QueryPlan(func(ctx context.Context, queryClient plantypes.QueryClient) error {
		var err error
		resp, err = queryClient.Plan(ctx, &plantypes.QueryPlanRequest{
			Id: planId,
		})
		return err
	})

	if err != nil {
		return plantypes.Plan{}, err
	}

	return resp.Plan, nil
}

func (c *QueryClient) ClaimLeafNode(planId uint64, roundId math.Int, leafNode string) (bool, error) {
	var resp *plantypes.QueryClaimLeafNodeResponse
	err := c.QueryPlan(func(ctx context.Context, queryClient plantypes.QueryClient) error {
		var err error
		resp, err = queryClient.ClaimLeafNode(ctx, &plantypes.QueryClaimLeafNodeRequest{
			Id:       planId,
			RoundId:  roundId,
			LeafNode: leafNode,
		})
		return err
	})

	if err != nil {
		return false, err
	}

	return resp.Success, nil
}
