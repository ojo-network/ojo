package keeper_test

import "math/rand"

func (s *IntegrationTestSuite) TestMissCounter() {
	app, ctx := s.app, s.ctx
	missCounter := uint64(rand.Intn(100))

	s.Require().Equal(app.OracleKeeper.GetMissCounter(ctx, valAddr), uint64(0))
	app.OracleKeeper.SetMissCounter(ctx, valAddr, missCounter)
	s.Require().Equal(app.OracleKeeper.GetMissCounter(ctx, valAddr), missCounter)

	app.OracleKeeper.DeleteMissCounter(ctx, valAddr)
	s.Require().Equal(app.OracleKeeper.GetMissCounter(ctx, valAddr), uint64(0))
}
