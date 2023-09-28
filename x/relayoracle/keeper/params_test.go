package keeper_test

import (
	"time"

	"github.com/ojo-network/ojo/x/relayoracle/types"
)

func (s *IntegrationTestSuite) TestGetParams() {
	app, ctx := s.app, s.ctx

	params := types.DefaultParams()
	params.PacketTimeout = uint64(10 * time.Minute)
	params.IbcRequestEnabled = false

	app.RelayOracle.SetParams(ctx, params)
	s.Require().Equal(app.RelayOracle.IbcRequestEnabled(ctx), params.IbcRequestEnabled)
	s.Require().Equal(app.RelayOracle.PacketTimeout(ctx), params.PacketTimeout)
}
