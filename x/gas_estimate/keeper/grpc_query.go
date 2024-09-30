package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/x/gas_estimate/types"
)

var _ types.QueryServer = querier{}

// Querier implements a QueryServer for the x/gas_estimate module.
type querier struct {
	Keeper
}

// NewQuerier returns an implementation of the gas_estimate QueryServer interface
// for the provided Keeper.
func NewQuerier(keeper Keeper) types.QueryServer {
	return &querier{Keeper: keeper}
}

// Params queries params of x/gas_estimate module.
func (q querier) Params(
	goCtx context.Context,
	_ *types.ParamsRequest,
) (*types.ParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := q.GetParams(ctx)
	return &types.ParamsResponse{Params: params}, nil
}

func (q querier) GasEstimate(
	goCtx context.Context,
	req *types.GasEstimateRequest,
) (*types.GasEstimateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	gasEstimate, err := q.GetGasEstimate(ctx, req.Network)
	if err != nil {
		return nil, err
	}
	return &types.GasEstimateResponse{GasEstimate: gasEstimate.GasEstimate}, nil
}
