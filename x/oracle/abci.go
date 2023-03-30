package oracle

import (
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/x/oracle/keeper"
	"github.com/ojo-network/ojo/x/oracle/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	params := k.GetParams(ctx)

	// Set all current active validators into the ValidatorRewardSet at
	// the beginning of a new Slash Window.
	if k.IsPeriodLastBlock(ctx, params.SlashWindow+1) {
		k.SetValidatorRewardSet(ctx)
	}

	if k.IsPeriodLastBlock(ctx, params.VotePeriod) {
		if err := CalcPrices(ctx, params, k); err != nil {
			return err
		}
	}

	// Slash oracle providers who missed voting over the threshold and reset
	// miss counters of all validators at the last block of slash window.
	if k.IsPeriodLastBlock(ctx, params.SlashWindow) {
		k.SlashAndResetMissCounters(ctx)
	}

	k.PruneAllPrices(ctx)
	go k.RecordEndBlockMetrics(ctx)

	return nil
}

func CalcPrices(ctx sdk.Context, params types.Params, k keeper.Keeper) error {
	// Build claim map over all validators in active set
	validatorClaimMap := make(map[string]types.Claim)
	powerReduction := k.StakingKeeper.PowerReduction(ctx)

	// Calculate total validator power
	var totalBondedPower int64
	for _, v := range k.StakingKeeper.GetBondedValidatorsByPower(ctx) {
		addr := v.GetOperator()
		power := v.GetConsensusPower(powerReduction)
		totalBondedPower += power
		validatorClaimMap[addr.String()] = types.NewClaim(power, 0, addr)
	}

	// voteTargets defines the symbol (ticker) denoms that we require votes on
	voteTargetDenoms := make([]string, 0)
	for _, v := range params.AcceptList {
		voteTargetDenoms = append(voteTargetDenoms, v.BaseDenom)
	}

	k.ClearExchangeRates(ctx)

	// NOTE: it filters out inactive or jailed validators
	ballotDenomSlice := k.OrganizeBallotByDenom(ctx, validatorClaimMap)
	threshold := k.VoteThreshold(ctx).MulInt64(types.MaxVoteThresholdMultiplier).TruncateInt64()

	// Iterate through ballots and update exchange rates; drop if not enough votes have been achieved.
	for _, ballotDenom := range ballotDenomSlice {
		// Increment Mandatory Win count if Denom in Mandatory list
		incrementWin := params.MandatoryList.Contains(ballotDenom.Denom)

		// If the asset is not in the mandatory or accept list, continue
		if !incrementWin && !params.AcceptList.Contains(ballotDenom.Denom) {
			ctx.Logger().Info("Unsupported denom, dropping ballot", "denom", ballotDenom)
			continue
		}

		// Calculate the portion of votes received as an integer, scaled up using the
		// same multiplier as the `threshold` computed above
		support := ballotDenom.Ballot.Power() * types.MaxVoteThresholdMultiplier / totalBondedPower
		if support < threshold {
			ctx.Logger().Info("Ballot voting power is under vote threshold, dropping ballot", "denom", ballotDenom)
			continue
		}

		// Get the current denom's reward band
		rewardBand, err := params.RewardBands.GetBandFromDenom(ballotDenom.Denom)
		if err != nil {
			return err
		}

		// Get median of exchange rates
		exchangeRate, err := Tally(ballotDenom.Ballot, rewardBand, validatorClaimMap, incrementWin)
		if err != nil {
			return err
		}

		// Set the exchange rate, emit ABCI event
		if err = k.SetExchangeRateWithEvent(ctx, ballotDenom.Denom, exchangeRate); err != nil {
			return err
		}

		if k.IsPeriodLastBlock(ctx, params.HistoricStampPeriod) {
			k.AddHistoricPrice(ctx, ballotDenom.Denom, exchangeRate)
		}

		// Calculate and stamp median/median deviation if median stamp period has passed
		if k.IsPeriodLastBlock(ctx, params.MedianStampPeriod) {
			if err = k.CalcAndSetHistoricMedian(ctx, ballotDenom.Denom); err != nil {
				return err
			}
		}
	}

	// Get validtors that can earn rewards in this Slash Window
	validatorRewardSet := types.ValidatorRewardSet{
		ValidatorMap: make(map[string]bool),
	}
	k.CurrentValidatorRewardSet(ctx, func(valRewardSet types.ValidatorRewardSet) bool {
		validatorRewardSet = valRewardSet
		return false
	})

	// update miss counting & slashing
	voteTargetsLen := len(params.MandatoryList)
	claimSlice := types.ClaimMapToSlice(validatorClaimMap, validatorRewardSet.ValidatorMap)
	for _, claim := range claimSlice {
		misses := uint64(voteTargetsLen - int(claim.MandatoryWinCount))
		if misses == 0 {
			continue
		}

		// Increase miss counter
		k.SetMissCounter(ctx, claim.Recipient, k.GetMissCounter(ctx, claim.Recipient)+misses)
	}

	// Distribute rewards to ballot winners
	k.RewardBallotWinners(
		ctx,
		int64(params.VotePeriod),
		int64(params.RewardDistributionWindow),
		voteTargetDenoms,
		claimSlice,
	)

	// Clear the ballot
	k.ClearBallots(ctx, params.VotePeriod)
	return nil
}

// Tally calculates and returns the median. It sets the set of voters to be
// rewarded, i.e. voted within a reasonable spread from the weighted median to
// the store. Note, the ballot is sorted by ExchangeRate.
func Tally(
	ballot types.ExchangeRateBallot,
	rewardBand sdk.Dec,
	validatorClaimMap map[string]types.Claim,
	incrementWin bool,
) (sdk.Dec, error) {
	median, err := ballot.Median()
	if err != nil {
		return sdk.ZeroDec(), err
	}
	standardDeviation, err := ballot.StandardDeviation(median)
	if err != nil {
		return sdk.ZeroDec(), err
	}

	// rewardSpread is the MAX((median * (rewardBand/2)), standardDeviation)
	rewardSpread := median.Mul(rewardBand.QuoInt64(2))
	rewardSpread = sdk.MaxDec(rewardSpread, standardDeviation)

	for _, tallyVote := range ballot {
		// Filter ballot winners. For voters, we filter out the tally vote iff:
		// (median - rewardSpread) <= ExchangeRate <= (median + rewardSpread)
		if (tallyVote.ExchangeRate.GTE(median.Sub(rewardSpread)) &&
			tallyVote.ExchangeRate.LTE(median.Add(rewardSpread))) ||
			!tallyVote.ExchangeRate.IsPositive() {

			key := tallyVote.Voter.String()
			claim := validatorClaimMap[key]

			if incrementWin {
				claim.MandatoryWinCount++
			}

			validatorClaimMap[key] = claim
		}
	}

	return median, nil
}
