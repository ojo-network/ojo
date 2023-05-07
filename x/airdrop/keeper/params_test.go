package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/x/airdrop/types"
)

func (s *IntegrationTestSuite) TestSetAndGetParams() {
	app, ctx := s.app, s.ctx

	delegationRequirement := sdk.MustNewDecFromStr("0.5")
	airdropFactor := sdk.MustNewDecFromStr("3")

	params := types.Params{
		ExpiryBlock:           uint64(10000),
		DelegationRequirement: &delegationRequirement,
		AirdropFactor:         &airdropFactor,
	}

	err := app.AirdropKeeper.SetParams(ctx, params)
	s.Require().NoError(err)

	params2, err := app.AirdropKeeper.GetParams(ctx)
	s.Require().NoError(err)

	s.Require().Equal(params2.ExpiryBlock, params.ExpiryBlock)
	s.Require().Equal(params2.DelegationRequirement, params.DelegationRequirement)
	s.Require().Equal(params2.AirdropFactor, params.AirdropFactor)
}
