package keeper

import (
	"fmt"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"

	oracleTypes "github.com/ojo-network/ojo/x/oracle/types"
	"github.com/ojo-network/ojo/x/relayoracle/types"
)

func (k Keeper) AddRequest(ctx sdk.Context, req types.Request) uint64 {
	id := k.GetNextRequestID(ctx)
	k.SetRequest(ctx, id, req)
	k.AddRequestIDToPendingList(ctx, id)
	return id
}

func (k Keeper) GetNextRequestID(ctx sdk.Context) uint64 {
	requestNumber := k.GetRequestCount(ctx)
	k.SetRequestCount(ctx, requestNumber+1)
	return requestNumber + 1
}

func (k Keeper) GetRequestCount(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	return sdk.BigEndianToUint64(store.Get(types.RequestCountKey))
}

func (k Keeper) SetRequestCount(ctx sdk.Context, count uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.RequestCountKey, sdk.Uint64ToBigEndian(count))
}

func (k Keeper) SetRequest(ctx sdk.Context, id uint64, request types.Request) {
	ctx.KVStore(k.storeKey).Set(types.RequestStoreKey(id), k.cdc.MustMarshal(&request))
}

func (k Keeper) GetRequest(ctx sdk.Context, id uint64) (types.Request, error) {
	var request types.Request
	bz := ctx.KVStore(k.storeKey).Get(types.RequestStoreKey(id))
	if bz == nil {
		return request, sdkerrors.Wrapf(types.ErrRequestNotFound, "id %d", id)
	}

	k.cdc.MustUnmarshal(bz, &request)

	return request, nil
}

func (k Keeper) SetResult(ctx sdk.Context, result types.Result) {
	ctx.KVStore(k.storeKey).Set(types.RequestStoreKey(result.RequestID), k.cdc.MustMarshal(&result))
}

// DeleteRequest removes the given data request from the store.
func (k Keeper) DeleteRequest(ctx sdk.Context, id uint64) {
	ctx.KVStore(k.storeKey).Delete(types.RequestStoreKey(id))
}

func (k Keeper) PrepareRequest(
	ctx sdk.Context,
	ibcChannel *types.IBCChannel,
	data *types.OracleRequestPacketData,
) (uint64, error) {
	req := types.NewRequest(data.GetCalldata(), data.GetClientID(), ibcChannel)
	return k.AddRequest(ctx, req), nil
}

func (k Keeper) GetPendingRequestList(ctx sdk.Context) []uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PendingRequestListKey)
	if len(bz) == 0 {
		return []uint64{}
	}

	var pending types.PendingRequestList
	k.cdc.MustUnmarshal(bz, &pending)

	return pending.RequestIds
}

func (k Keeper) SetPendingRequestList(ctx sdk.Context, reqIDS []uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PendingRequestListKey, k.cdc.MustMarshal(&types.PendingRequestList{RequestIds: reqIDS}))
}

func (k Keeper) AddRequestIDToPendingList(ctx sdk.Context, reqID uint64) {
	var pending types.PendingRequestList
	k.cdc.MustUnmarshal(ctx.KVStore(k.storeKey).Get(types.PendingRequestListKey), &pending)
	pending.RequestIds = append(pending.RequestIds, reqID)

	ctx.KVStore(k.storeKey).Set(types.PendingRequestListKey, k.cdc.MustMarshal(&pending))
}

func (k Keeper) FlushPendingRequestList(ctx sdk.Context) {
	ctx.KVStore(k.storeKey).Delete(types.PendingRequestListKey)
}

func (k Keeper) ResolveRequest(ctx sdk.Context, reqID uint64) {
	req, err := k.GetRequest(ctx, reqID)
	if err != nil {
		panic(err)
	}

	result, status := k.ProcessRequestCalldata(ctx, req.GetRequestCallData())
	k.ProcessResult(ctx, reqID, status, result)

	err = ctx.EventManager().EmitTypedEvents(&types.EventRequestResolve{
		RequestId: reqID,
		Status:    status,
	})

	if err != nil {
		panic(err)
	}
}

