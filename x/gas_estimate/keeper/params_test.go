package keeper_test

import "github.com/ojo-network/ojo/x/gas_estimate/types"

func (s *IntegrationTestSuite) TestSetAndGetParams() {
	app, ctx := s.app, s.ctx

	params := types.Params{
		ContractRegistry: []*types.Contract{
			{
				Address: "0x0",
				Network: "Ethereum",
			},
		},
	}

	app.GasEstimateKeeper.SetParams(ctx, params)

	params2 := app.GasEstimateKeeper.GetParams(ctx)

	s.Require().Equal(params, params2)
}
