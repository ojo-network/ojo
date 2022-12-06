package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (s *IntegrationTestSuite) TestSlashAndResetMissCounters() {
	// initial setup
	addr, addr2 := valAddr, valAddr2
	amt := sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction)

	s.Require().Equal(amt, s.app.StakingKeeper.Validator(s.ctx, addr).GetBondedTokens())
	s.Require().Equal(amt, s.app.StakingKeeper.Validator(s.ctx, addr2).GetBondedTokens())

	votePeriodsPerWindow := sdk.NewDec(int64(s.app.OracleKeeper.SlashWindow(s.ctx))).QuoInt64(int64(s.app.OracleKeeper.VotePeriod(s.ctx))).TruncateInt64()
	numberOfAssets := int64(len(s.app.OracleKeeper.GetParams(s.ctx).AcceptList))
	possibleVotesPerSlashWindow := votePeriodsPerWindow * numberOfAssets
	slashFraction := s.app.OracleKeeper.SlashFraction(s.ctx)
	minValidVotes := s.app.OracleKeeper.MinValidPerWindow(s.ctx).MulInt64(votePeriodsPerWindow * numberOfAssets).TruncateInt64()
	// Case 1, no slash
	s.app.OracleKeeper.SetMissCounter(s.ctx, valAddr, uint64(possibleVotesPerSlashWindow-minValidVotes))
	s.app.OracleKeeper.SlashAndResetMissCounters(s.ctx)
	staking.EndBlocker(s.ctx, *s.app.StakingKeeper)

	validator, _ := s.app.StakingKeeper.GetValidator(s.ctx, valAddr)
	s.Require().Equal(amt, validator.GetBondedTokens())

	// Case 2, slash
	s.app.OracleKeeper.SetMissCounter(s.ctx, valAddr, uint64(possibleVotesPerSlashWindow-minValidVotes+1))
	s.app.OracleKeeper.SlashAndResetMissCounters(s.ctx)
	validator, _ = s.app.StakingKeeper.GetValidator(s.ctx, valAddr)
	s.Require().Equal(amt.Sub(slashFraction.MulInt(amt).TruncateInt()), validator.GetBondedTokens())
	s.Require().True(validator.Jailed)

	// Case 3, slash unbonded validator
	validator, _ = s.app.StakingKeeper.GetValidator(s.ctx, valAddr)
	validator.Status = stakingTypes.Unbonded
	validator.Jailed = false
	validator.Tokens = amt
	s.app.StakingKeeper.SetValidator(s.ctx, validator)

	s.app.OracleKeeper.SetMissCounter(s.ctx, valAddr, uint64(possibleVotesPerSlashWindow-minValidVotes+1))
	s.app.OracleKeeper.SlashAndResetMissCounters(s.ctx)
	validator, _ = s.app.StakingKeeper.GetValidator(s.ctx, valAddr)
	s.Require().Equal(amt, validator.Tokens)
	s.Require().False(validator.Jailed)

	// Case 4, slash jailed validator
	validator, _ = s.app.StakingKeeper.GetValidator(s.ctx, valAddr)
	validator.Status = stakingTypes.Bonded
	validator.Jailed = true
	validator.Tokens = amt
	s.app.StakingKeeper.SetValidator(s.ctx, validator)

	s.app.OracleKeeper.SetMissCounter(s.ctx, valAddr, uint64(possibleVotesPerSlashWindow-minValidVotes+1))
	s.app.OracleKeeper.SlashAndResetMissCounters(s.ctx)
	validator, _ = s.app.StakingKeeper.GetValidator(s.ctx, valAddr)
	s.Require().Equal(amt, validator.Tokens)
}
