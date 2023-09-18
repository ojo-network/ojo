package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ojo-network/ojo/x/relayoracle/types"
)

var _ types.QueryServer = querier{}

// Querier implements a QueryServer for the x/relayoracle module.
type querier struct {
	Keeper
}

// NewQuerier returns an implementation of the relayoracle QueryServer interface
// for the provided Keeper.
func NewQuerier(keeper Keeper) types.QueryServer {
	return &querier{Keeper: keeper}
}

func (k Keeper) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request.go")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}
