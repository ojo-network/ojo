package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/x/oracle/types"
)

// Test the reward giving mechanism
func (s *IntegrationTestSuite) TestRewardBallotWinners() {
	app, ctx := s.app, s.ctx

	// Add claim pools
	claims := []types.Claim{
		types.NewClaim(10, 10, 0, valAddr),
		types.NewClaim(20, 20, 0, valAddr2),
	}

	missCounters := []types.MissCounter{
		{ValidatorAddress: valAddr.String(), MissCounter: uint64(10)},
		{ValidatorAddress: valAddr2.String(), MissCounter: uint64(20)},
	}

	for _, mc := range missCounters {
		operator, _ := sdk.ValAddressFromBech32(mc.ValidatorAddress)
		app.OracleKeeper.SetMissCounter(ctx, operator, mc.MissCounter)
	}

	// Prepare reward pool
	givingAmt := sdk.NewCoins(sdk.NewInt64Coin(types.OjoDenom, 30000000))
	err := app.BankKeeper.MintCoins(ctx, "oracle", givingAmt)
	s.Require().NoError(err)

	var voteTargets []string
	params := app.OracleKeeper.GetParams(ctx)
	for _, v := range params.AcceptList {
		voteTargets = append(voteTargets, v.SymbolDenom)
	}

	votePeriodsPerWindow := sdk.NewDec((int64)(app.OracleKeeper.RewardDistributionWindow(ctx))).
		QuoInt64((int64)(app.OracleKeeper.VotePeriod(ctx))).
		TruncateInt64()
	app.OracleKeeper.RewardBallotWinners(ctx, (int64)(app.OracleKeeper.VotePeriod(ctx)), (int64)(app.OracleKeeper.RewardDistributionWindow(ctx)), voteTargets, claims)
	outstandingRewardsDecVal1 := app.DistrKeeper.GetValidatorOutstandingRewardsCoins(ctx, valAddr)
	outstandingRewardsVal1, _ := outstandingRewardsDecVal1.TruncateDecimal()
	outstandingRewardsDecVal2 := app.DistrKeeper.GetValidatorOutstandingRewardsCoins(ctx, valAddr2)
	outstandingRewardsVal2, _ := outstandingRewardsDecVal2.TruncateDecimal()
	s.Require().Equal(sdk.NewDecFromInt(givingAmt.AmountOf(types.OjoDenom)).QuoInt64(votePeriodsPerWindow).QuoInt64(3).TruncateInt(),
		outstandingRewardsVal1.AmountOf(types.OjoDenom))
	s.Require().Equal(sdk.NewDecFromInt(givingAmt.AmountOf(types.OjoDenom)).QuoInt64(votePeriodsPerWindow).QuoInt64(3).MulInt64(2).TruncateInt(),
		outstandingRewardsVal2.AmountOf(types.OjoDenom))
}

func (s *IntegrationTestSuite) TestRewardBallotWinnersZeroMissCounters() {
	app, ctx := s.app, s.ctx

	// Add claim pools
	claims := []types.Claim{
		types.NewClaim(10, 10, 0, valAddr),
		types.NewClaim(20, 20, 0, valAddr2),
	}

	// Prepare reward pool
	givingAmt := sdk.NewCoins(sdk.NewInt64Coin(types.OjoDenom, 30000000))
	err := app.BankKeeper.MintCoins(ctx, "oracle", givingAmt)
	s.Require().NoError(err)

	var voteTargets []string
	params := app.OracleKeeper.GetParams(ctx)
	for _, v := range params.AcceptList {
		voteTargets = append(voteTargets, v.SymbolDenom)
	}

	votePeriodsPerWindow := sdk.NewDec((int64)(app.OracleKeeper.RewardDistributionWindow(ctx))).
		QuoInt64((int64)(app.OracleKeeper.VotePeriod(ctx))).
		TruncateInt64()
	app.OracleKeeper.RewardBallotWinners(ctx, (int64)(app.OracleKeeper.VotePeriod(ctx)), (int64)(app.OracleKeeper.RewardDistributionWindow(ctx)), voteTargets, claims)
	outstandingRewardsDecVal1 := app.DistrKeeper.GetValidatorOutstandingRewardsCoins(ctx, valAddr)
	outstandingRewardsVal1, _ := outstandingRewardsDecVal1.TruncateDecimal()
	outstandingRewardsDecVal2 := app.DistrKeeper.GetValidatorOutstandingRewardsCoins(ctx, valAddr2)
	outstandingRewardsVal2, _ := outstandingRewardsDecVal2.TruncateDecimal()
	s.Require().Equal(sdk.NewDecFromInt(givingAmt.AmountOf(types.OjoDenom)).QuoInt64(votePeriodsPerWindow).QuoInt64(2).TruncateInt(),
		outstandingRewardsVal1.AmountOf(types.OjoDenom))
	s.Require().Equal(sdk.NewDecFromInt(givingAmt.AmountOf(types.OjoDenom)).QuoInt64(votePeriodsPerWindow).QuoInt64(2).TruncateInt(),
		outstandingRewardsVal2.AmountOf(types.OjoDenom))
}
