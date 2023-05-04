package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/x/airdrop/types"
)

var _ types.QueryServer = querier{}

// Querier implements a QueryServer for the x/airdrop module.
type querier struct {
	Keeper
}

// NewQuerier returns an implementation of the airdrop QueryServer interface
// for the provided Keeper.
func NewQuerier(keeper Keeper) types.QueryServer {
	return &querier{Keeper: keeper}
}

// Params queries params of x/airdrop module.
func (q querier) Params(
	goCtx context.Context,
	req *types.ParamsRequest,
) (*types.ParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params, err := q.GetParams(ctx)
	if err != nil {
		return nil, err
	}
	return &types.ParamsResponse{Params: params}, nil
}

func (q querier) AirdropAccount(
	goCtx context.Context,
	req *types.AirdropAccountRequest,
) (*types.AirdropAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	airdropAccount, err := q.GetAirdropAccount(ctx, req.Address)
	if err != nil {
		return nil, err
	}
	return &types.AirdropAccountResponse{AirdropAccount: &airdropAccount}, nil
}
