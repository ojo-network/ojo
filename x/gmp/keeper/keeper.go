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
	IBCKeeper    *ibctransfer.Keeper
	// the address capable of executing a MsgSetParams message. Typically, this
	// should be the x/gov module account.
	authority string
}

// NewKeeper constructs a new keeper for gmp module.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	oracleKeeper types.OracleKeeper,
	authority string,
	ibcKeeper ibctransfer.Keeper,
) Keeper {
	return Keeper{
		cdc:          cdc,
		storeKey:     storeKey,
		authority:    authority,
		oracleKeeper: oracleKeeper,
		IBCKeeper:    &ibcKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// RelayPrice submits an IBC transfer with a MsgRelayPrice payload.
// This is so the IBC Transfer module can then use BuildGmpRequest
// and perform the GMP request.
func (k Keeper) RelayPrice(
	goCtx context.Context,
	msg *types.MsgRelayPrice,
) (*types.MsgRelayPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)

	bz, err := msg.Marshal()
	if err != nil {
		return nil, err
	}

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
	_, err = k.IBCKeeper.Transfer(ctx, transferMsg)
	if err != nil {
		return &types.MsgRelayPriceResponse{}, err
	}

	return &types.MsgRelayPriceResponse{}, nil
}

func (k Keeper) BuildGmpRequest(
	goCtx context.Context,
	msg *types.MsgRelayPrice,
) (*ibctransfertypes.MsgTransfer, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)

	prices := []types.PriceData{}
	for _, denom := range msg.Denoms {
		// get exchange rate
		rate, err := k.oracleKeeper.GetExchangeRate(ctx, denom)
		if err != nil {
			k.Logger(ctx).With(err).Error("attempting to relay unavailable denom")
			continue
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
			return &ibctransfertypes.MsgTransfer{}, err
		}

		price, err := types.NewPriceData(
			denom,
			rate,
			big.NewInt(msg.Timestamp),
			medianData,
		)
		if err != nil {
			k.Logger(ctx).With(err).Error("unable to relay price to gmp")
			continue
		}
		prices = append(prices, price)
	}

	// convert commandSelector to [4]byte
	var commandSelector [4]byte
	copy(commandSelector[:], msg.CommandSelector)

	encoder := types.NewGMPEncoder(
		prices,
		msg.Denoms,
		common.HexToAddress(msg.ClientContractAddress),
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
		DestinationAddress: msg.OjoContractAddress,
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
	return transferMsg, nil
}
