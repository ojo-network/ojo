package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/ojo-network/ojo/x/oracle/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the oracle MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (ms msgServer) AggregateExchangeRatePrevote(
	goCtx context.Context,
	msg *types.MsgAggregateExchangeRatePrevote,
) (*types.MsgAggregateExchangeRatePrevoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	valAddr, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return nil, err
	}
	feederAddr, err := sdk.AccAddressFromBech32(msg.Feeder)
	if err != nil {
		return nil, err
	}

	if err := ms.ValidateFeeder(ctx, feederAddr, valAddr); err != nil {
		return nil, err
	}

	// Ensure prevote wasn't already submitted
	if ms.HasAggregateExchangeRatePrevote(ctx, valAddr) {
		return nil, types.ErrExistingPrevote
	}

	// Convert hex string to votehash
	voteHash, err := types.AggregateVoteHashFromHexString(msg.Hash)
	if err != nil {
		return nil, types.ErrInvalidHash.Wrap(err.Error())
	}

	aggregatePrevote := types.NewAggregateExchangeRatePrevote(voteHash, valAddr, uint64(ctx.BlockHeight()))
	ms.SetAggregateExchangeRatePrevote(ctx, valAddr, aggregatePrevote)

	return &types.MsgAggregateExchangeRatePrevoteResponse{}, nil
}

func (ms msgServer) AggregateExchangeRateVote(
	goCtx context.Context,
	msg *types.MsgAggregateExchangeRateVote,
) (*types.MsgAggregateExchangeRateVoteResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	valAddr, err := sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return nil, err
	}
	feederAddr, err := sdk.AccAddressFromBech32(msg.Feeder)
	if err != nil {
		return nil, err
	}
	if err := ms.ValidateFeeder(ctx, feederAddr, valAddr); err != nil {
		return nil, err
	}

	params := ms.GetParams(ctx)
	aggregatePrevote, err := ms.GetAggregateExchangeRatePrevote(ctx, valAddr)
	if err != nil {
		return nil, types.ErrNoAggregatePrevote.Wrap(msg.Validator)
	}

	// Check a msg is submitted proper period
	if (uint64(ctx.BlockHeight())/params.VotePeriod)-(aggregatePrevote.SubmitBlock/params.VotePeriod) != 1 {
		return nil, types.ErrRevealPeriodMissMatch
	}

	exchangeRates, err := types.ParseExchangeRateDecCoins(msg.ExchangeRates)
	if err != nil {
		return nil, types.ErrInvalidExchangeRate.Wrap(err.Error())
	}

	// Verify a exchange rate with aggregate prevote hash
	hash := types.GetAggregateVoteHash(msg.Salt, msg.ExchangeRates, valAddr)
	if aggregatePrevote.Hash != hash.String() {
		return nil, types.ErrVerificationFailed.Wrapf("must be given %s not %s", aggregatePrevote.Hash, hash)
	}

	// Filter out rates which aren't included in the AcceptList
	filteredDecCoins := sdk.DecCoins{}
	for _, decCoin := range exchangeRates {
		if params.AcceptList.Contains(decCoin.Denom) {
			filteredDecCoins = append(filteredDecCoins, decCoin)
		}
	}

	// Move aggregate prevote to aggregate vote with given exchange rates
	ms.SetAggregateExchangeRateVote(ctx, valAddr, types.NewAggregateExchangeRateVote(filteredDecCoins, valAddr))
	ms.DeleteAggregateExchangeRatePrevote(ctx, valAddr)

	return &types.MsgAggregateExchangeRateVoteResponse{}, nil
}

func (ms msgServer) DelegateFeedConsent(
	goCtx context.Context,
	msg *types.MsgDelegateFeedConsent,
) (*types.MsgDelegateFeedConsentResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	operatorAddr, err := sdk.ValAddressFromBech32(msg.Operator)
	if err != nil {
		return nil, err
	}

	delegateAddr, err := sdk.AccAddressFromBech32(msg.Delegate)
	if err != nil {
		return nil, err
	}

	val := ms.StakingKeeper.Validator(ctx, operatorAddr)
	if val == nil {
		return nil, stakingtypes.ErrNoValidatorFound.Wrap(msg.Operator)
	}

	ms.SetFeederDelegation(ctx, operatorAddr, delegateAddr)
	err = ctx.EventManager().EmitTypedEvent(&types.EventDelegateFeedConsent{
		Operator: msg.Operator, Delegate: msg.Delegate,
	})

	return &types.MsgDelegateFeedConsentResponse{}, err
}

