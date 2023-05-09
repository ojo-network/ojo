package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ojo-network/ojo/x/airdrop/types"
)

type msgServer struct {
	keeper Keeper
}

// NewMsgServerImpl returns an implementation of the airdrop MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{keeper: keeper}
}

// SetParams implements MsgServer.SetParams method.
// It defines a method to update the x/slashing module parameters.
func (ms msgServer) SetParams(goCtx context.Context, msg *types.MsgSetParams) (*types.MsgSetParamsResponse, error) {
	if ms.keeper.authority != msg.Authority {
		err := errors.Wrapf(
			govtypes.ErrInvalidSigner,
			"invalid authority; expected %s, got %s",
			ms.keeper.authority,
			msg.Authority,
		)
		return nil, err
	}

	if err := msg.Params.Validate(); err != nil {
		return nil, err
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := ms.keeper.SetParams(ctx, *msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgSetParamsResponse{}, nil
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
