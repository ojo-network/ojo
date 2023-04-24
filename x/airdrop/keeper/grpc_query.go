package keeper

import (
	"context"

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
	ctx context.Context,
	req *types.ParamsRequest,
) (*types.ParamsResponse, error) {
	panic("not implemented")
}

func (q querier) AirdropAccount(
	ctx context.Context,
	req *types.AirdropAccountRequest,
) (*types.AirdropAccountResponse, error) {
	panic("not implemented")
}
