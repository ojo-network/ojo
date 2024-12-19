package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/x/symbiotic/types"
)

var _ types.QueryServer = querier{}

// Querier implements a QueryServer for the x/symbiotic module.
type querier struct {
	Keeper
}

// NewQuerier returns an implementation of the symbiotic QueryServer interface
// for the provided Keeper.
func NewQuerier(keeper Keeper) types.QueryServer {
	return &querier{Keeper: keeper}
}

// Params queries params of x/symbiotic module.
func (q querier) Params(
	goCtx context.Context,
	_ *types.ParamsRequest,
) (*types.ParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := q.GetParams(ctx)
	return &types.ParamsResponse{Params: params}, nil
}

// Cached Block Hashes currently in the store.
func (q querier) CachedBlockHashes(
	goCtx context.Context,
	_ *types.CachedBlockHashesRequest,
) (*types.CachedBlockHashesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	cachedBlockHashes := q.GetAllCachedBlockHashes(ctx)
	return &types.CachedBlockHashesResponse{CachedBlockHashes: cachedBlockHashes}, nil
}
