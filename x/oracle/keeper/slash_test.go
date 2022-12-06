package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

func (s *IntegrationTestSuite) TestSlashAndResetMissCounters() {
	initialTokens := sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction)
	s.Require().Equal(initialTokens, s.app.StakingKeeper.Validator(s.ctx, valAddr).GetBondedTokens())

	var (
		slashFraction              = s.app.OracleKeeper.SlashFraction(s.ctx)
		possibleWinsPerSlashWindow = s.app.OracleKeeper.PossibleWinsPerSlashWindow(s.ctx)
		minValidPerWindow          = s.app.OracleKeeper.MinValidPerWindow(s.ctx)
		minValidVotes              = minValidPerWindow.MulInt64(possibleWinsPerSlashWindow).TruncateInt()
		maxMissesBeforeSlash       = sdk.NewInt(possibleWinsPerSlashWindow).Sub(minValidVotes).Uint64()
	)

	testCases := []struct {
		name         string
		missCounter  uint64
		status       stakingtypes.BondStatus
		jailedBefore bool
		jailedAfter  bool
		slashed      bool
	}{
		{
			name:         "bonded validator above minValidVotes",
			missCounter:  maxMissesBeforeSlash,
			status:       stakingtypes.Bonded,
			jailedBefore: false,
			jailedAfter:  false,
			slashed:      false,
		},
		{
			name:         "bonded validator below minValidVotes",
			missCounter:  maxMissesBeforeSlash + 1,
			status:       stakingtypes.Bonded,
			jailedBefore: false,
			jailedAfter:  true,
			slashed:      true,
		},
		{
			name:         "unBonded validator below minValidVotes",
			missCounter:  maxMissesBeforeSlash + 1,
			status:       stakingtypes.Unbonded,
			jailedBefore: false,
			jailedAfter:  false,
			slashed:      false,
		},
		{
			name:         "jailed validator below minValidVotes",
			missCounter:  maxMissesBeforeSlash + 1,
			status:       stakingtypes.Bonded,
			jailedBefore: true,
			jailedAfter:  true,
			slashed:      false,
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			validator, _ := s.app.StakingKeeper.GetValidator(s.ctx, valAddr)
			validator.Status = tc.status
			validator.Jailed = tc.jailedBefore
			validator.Tokens = initialTokens
			s.app.StakingKeeper.SetValidator(s.ctx, validator)

			s.app.OracleKeeper.SetMissCounter(s.ctx, valAddr, tc.missCounter)
			s.app.OracleKeeper.SlashAndResetMissCounters(s.ctx)

			expectedTokens := initialTokens
			if tc.slashed {
				expectedTokens = initialTokens.Sub(slashFraction.MulInt(initialTokens).TruncateInt())
			}

			validator, _ = s.app.StakingKeeper.GetValidator(s.ctx, valAddr)
			s.Require().Equal(expectedTokens, validator.Tokens)
			s.Require().Equal(tc.jailedAfter, validator.Jailed)
		})
	}

}
