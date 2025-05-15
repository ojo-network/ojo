package keeper

import (
	"context"
	"sort"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ojo-network/ojo/util"
	"github.com/ojo-network/ojo/x/oracle/types"
)

var _ types.QueryServer = querier{}

// Querier implements a QueryServer for the x/oracle module.
type querier struct {
	Keeper
}

// NewQuerier returns an implementation of the oracle QueryServer interface
// for the provided Keeper.
func NewQuerier(keeper Keeper) types.QueryServer {
	return &querier{Keeper: keeper}
}

// Params queries params of x/oracle module.
func (q querier) Params(
	goCtx context.Context,
	req *types.QueryParams,
) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	params := q.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

// ExchangeRates queries exchange rates of all denoms, or, if specified, returns
// a single denom.
func (q querier) ExchangeRates(
	goCtx context.Context,
	req *types.QueryExchangeRates,
) (*types.QueryExchangeRatesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	var exchangeRates sdk.DecCoins

	if len(req.Denom) > 0 {
		exchangeRate, err := q.GetExchangeRate(ctx, req.Denom)
		if err != nil {
			return nil, err
		}

		exchangeRates = exchangeRates.Add(sdk.DecCoin{Denom: req.Denom, Amount: exchangeRate})
	} else {
		q.IterateExchangeRates(ctx, func(denom string, rate math.LegacyDec) (stop bool) {
			exchangeRates = exchangeRates.Add(sdk.DecCoin{Denom: denom, Amount: rate})
			return false
		})
	}

	return &types.QueryExchangeRatesResponse{ExchangeRates: exchangeRates}, nil
}

// ActiveExchangeRates queries all denoms for which exchange rates exist.
func (q querier) ActiveExchangeRates(
	goCtx context.Context,
	req *types.QueryActiveExchangeRates,
) (*types.QueryActiveExchangeRatesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	denoms := []string{}
	q.IterateExchangeRates(ctx, func(denom string, _ math.LegacyDec) (stop bool) {
		denoms = append(denoms, denom)
		return false
	})

	return &types.QueryActiveExchangeRatesResponse{ActiveRates: denoms}, nil
}

// FeederDelegation queries the account address to which the validator operator
// delegated oracle vote rights.
func (q querier) FeederDelegation(
	goCtx context.Context,
	req *types.QueryFeederDelegation,
) (*types.QueryFeederDelegationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	feederAddr, err := q.GetFeederDelegation(ctx, valAddr)
	if err != nil {
		return nil, err
	}

	return &types.QueryFeederDelegationResponse{
		FeederAddr: feederAddr.String(),
	}, nil
}

// MissCounter queries oracle miss counter of a validator.
func (q querier) MissCounter(
	goCtx context.Context,
	req *types.QueryMissCounter,
) (*types.QueryMissCounterResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryMissCounterResponse{
		MissCounter: q.GetMissCounter(ctx, valAddr.String()),
	}, nil
}

// SlashWindow queries the current slash window progress of the oracle.
func (q querier) SlashWindow(
	goCtx context.Context,
	req *types.QuerySlashWindow,
) (*types.QuerySlashWindowResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	params := q.GetParams(ctx)

	return &types.QuerySlashWindowResponse{
		WindowProgress: (util.SafeInt64ToUint64(ctx.BlockHeight()) % params.SlashWindow) /
			params.VotePeriod,
	}, nil
}

// AggregatePrevote queries an aggregate prevote of a validator.
func (q querier) AggregatePrevote(
	goCtx context.Context,
	req *types.QueryAggregatePrevote,
) (*types.QueryAggregatePrevoteResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	prevote, err := q.GetAggregateExchangeRatePrevote(ctx, valAddr)
	if err != nil {
		return nil, err
	}

	return &types.QueryAggregatePrevoteResponse{
		AggregatePrevote: prevote,
	}, nil
}

// AggregatePrevotes queries aggregate prevotes of all validators
func (q querier) AggregatePrevotes(
	goCtx context.Context,
	req *types.QueryAggregatePrevotes,
) (*types.QueryAggregatePrevotesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	var prevotes []types.AggregateExchangeRatePrevote
	q.IterateAggregateExchangeRatePrevotes(ctx, func(_ sdk.ValAddress, prevote types.AggregateExchangeRatePrevote) bool {
		prevotes = append(prevotes, prevote)
		return false
	})

	return &types.QueryAggregatePrevotesResponse{
		AggregatePrevotes: prevotes,
	}, nil
}

// AggregateVote queries an aggregate vote of a validator
func (q querier) AggregateVote(
	goCtx context.Context,
	req *types.QueryAggregateVote,
) (*types.QueryAggregateVoteResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	valAddr, err := sdk.ValAddressFromBech32(req.ValidatorAddr)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	vote, err := q.GetAggregateExchangeRateVote(ctx, valAddr.String())
	if err != nil {
		return nil, err
	}

	return &types.QueryAggregateVoteResponse{
		AggregateVote: vote,
	}, nil
}

// AggregateVotes queries aggregate votes of all validators
func (q querier) AggregateVotes(
	goCtx context.Context,
	req *types.QueryAggregateVotes,
) (*types.QueryAggregateVotesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	var votes []types.AggregateExchangeRateVote
	q.IterateAggregateExchangeRateVotes(ctx, func(_ string, vote types.AggregateExchangeRateVote) bool {
		votes = append(votes, vote)
		return false
	})

	return &types.QueryAggregateVotesResponse{
		AggregateVotes: votes,
	}, nil
}

