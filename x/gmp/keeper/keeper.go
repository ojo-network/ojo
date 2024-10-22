package keeper

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"cosmossdk.io/errors"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/ojo-network/ojo/app/ibctransfer"
	"github.com/ojo-network/ojo/util"
	"github.com/ojo-network/ojo/x/gmp/types"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
)

type Keeper struct {
	cdc               codec.BinaryCodec
	storeKey          storetypes.StoreKey
	oracleKeeper      types.OracleKeeper
	IBCKeeper         *ibctransfer.Keeper
	BankKeeper        types.BankKeeper
	GasEstimateKeeper types.GasEstimateKeeper
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
	bankKeeper types.BankKeeper,
	gasEstimateKeeper types.GasEstimateKeeper,
) Keeper {
	return Keeper{
		cdc:               cdc,
		storeKey:          storeKey,
		authority:         authority,
		oracleKeeper:      oracleKeeper,
		IBCKeeper:         &ibcKeeper,
		BankKeeper:        bankKeeper,
		GasEstimateKeeper: gasEstimateKeeper,
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
		util.SafeInt64ToUint64(ctx.BlockTime().Add(time.Duration(params.GmpTimeout)*time.Hour).UnixNano()),
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
		// convert them to a MedianData struct
		medianData, err := types.NewMedianData(medians, deviations)
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
		Fee: &types.GmpFee{
			Amount:    msg.Token.Amount.String(),
			Recipient: params.FeeRecipient,
		},
	}
	bz, err := json.Marshal(&message)
	if err != nil {
		k.Logger(ctx).With(err).Error("error marshaling GMP message")
		return nil, nil
	}

	// submit IBC transfer
	transferMsg := ibctransfertypes.NewMsgTransfer(
		ibctransfertypes.PortID,
		params.GmpChannel,
		msg.Token,
		msg.Relayer,
		params.GmpAddress,
		clienttypes.ZeroHeight(),
		util.SafeInt64ToUint64(ctx.BlockTime().Add(time.Duration(params.GmpTimeout)*time.Hour).UnixNano()),
		string(bz),
	)
	return transferMsg, nil
}

func (k Keeper) SetPayment(
	ctx sdk.Context,
	payment types.Payment,
) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshal(&payment)
	store.Set(types.PaymentKey(payment.Relayer, payment.Denom), bz)
}

func (k Keeper) DeletePayment(
	ctx sdk.Context,
	payment types.Payment,
) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.PaymentKey(payment.Relayer, payment.Denom))
}

func (k Keeper) GetPayment(
	ctx sdk.Context,
	authority string,
	denom string,
) (types.MsgCreatePayment, error) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.PaymentKey(authority, denom))
	payment := types.MsgCreatePayment{}
	k.cdc.MustUnmarshal(bz, &payment)
	return payment, nil
}

func (k Keeper) GetAllPayments(
	ctx sdk.Context,
) []types.Payment {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.PaymentKeyPrefix)
	defer iterator.Close()

	payments := []types.Payment{}
	for ; iterator.Valid(); iterator.Next() {
		payment := types.Payment{}
		k.cdc.MustUnmarshal(iterator.Value(), &payment)
		payments = append(payments, payment)
	}
	return payments
}

func (k Keeper) ProcessPayment(
	goCtx context.Context,
	payment types.Payment,
) error {
	ctx := sdk.UnwrapSDKContext(goCtx)
	k.Logger(ctx).Info("processing payment", "payment", payment)
	gasEstimateParams := k.GasEstimateKeeper.GetParams(ctx)
	// get gmp contract address on receiving chain
	contractAddress := ""
	for _, contract := range gasEstimateParams.ContractRegistry {
		if payment.DestinationChain == contract.Network {
			contractAddress = contract.Address
		}
	}
	// if contract address not found, return
	if contractAddress == "" {
		k.Logger(ctx).Error("contract address not found for chain", "chain", payment.DestinationChain)
		return fmt.Errorf("contract address not found for chain %s", payment.DestinationChain)
	}

	gasAmount := math.NewInt(k.GetParams(ctx).DefaultGasEstimate)
	gasEstimate, err := k.GasEstimateKeeper.GetGasEstimate(ctx, payment.DestinationChain)
	if err != nil {
		k.Logger(ctx).With(err).Error("error getting gas estimate. using default gas estimates")
	} else {
		gasAmount = math.NewInt(gasEstimate.GasEstimate)
	}

	coins := sdk.Coin{
		Denom:  payment.Token.Denom,
		Amount: gasAmount,
	}

	// if payment.Token.Amount is less than coins.Amount, return funds and delete payment
	relayerAddr, err := sdk.AccAddressFromBech32(payment.Relayer)
	if err != nil {
		k.Logger(ctx).With(err).Error("error getting relayer address")
		return err
	}
	if payment.Token.Amount.LTE(coins.Amount) {
		k.Logger(ctx).With(err).Debug("payment amount is less than gas estimate, returning funds and deleting payment")
		err := k.BankKeeper.SendCoinsFromModuleToAccount(
			ctx,
			types.ModuleName,
			relayerAddr,
			sdk.NewCoins(payment.Token),
		)
		if err != nil {
			return err
		}
		k.DeletePayment(ctx, payment)
		return nil
	}

	msg := types.NewMsgRelay(
		authtypes.NewModuleAddress(types.ModuleName).String(),
		payment.DestinationChain,
		contractAddress,
		types.EmptyContract,
		coins,
		[]string{payment.Denom},
		types.EmptyByteSlice,
		types.EmptyByteSlice,
		ctx.BlockTime().Unix(),
	)
	_, err = k.RelayPrice(goCtx, msg)
	if err != nil {
		k.Logger(ctx).With(err).Error("error relaying price")
		return err
	}
	k.Logger(ctx).Info("relay price submitted", "MsgRelayPrice", msg)

	// update payment in the store with the amount paid
	payment.Token.Amount = payment.Token.Amount.Sub(coins.Amount)
	lastPrice, err := k.oracleKeeper.GetExchangeRate(ctx, payment.Denom)
	if err != nil {
		k.Logger(ctx).With(err).Error("error getting exchange rate")
		return err
	}
	payment.LastPrice = lastPrice
	payment.LastBlock = ctx.BlockHeight()

	k.SetPayment(ctx, payment)
	k.Logger(ctx).Info("payment updated", "payment", payment)

	return nil
}

func (k Keeper) CreatePayment(
	goCtx context.Context,
	msg *types.MsgCreatePayment,
) (*types.MsgCreatePaymentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// make sure the destination chain is valid
	gasEstimateParams := k.GasEstimateKeeper.GetParams(ctx)
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
	_, err := k.oracleKeeper.GetExchangeRate(ctx, msg.Payment.Denom)
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
	err = k.BankKeeper.SendCoinsFromAccountToModule(ctx, address, types.ModuleName, coins)
	if err != nil {
		return nil, err
	}

	// Create a payment record in the KV store
	msg.Payment.Relayer = msg.Relayer
	k.SetPayment(ctx, *msg.Payment)
	return &types.MsgCreatePaymentResponse{}, nil
}
