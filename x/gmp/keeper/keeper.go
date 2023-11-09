package keeper

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/cometbft/cometbft/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/ojo-network/ojo/app/ibctransfer"
	"github.com/ojo-network/ojo/x/gmp/types"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	storeKey     storetypes.StoreKey
	oracleKeeper types.OracleKeeper
	ibcKeeper    ibctransfer.Keeper
	// the address capable of executing a MsgSetParams message. Typically, this
	// should be the x/gov module account.
	authority string
}

// NewKeeper constructs a new keeper for gmp module.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	oracleKeeper types.OracleKeeper,
	ibcKeeper ibctransfer.Keeper,
	authority string,
) Keeper {
	return Keeper{
		cdc:          cdc,
		storeKey:     storeKey,
		authority:    authority,
		oracleKeeper: oracleKeeper,
		ibcKeeper:    ibcKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// RelayPrice
func (k Keeper) RelayPrice(
	goCtx context.Context,
	msg *types.MsgRelayPrice,
) (*types.MsgRelayPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)

	// encode oracle data
	rates := []types.PriceFeedData{}
	for _, denom := range msg.Denoms {
		rate, err := k.oracleKeeper.GetExchangeRate(ctx, denom)
		if err != nil {
			return &types.MsgRelayPriceResponse{}, err
		}

		priceFeed, err := types.NewPriceFeedData(
			denom,
			rate,
			// TODO: replace with actual resolve time & id
			// Ref: https://github.com/ojo-network/ojo/issues/309
			big.NewInt(1),
			big.NewInt(1),
		)
		if err != nil {
			k.Logger(ctx).With(err).Error("unable to relay price to gmp")
			continue
		}

		rates = append(rates, priceFeed)
	}

	// TODO: fill with actual disableResolve option
	// Ref: https://github.com/ojo-network/ojo/issues/309
	payload, err := types.EncodeABI("postPrices", rates, false)
	if err != nil {
		return nil, err
	}

	// package GMP
	message := types.GmpMessage{
		DestinationChain:   msg.DestinationChain,
		DestinationAddress: msg.DestinationAddress,
		Payload:            payload,
		Type:               types.TypeGeneralMessage,
	}
	bz, err := message.Marshal()
	if err != nil {
		return nil, err
	}

	// submit IBC transfer
	transferMsg := ibctransfertypes.NewMsgTransfer(
		ibctransfertypes.PortID,
		params.GmpChannel,
		msg.Token,
		msg.Relayer,
		params.GmpAddress,
		clienttypes.ZeroHeight(),
		uint64(ctx.BlockTime().Add(time.Duration(params.GmpTimeout)*time.Hour).UnixNano()),
		string(bz),
	)

	_, err = k.ibcKeeper.Transfer(ctx, transferMsg)
	if err != nil {
		return &types.MsgRelayPriceResponse{}, err
	}

	return &types.MsgRelayPriceResponse{}, nil
}
