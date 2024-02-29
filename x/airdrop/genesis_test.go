package airdrop_test

import (
	"gotest.tools/v3/assert"

	"cosmossdk.io/math"

	"github.com/ojo-network/ojo/x/airdrop"
	"github.com/ojo-network/ojo/x/airdrop/types"
)

func (s *IntegrationTestSuite) TestGenesis_InitGenesis() {
	keeper, ctx := s.app.AirdropKeeper, s.ctx

	genesisState := types.GenesisState{
		Params:          types.DefaultParams(),
		AirdropAccounts: []*types.AirdropAccount{},
	}

	s.Assertions.NotPanics(func() { airdrop.InitGenesis(ctx, keeper, genesisState) })
}

func (s *IntegrationTestSuite) TestGenesis_ExportGenesis() {
	keeper, ctx := s.app.AirdropKeeper, s.ctx

	params := types.DefaultParams()

	airdropAccounts := []*types.AirdropAccount{
		{
			VestingEndTime: 100,
			OriginAddress:  "ojo1ner6kc63xl903wrv2n8p9mtun79gegjld93lx0",
			OriginAmount:   math.NewInt(100).Uint64(),
		},
	}

	genesisState := types.GenesisState{
		Params:          params,
		AirdropAccounts: airdropAccounts,
	}

	airdrop.InitGenesis(ctx, keeper, genesisState)

	result := airdrop.ExportGenesis(ctx, keeper)

	assert.DeepEqual(s.T(), params, result.Params)
	assert.DeepEqual(s.T(), airdropAccounts, result.AirdropAccounts)
}
