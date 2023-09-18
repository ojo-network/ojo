package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	host "github.com/cosmos/ibc-go/v7/modules/core/24-host"

	sdkerrors "cosmossdk.io/errors"

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
	ctx.KVStore(k.storeKey).Set(types.ResultStoreKey(id), k.cdc.MustMarshal(&request))
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

	expiryTime := k.PacketExpiryBlockCount(ctx)
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

		if _, err := k.channelKeeper.SendPacket(
			ctx,
			channelCap,
			sourcePort,
			sourceChannel,
			clienttypes.NewHeight(0, 0),
			uint64(ctx.BlockTime().UnixNano())+expiryTime,
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

func (k Keeper) ProcessRequestCalldata(ctx sdk.Context, request []byte) (result []byte, status types.ResolveStatus) {
	var priceRequest types.RequestPrice
	err := k.cdc.Unmarshal(request, &priceRequest)
	if err != nil {
		return nil, types.RESOLVE_STATUS_FAILURE
	}

	switch priceRequest.Request {
	case types.PRICE_REQUEST_RATE:
		price, err := k.oracleKeeper.GetExchangeRate(ctx, priceRequest.GetDenom())
		if err != nil {
			return nil, types.RESOLVE_STATUS_FAILURE
		}

		result, err = price.Marshal()
		if err != nil {
			return nil, types.RESOLVE_STATUS_FAILURE
		}
	case types.PRICE_REQUEST_MEDIAN:
		medians := k.oracleKeeper.HistoricMedians(ctx, priceRequest.GetDenom(), 0)
		medians.Sort()
		var priceStamps types.PriceStamp
		for _, median := range medians {
			priceStamps.ExchangeRate = append(priceStamps.ExchangeRate, *median.ExchangeRate)
			priceStamps.BlockNum = append(priceStamps.BlockNum, median.BlockNum)
		}

		result, err = priceStamps.Marshal()
		if err != nil {
			return nil, types.RESOLVE_STATUS_FAILURE
		}

	case types.PRICE_REQUEST_DEVIATION:
		var priceStamp types.PriceStamp
		deviation, err := k.oracleKeeper.HistoricMedianDeviation(ctx, priceRequest.GetDenom())
		if err != nil {
			return nil, types.RESOLVE_STATUS_FAILURE
		}

		priceStamp.ExchangeRate = []sdk.DecCoin{*deviation.ExchangeRate}
		priceStamp.BlockNum = []uint64{deviation.BlockNum}

		result, err = priceStamp.Marshal()
		if err != nil {
			return nil, types.RESOLVE_STATUS_FAILURE
		}

	default:
		return nil, types.RESOLVE_STATUS_FAILURE
	}

	return result, types.RESOLVE_STATUS_SUCCESS
}
