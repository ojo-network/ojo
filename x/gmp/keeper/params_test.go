package keeper_test

import "github.com/ojo-network/ojo/x/gmp/types"

func (s *IntegrationTestSuite) TestSetAndGetParams() {
	app, ctx := s.app, s.ctx

	params := types.Params{
		GmpChannel: "channel-101",
		GmpAddress: "gmpaddress",
		GmpTimeout: int64(101),
		FeeRecipient: "feerecipient",
	}

	app.GmpKeeper.SetParams(ctx, params)

	params2 := app.GmpKeeper.GetParams(ctx)

	s.Require().Equal(params2.GmpAddress, params.GmpAddress)
	s.Require().Equal(params2.GmpChannel, params.GmpChannel)
	s.Require().Equal(params2.GmpTimeout, params.GmpTimeout)
	s.Require().Equal(params2.FeeRecipient, params.FeeRecipient)
}
