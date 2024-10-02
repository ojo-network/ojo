package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ojo-network/ojo/x/gmp/types"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
)

type msgServer struct {
	keeper Keeper
}

// NewMsgServerImpl returns an implementation of the gmp MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{keeper: keeper}
}

// SetParams implements MsgServer.SetParams method.
// It defines a method to update the x/gmp module parameters.
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

// Relay implements MsgServer.Relay method.
// It defines a method to relay over GMP to recipient chains.
func (ms msgServer) RelayPrice(
	goCtx context.Context,
	msg *types.MsgRelayPrice,
) (*types.MsgRelayPriceResponse, error) {
	return ms.keeper.RelayPrice(goCtx, msg)
}

func (ms msgServer) CreatePayment(
	goCtx context.Context,
	msg *types.MsgCreatePayment,
) (*types.MsgCreatePaymentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// make sure the destination chain is valid
	gasEstimateParams := ms.keeper.GasEstimateKeeper.GetParams(ctx)
	isValidChain := false
	for _, chain := range gasEstimateParams.ContractRegistry {
		if chain.Network == msg.Payment.DestinationChain {
			isValidChain = true
			break
		}
	}
	if !isValidChain {
		return nil, errors.Wrapf(
			types.ErrInvalidDestinationChain,
			"destination chain %s not found in contract registry",
			msg.Payment.DestinationChain,
		)
	}

	// make sure the denom is active in the oracle
	_, err := ms.keeper.oracleKeeper.GetExchangeRate(ctx, msg.Payment.Denom)
	if err != nil {
		return nil, errors.Wrapf(
			oracletypes.ErrUnknownDenom,
			"denom %s not active in the oracle",
			msg.Payment.Denom,
		)
	}

	// send tokens from msg to the module account
	address, err := sdk.AccAddressFromBech32(msg.Relayer)
	if err != nil {
		return nil, err
	}
	coins := sdk.NewCoins(msg.Payment.Token)
	err = ms.keeper.BankKeeper.SendCoinsFromAccountToModule(ctx, address, types.ModuleName, coins)
	if err != nil {
		return nil, err
	}

	// Create a payment record in the KV store
	msg.Payment.Relayer = msg.Relayer
	ms.keeper.SetPayment(ctx, *msg.Payment)
	return &types.MsgCreatePaymentResponse{}, nil
}
