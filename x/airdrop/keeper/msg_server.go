package keeper

import (
	"context"
	"fmt"

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
	// TODO - require authority signature
	ctx := sdk.UnwrapSDKContext(goCtx)
	airdropAccount := types.AirdropAccount{
		OriginAddress:  msg.Address,
		OriginAmount:   msg.TokensToReceive,
		VestingEndTime: msg.VestingEndTime,
	}
	fmt.Println(airdropAccount)
	ms.keeper.SetAirdropAccount(ctx, airdropAccount)
	return &types.MsgCreateAirdropAccountResponse{}, nil
}

// ClaimAirdrop implements MsgServer.ClaimAirdrop method.
// It defines a method to claim an airdrop.
func (ms msgServer) ClaimAirdrop(
	goCtx context.Context,
	msg *types.MsgClaimAirdrop,
) (*types.MsgClaimAirdropResponse, error) {
	// TODO - require signature from claim address
	ctx := sdk.UnwrapSDKContext(goCtx)
	claimAddress := sdk.AccAddress(msg.ToAddress)
	airdropAccount, err := ms.keeper.GetAirdropAccount(ctx, msg.FromAddress)
	if err != nil {
		return nil, err
	}

	// Check if already claimed
	if airdropAccount.ClaimAddress != "" {
		return nil, errors.Wrapf(
			types.ErrAirdropAlreadyClaimed,
			"already claimed by address %s",
			airdropAccount.ClaimAddress,
		)
	}

	// Check if past expiry block
	if ctx.BlockHeight() > int64(ms.keeper.GetParams(ctx).ExpiryBlock) {
		return nil, types.ErrAirdropExpired
	}

	// Check delegation requirement
	delegations := ms.keeper.stakingKeeper.GetDelegatorDelegations(ctx, claimAddress, 999)
	totalShares := sdk.ZeroDec()
	for _, delegation := range delegations {
		totalShares = totalShares.Add(delegation.Shares)
	}
	if totalShares.LT(*ms.keeper.GetParams(ctx).DelegationRequirement) {
		return nil, types.ErrInsufficientDelegation
	}

	// Mint and send tokens
	claimAmount := airdropAccount.OriginAmount.Mul(*ms.keeper.GetParams(ctx).AirdropFactor)
	claimAmount.Abs()

	claimDecCoin := sdk.NewCoins(sdk.NewCoin("ojo", claimAmount.TruncateInt()))
	ms.keeper.bankKeeper.MintCoins(ctx, types.ModuleName, claimDecCoin)
	ms.keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, claimAddress, claimDecCoin)

	// Update AirdropAccount
	airdropAccount.ClaimAddress = msg.ToAddress
	airdropAccount.ClaimAmount = &claimAmount
	ms.keeper.SetAirdropAccount(ctx, airdropAccount)

	return &types.MsgClaimAirdropResponse{}, nil
}
