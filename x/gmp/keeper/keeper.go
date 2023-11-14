package keeper

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/ethereum/go-ethereum/common"

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

	prices := []types.PriceData{}
	for _, denom := range msg.Denoms {
		// get exchange rate
		rate, err := k.oracleKeeper.GetExchangeRate(ctx, denom)
		if err != nil {
			return &types.MsgRelayPriceResponse{}, err
		}

		// get any available median and standard deviation data
		medians := k.oracleKeeper.HistoricMedians(
			ctx,
			denom,
			k.oracleKeeper.MaximumMedianStamps(ctx),
		)
		deviations := k.oracleKeeper.HistoricDeviations(
			ctx,
			denom,
			k.oracleKeeper.MaximumMedianStamps(ctx),
		)
		// convert them to a medianData slice
		medianData, err := types.NewMediansSlice(medians, deviations)
		if err != nil {
			return &types.MsgRelayPriceResponse{}, err
		}

		priceFeed, err := types.NewPriceData(
			denom,
			rate,
			big.NewInt(msg.Timestamp),
			medianData,
		)
		if err != nil {
			k.Logger(ctx).With(err).Error("unable to relay price to gmp")
			continue
		}

		prices = append(prices, priceFeed)
	}

	// convert commandSelector to [4]byte
	var commandSelector [4]byte
	copy(commandSelector[:], msg.CommandSelector)

	encoder := types.NewGMPEncoder(
		prices,
		msg.Denoms,
		common.HexToAddress(msg.ContractAddress),
		commandSelector,
		msg.CommandParams,
	)
	payload, err := encoder.GMPEncode()
	if err != nil {
		return nil, err
	}

	// package GMP
	message := types.GmpMessage{
		DestinationChain:   msg.DestinationChain,
		DestinationAddress: msg.ContractAddress,
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