func (k Keeper) ProcessResult(ctx sdk.Context, requestID uint64, status types.ResolveStatus, result []byte) {
	req, err := k.GetRequest(ctx, requestID)
	if err != nil {
		panic(err)
	}

	k.SetResult(ctx, types.Result{
		RequestID:       requestID,
		RequestCallData: req.RequestCallData,
		ClientID:        req.ClientID,
		RequestHeight:   req.RequestHeight,
		RequestTime:     req.RequestHeight,
		Status:          status,
		Result:          result,
	})

	expiry := k.PacketTimeout(ctx)
	if req.IBCChannel != nil {
		sourceChannel := req.IBCChannel.ChannelId
		sourcePort := req.IBCChannel.PortId

		channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, sourceChannel))
		if !ok {
			err = ctx.EventManager().EmitTypedEvent(&types.EventPackedSendFailed{
				Error: fmt.Sprintf("channel not found for port ID (%s) channel ID (%s)", sourcePort, sourceChannel),
			})
			if err != nil {
				panic(err)
			}
		}

		packetData := types.NewOracleResponsePacketData(
			req.ClientID, requestID, req.RequestTime, ctx.BlockTime().Unix(), status, result,
		)

		//TODO: remove log
		ctx.Logger().Info("inside save request", "source channel", sourceChannel,
			"source port", sourcePort, "packet data", packetData.String(),
			"binary", packetData.ToBytes(), "json", types.ModuleCdc.MustMarshalJSON(packetData))

		if _, err := k.channelKeeper.SendPacket(
			ctx,
			channelCap,
			sourcePort,
			sourceChannel,
			clienttypes.NewHeight(0, 0),
			uint64(ctx.BlockTime().UnixNano()+int64(expiry)),
			packetData.ToBytes(),
		); err != nil {
			err = ctx.EventManager().EmitTypedEvents(&types.EventPackedSendFailed{
				Error: fmt.Sprintf("unable to send packet %s", err),
			})

			if err != nil {
				panic(err)
			}
		}
	}
}

func (k Keeper) ProcessRequestCalldata(
	ctx sdk.Context,
	requestEncoded []byte,
) (resultEncoded []byte, status types.ResolveStatus) {
	var request types.RequestPrice
	err := k.cdc.Unmarshal(requestEncoded, &request)
	if err != nil {
		return nil, types.RESOLVE_STATUS_FAILURE
	}

	switch request.Request {
	case types.PRICE_REQUEST_RATE:
		priceStamps, err := k.oracleKeeper.IterateExchangeRatesForDenoms(ctx, request.GetDenoms())
		if err != nil {
			return nil, types.RESOLVE_STATUS_FAILURE
		}

		result := types.OracleRequestResult{}
		for _, price := range priceStamps {
			result.PriceData = append(result.PriceData, types.OracleData{
				ExchangeRate: []sdk.DecCoin{*price.ExchangeRate},
				BlockNum:     []uint64{price.BlockNum},
			})
		}

		resultEncoded, err = result.Marshal()
		if err != nil {
			return nil, types.RESOLVE_STATUS_FAILURE
		}

	case types.PRICE_REQUEST_MEDIAN:
		numStamps := k.oracleKeeper.MaximumMedianStamps(ctx)
		mediansPriceStamps := k.oracleKeeper.IterateHistoricPricesForDenoms(
			ctx,
			oracleTypes.KeyPrefixMedian,
			request.GetDenoms(),
			uint(numStamps),
		)
		mediansPriceStamps.Sort()
		medianData := mediansPriceStamps.MapDenoms()

		result := types.OracleRequestResult{}
		for _, medians := range medianData {
			result.PriceData = append(result.PriceData, types.OracleData{
				ExchangeRate: medians.Rates,
				BlockNum:     medians.BlockNums,
			})
		}

		resultEncoded, err = result.Marshal()
		if err != nil {
			return nil, types.RESOLVE_STATUS_FAILURE
		}

	case types.PRICE_REQUEST_DEVIATION:
		numStamps := k.oracleKeeper.MaximumMedianStamps(ctx)
		deviationsPriceStamps := k.oracleKeeper.IterateHistoricPricesForDenoms(
			ctx,
			oracleTypes.KeyPrefixMedianDeviation,
			request.GetDenoms(),
			uint(numStamps),
		)
		deviationsPriceStamps.Sort()
		deviationData := deviationsPriceStamps.MapDenoms()

		result := types.OracleRequestResult{}
		for _, deviations := range deviationData {
			result.PriceData = append(result.PriceData, types.OracleData{
				ExchangeRate: deviations.Rates,
				BlockNum:     deviations.BlockNums,
			})
		}

		resultEncoded, err = result.Marshal()
		if err != nil {
			return nil, types.RESOLVE_STATUS_FAILURE
		}

	default:
		return nil, types.RESOLVE_STATUS_FAILURE
	}

	return resultEncoded, types.RESOLVE_STATUS_SUCCESS
}
