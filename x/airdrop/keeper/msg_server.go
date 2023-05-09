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
// It defines a method to update the x/airdrop module parameters.
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
	ms.keeper.SetParams(ctx, *msg.Params)

	return &types.MsgSetParamsResponse{}, nil
}

// CreateAirdropAccount implements MsgServer.CreateAirdropAccount method.
// It defines a method to create an airdrop account.
func (ms msgServer) CreateAirdropAccount(
	goCtx context.Context,
	msg *types.MsgCreateAirdropAccount,
) (*types.MsgCreateAirdropAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	airdropAccount := &types.AirdropAccount{
		OriginAddress:  msg.Address,
		OriginAmount:   msg.TokensToReceive,
		VestingEndTime: msg.VestingEndTime,
	}

	ms.keeper.CreateOriginAccount(ctx, airdropAccount)
	err := ms.keeper.MintOriginTokens(ctx, airdropAccount)
	if err != nil {
		return nil, err
	}
	err = ms.keeper.SendOriginTokens(ctx, airdropAccount)
	if err != nil {
		return nil, err
	}
	ms.keeper.SetAirdropAccount(ctx, airdropAccount)

	return &types.MsgCreateAirdropAccountResponse{}, nil
}

// ClaimAirdrop implements MsgServer.ClaimAirdrop method.
// It defines a method to claim an airdrop.
func (ms msgServer) ClaimAirdrop(
	goCtx context.Context,
	msg *types.MsgClaimAirdrop,
) (*types.MsgClaimAirdropResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	airdropAccount, err := ms.keeper.GetAirdropAccount(ctx, msg.FromAddress)
	if err != nil {
		return nil, err
	}

	if err := airdropAccount.VerifyNotClaimed(); err != nil {
		return nil, err
	}

	// Check if past expiry block
	if ctx.BlockHeight() > int64(ms.keeper.GetParams(ctx).ExpiryBlock) {
		return nil, types.ErrAirdropExpired
	}

	airdropAccount.ClaimAddress = msg.ToAddress
	if err := ms.keeper.VerifyDelegationRequirement(ctx, airdropAccount); err != nil {
		return nil, err
	}

	ms.keeper.SetClaimAmount(ctx, airdropAccount)
	if err = ms.keeper.MintClaimTokensToAirdrop(ctx, airdropAccount); err != nil {
		return nil, err
	}
	ms.keeper.CreateClaimAccount(ctx, airdropAccount)
	if err = ms.keeper.SendClaimTokens(ctx, airdropAccount); err != nil {
		return nil, err
	}
	ms.keeper.SetAirdropAccount(ctx, airdropAccount)

	return &types.MsgClaimAirdropResponse{}, nil
}
