package keeper

import (
	"context"

	"github.com/ojo-network/ojo/x/airdrop/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the oracle MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (ms msgServer) CreateAirdropAccount(
	ctx context.Context,
	msg *types.MsgCreateAirdropAccount,
) (*types.MsgCreateAirdropAccountResponse, error) {
	panic("not implemented")
}

func (ms msgServer) ClaimAirdrop(
	ctx context.Context,
	msg *types.MsgClaimAirdrop,
) (*types.MsgClaimAirdropResponse, error) {
	panic("not implemented")
}
