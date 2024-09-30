package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	ojoapp "github.com/ojo-network/ojo/app"
	appparams "github.com/ojo-network/ojo/app/params"
	"github.com/ojo-network/ojo/tests/integration"
	"github.com/ojo-network/ojo/x/gas_estimate/keeper"
	"github.com/ojo-network/ojo/x/gas_estimate/types"
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
	s.msgServer = keeper.NewMsgServerImpl(s.app.GasEstimateKeeper)
}
