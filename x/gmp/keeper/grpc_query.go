package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/x/gmp/types"
)

var _ types.QueryServer = querier{}

// Querier implements a QueryServer for the x/gmp module.
type querier struct {
	Keeper
}

// NewQuerier returns an implementation of the gmp QueryServer interface
// for the provided Keeper.
func NewQuerier(keeper Keeper) types.QueryServer {
	return &querier{Keeper: keeper}
}

// Params queries params of x/gmp module.
func (q querier) Params(
	goCtx context.Context,
	_ *types.ParamsRequest,
) (*types.ParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := q.GetParams(ctx)
	return &types.ParamsResponse{Params: params}, nil
}

func (q querier) AllPayments(
	goCtx context.Context,
	_ *types.AllPaymentsRequest,
) (*types.AllPaymentsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	payments := q.Keeper.GetAllPayments(ctx)
	return &types.AllPaymentsResponse{Payments: payments}, nil
}
