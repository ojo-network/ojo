package abci_test

import (
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ojo-network/ojo/tests/integration"
	"github.com/stretchr/testify/suite"

	ojoapp "github.com/ojo-network/ojo/app"
)

type IntegrationTestSuite struct {
	suite.Suite

	ctx  sdk.Context
	app  *ojoapp.App
	keys []integration.TestValidatorKey
}

func (s *IntegrationTestSuite) SetupTest() {
	s.app, s.ctx, s.keys = integration.SetupAppWithContext(s.T())
	s.app.OracleKeeper.SetVoteThreshold(s.ctx, math.LegacyMustNewDecFromStr("0.4"))
}

func TestAbciTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