func (ms msgServer) GovUpdateParams(
	goCtx context.Context,
	msg *types.MsgGovUpdateParams,
) (*types.MsgGovUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	for _, key := range msg.Keys {
		switch key {
		case string(types.KeyVotePeriod):
			ms.SetVotePeriod(ctx, msg.Changes.VotePeriod)

		case string(types.KeyVoteThreshold):
			ms.SetVoteThreshold(ctx, msg.Changes.VoteThreshold)

		case string(types.KeyRewardBand):
			ms.SetRewardBand(ctx, msg.Changes.RewardBand)

		case string(types.KeyRewardDistributionWindow):
			if msg.Changes.RewardDistributionWindow < ms.Keeper.VotePeriod(ctx) {
				return nil, fmt.Errorf("oracle parameter RewardDistributionWindow must be greater than or equal with VotePeriod")
			}
			ms.SetRewardDistributionWindow(ctx, msg.Changes.RewardDistributionWindow)

		case string(types.KeyAcceptList):
			if !msg.Changes.AcceptList.ContainDenoms(ms.Keeper.MandatoryList(ctx)) {
				return nil, fmt.Errorf("denom in MandatoryList not present in AcceptList")
			}
			ms.SetAcceptList(ctx, msg.Changes.AcceptList)

		case string(types.KeyMandatoryList):
			if !ms.Keeper.AcceptList(ctx).ContainDenoms(msg.Changes.MandatoryList) {
				return nil, fmt.Errorf("denom in MandatoryList not present in AcceptList")
			}
			ms.SetMandatoryList(ctx, msg.Changes.MandatoryList)

		case string(types.KeySlashFraction):
			ms.SetSlashFraction(ctx, msg.Changes.SlashFraction)

		case string(types.KeySlashWindow):
			if msg.Changes.SlashWindow < ms.Keeper.VotePeriod(ctx) {
				return nil, fmt.Errorf("oracle parameter SlashWindow must be greater than or equal with VotePeriod")
			}
			ms.SetSlashWindow(ctx, msg.Changes.SlashWindow)

		case string(types.KeyMinValidPerWindow):
			ms.SetMinValidPerWindow(ctx, msg.Changes.MinValidPerWindow)

		case string(types.KeyHistoricStampPeriod):
			if msg.Changes.HistoricStampPeriod > ms.Keeper.MedianStampPeriod(ctx) {
				return nil, fmt.Errorf("oracle parameter HistoricStampPeriod must be less than or equal with MedianStampPeriod")
			}
			if msg.Changes.HistoricStampPeriod%ms.Keeper.VotePeriod(ctx) != 0 {
				return nil, fmt.Errorf("oracle parameters HistoricStampPeriod must be exact multiples of VotePeriod")
			}
			ms.SetHistoricStampPeriod(ctx, msg.Changes.HistoricStampPeriod)

		case string(types.KeyMedianStampPeriod):
			if msg.Changes.MedianStampPeriod < ms.Keeper.HistoricStampPeriod(ctx) {
				return nil, fmt.Errorf("oracle parameter MedianStampPeriod must be greater than or equal with HistoricStampPeriod")
			}
			if msg.Changes.MedianStampPeriod%ms.Keeper.VotePeriod(ctx) != 0 {
				return nil, fmt.Errorf("oracle parameters MedianStampPeriod must be exact multiples of VotePeriod")
			}
			ms.SetMedianStampPeriod(ctx, msg.Changes.MedianStampPeriod)

		case string(types.KeyMaximumPriceStamps):
			ms.SetMaximumPriceStamps(ctx, msg.Changes.MaximumPriceStamps)

		case string(types.KeyMaximumMedianStamps):
			ms.SetMaximumMedianStamps(ctx, msg.Changes.MaximumMedianStamps)

		default:
			return nil, fmt.Errorf("%s is not an existing orcale param key", key)
		}
	}

	return &types.MsgGovUpdateParamsResponse{}, nil
}
