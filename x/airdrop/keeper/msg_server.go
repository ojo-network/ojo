package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"

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
	// TODO - require genesis signature

	// Create Continuous Vesting Account
	baseAccount := &authtypes.BaseAccount{
		Address: msg.Address,
	}
	authvesting.NewContinuousVestingAccount(
		baseAccount,
		sdk.NewCoins(sdk.NewCoin("ojo", sdk.NewIntFromUint64(msg.TokensToReceive))),
		msg.VestingStartTime,
		msg.VestingEndTime,
	)

	// Create AirdropAccount entry
	ctx := sdk.UnwrapSDKContext(goCtx)
	airdropAccount := types.AirdropAccount{
		OriginAddress:    msg.Address,
		OriginAmount:     msg.TokensToReceive,
		VestingStartTime: msg.VestingStartTime,
		VestingEndTime:   msg.VestingEndTime,
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
	claimAmount := ms.keeper.GetParams(ctx).AirdropFactor.MulInt64(int64(airdropAccount.OriginAmount))
	claimAmount.Abs()
	intClaimAmount := claimAmount.TruncateInt64()

	claimDecCoin := sdk.NewCoins(sdk.NewCoin("ojo", claimAmount.TruncateInt()))
	ms.keeper.bankKeeper.MintCoins(ctx, types.ModuleName, claimDecCoin)
	ms.keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, claimAddress, claimDecCoin)

	// Create delayed vesting account
	baseAccount := &authtypes.BaseAccount{
		Address: msg.ToAddress,
	}
	authvesting.NewDelayedVestingAccount(
		baseAccount,
		sdk.NewCoins(sdk.NewCoin("ojo", claimAmount.TruncateInt())),
		airdropAccount.VestingEndTime,
	)

	// Update AirdropAccount
	airdropAccount.ClaimAddress = msg.ToAddress
	airdropAccount.ClaimAmount = uint64(intClaimAmount)
	ms.keeper.SetAirdropAccount(ctx, airdropAccount)

	return &types.MsgClaimAirdropResponse{}, nil
}
