package keeper_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	ojoapp "github.com/ojo-network/ojo/app"
	appparams "github.com/ojo-network/ojo/app/params"
	"github.com/ojo-network/ojo/tests/integration"
	"github.com/ojo-network/ojo/x/gmp/keeper"
	"github.com/ojo-network/ojo/x/gmp/types"
)

const (
	displayDenom string = appparams.DisplayDenom
	bondDenom    string = appparams.BondDenom
)

type IntegrationTestSuite struct {
	suite.Suite

	ctx       sdk.Context
	app       *ojoapp.App
	keys      []integration.TestValidatorKey
	msgServer types.MsgServer
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.app, s.ctx, s.keys = integration.SetupAppWithContext(s.T())
	s.msgServer = keeper.NewMsgServerImpl(s.app.GmpKeeper)
	s.SetOraclePrices()
}

func (s *IntegrationTestSuite) SetOraclePrices() {
	app, ctx := s.app, s.ctx
	app.OracleKeeper.SetExchangeRate(ctx, "ATOM", math.LegacyNewDecWithPrec(1, 1))
	app.OracleKeeper.SetExchangeRate(ctx, "BTC", math.LegacyNewDecWithPrec(1, 3))
}
