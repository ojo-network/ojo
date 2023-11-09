package gmp_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"gotest.tools/v3/assert"

	"github.com/stretchr/testify/suite"

	ojoapp "github.com/ojo-network/ojo/app"
	"github.com/ojo-network/ojo/tests/integration"
	"github.com/ojo-network/ojo/x/gmp"
	"github.com/ojo-network/ojo/x/gmp/types"
)

type IntegrationTestSuite struct {
	suite.Suite

	ctx  sdk.Context
	app  *ojoapp.App
	keys []integration.TestValidatorKey
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (s *IntegrationTestSuite) SetupTest() {
	s.app, s.ctx, s.keys = integration.SetupAppWithContext(s.T())
}

func (s *IntegrationTestSuite) TestGenesis_InitGenesis() {
	keeper, ctx := s.app.GmpKeeper, s.ctx

	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
	}

	s.Assertions.NotPanics(func() { gmp.InitGenesis(ctx, keeper, genesisState) })
}

func (s *IntegrationTestSuite) TestGenesis_ExportGenesis() {
	keeper, ctx := s.app.GmpKeeper, s.ctx

	params := types.DefaultParams()

	genesisState := types.GenesisState{
		Params: params,
	}

	gmp.InitGenesis(ctx, keeper, genesisState)

	result := gmp.ExportGenesis(ctx, keeper)

	assert.DeepEqual(s.T(), params, result.Params)
}