// Medians queries medians of all denoms, or, if specified, returns
// a single median.
func (q querier) Medians(
	goCtx context.Context,
	req *types.QueryMedians,
) (*types.QueryMediansResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	medians := types.PriceStamps{}

	if len(req.Denom) > 0 {
		if req.NumStamps == 0 {
			return nil, status.Error(codes.InvalidArgument, "parameter NumStamps must be greater than 0")
		}

		if req.NumStamps > util.SafeUint64ToUint32(q.MaximumMedianStamps(ctx)) {
			req.NumStamps = util.SafeUint64ToUint32(q.MaximumMedianStamps(ctx))
		}

		medians = q.HistoricMedians(ctx, req.Denom, uint64(req.NumStamps))
	} else {
		medians = q.AllMedianPrices(ctx)
	}

	return &types.QueryMediansResponse{Medians: *medians.Sort()}, nil
}

// MedianDeviations queries median deviations of all denoms, or, if specified, returns
// a single median deviation.
func (q querier) MedianDeviations(
	goCtx context.Context,
	req *types.QueryMedianDeviations,
) (*types.QueryMedianDeviationsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	medianDeviations := types.PriceStamps{}

	if len(req.Denom) > 0 {
		price, err := q.HistoricMedianDeviation(ctx, req.Denom)
		if err != nil {
			return nil, err
		}
		medianDeviations = append(medianDeviations, *price)
	} else {
		medianDeviations = q.AllMedianDeviationPrices(ctx)
	}

	return &types.QueryMedianDeviationsResponse{MedianDeviations: *medianDeviations.Sort()}, nil
}

// ValidatorRewardSet queries the list of validators that can earn rewards in
// the current Slash Window.
func (q querier) ValidatorRewardSet(
	goCtx context.Context,
	req *types.QueryValidatorRewardSet,
) (*types.QueryValidatorRewardSetResponse, error) {
	return &types.QueryValidatorRewardSetResponse{}, nil
}

func (q querier) LatestPrices(goCtx context.Context, req *types.QueryLatestPricesRequest) (*types.QueryLatestPricesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	var prices []*types.LatestPrice

	for _, denom := range req.Denoms {
		tokenLatestPrice := q.GetDenomPrice(ctx, denom)
		prices = append(prices, &types.LatestPrice{
			Denom:       denom,
			LatestPrice: tokenLatestPrice.Dec(),
		})
	}

	return &types.QueryLatestPricesResponse{
		Prices: prices,
	}, nil
}

func (q querier) AllLatestPrices(goCtx context.Context, req *types.QueryAllLatestPricesRequest) (*types.QueryAllLatestPricesResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	var prices []*types.LatestPrice

	assetInfos := q.GetAllAssetInfo(ctx)

	for _, assetInfo := range assetInfos {
		tokenLatestPrice := q.GetDenomPrice(ctx, assetInfo.Denom)
		prices = append(prices, &types.LatestPrice{
			Denom:       assetInfo.Denom,
			LatestPrice: tokenLatestPrice.Dec(),
		})
	}

	return &types.QueryAllLatestPricesResponse{
		Prices: prices,
	}, nil
}

func (q querier) PriceHistory(goCtx context.Context, req *types.QueryPriceHistoryRequest) (*types.QueryPriceHistoryResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	prices := q.GetAllAssetPrices(ctx, req.Asset)

	sort.Slice(prices, func(i, j int) bool {
		return prices[i].Timestamp > prices[j].Timestamp
	})

	return &types.QueryPriceHistoryResponse{
		Prices: prices,
	}, nil
}

func (q querier) Pool(goCtx context.Context, req *types.QueryPoolRequest) (*types.QueryPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	pool, found := q.GetPool(ctx, req.PoolId)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryPoolResponse{
		Pool: pool,
	}, nil
}

func (q querier) PoolAll(goCtx context.Context, req *types.QueryPoolAllRequest) (*types.QueryPoolAllResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	pools := q.GetAllPool(ctx)

	return &types.QueryPoolAllResponse{Pool: pools}, nil
}

func (q querier) AccountedPoolAll(
	goCtx context.Context,
	req *types.QueryAccountedPoolAllRequest,
) (*types.QueryAccountedPoolAllResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	accountedPools := q.GetAllAccountedPool(ctx)

	return &types.QueryAccountedPoolAllResponse{AccountedPool: accountedPools}, nil
}

func (q querier) AccountedPool(
	goCtx context.Context,
	req *types.QueryAccountedPoolRequest,
) (*types.QueryAccountedPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := q.GetAccountedPool(ctx, req.PoolId)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryAccountedPoolResponse{AccountedPool: val}, nil
}

func (q querier) AssetInfoAll(
	goCtx context.Context,
	req *types.QueryAssetInfoAllRequest,
) (*types.QueryAssetInfoAllResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	assetInfos := q.GetAllAssetInfo(ctx)

	return &types.QueryAssetInfoAllResponse{AssetInfo: assetInfos}, nil
}

func (q querier) AssetInfo(
	goCtx context.Context,
	req *types.QueryAssetInfoRequest,
) (*types.QueryAssetInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := q.GetAssetInfo(ctx, req.Denom)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryAssetInfoResponse{AssetInfo: val}, nil
}
