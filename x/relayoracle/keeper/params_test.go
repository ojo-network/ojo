package keeper_test

import (
	"github.com/ojo-network/ojo/x/relayoracle/types"
)

func (s *IntegrationTestSuite) TestGetParams() {
	app, ctx := s.app, s.ctx

	params := types.DefaultParams()
	params.PacketExpiryBlockCount = 20
	params.IbcRequestEnabled = false

	app.RelayOracle.SetParams(ctx, params)
	s.Require().Equal(app.RelayOracle.IbcRequestEnabled(ctx), params.IbcRequestEnabled)
	s.Require().Equal(app.RelayOracle.PacketExpiryBlockCount(ctx), params.PacketExpiryBlockCount)
}
