package keeper_test

import (
	types1 "github.com/cosmos/cosmos-sdk/codec/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"

	"github.com/ojo-network/ojo/x/oracle/keeper"
	"github.com/ojo-network/ojo/x/oracle/types"
)

func (s *IntegrationTestSuite) TestMigrateProposal() {
	ctx := s.ctx
	cdc := s.app.AppCodec()
	storeKey := s.app.GetKey(govtypes.StoreKey)

	// create legacy prop and set it in store
	legacyMsg := types.MsgLegacyGovUpdateParams{
		Authority:   "ojo10d07y265gmmuvt4z0w9aw880jnsr700jcz4krc",
		Title:       "title",
		Description: "desc",
		Keys: []string{
			"VotePeriod",
		},
		Changes: types.Params{
			VotePeriod: 5,
		},
	}
	bz, err := cdc.Marshal(&legacyMsg)
	s.Require().NoError(err)
	prop := govv1.Proposal{
		Id: 1,
		Messages: []*types1.Any{
			{
				TypeUrl:          "/ojo.oracle.v1.MsgGovUpdateParams",
				Value:            bz,
				XXX_unrecognized: []byte{},
			},
		},
		Status: govv1.ProposalStatus_PROPOSAL_STATUS_PASSED,
	}
	s.app.GovKeeper.SetProposal(ctx, prop)

	// try to retreive proposal and fail
	_, err = s.app.GovKeeper.Proposals.Get(ctx, prop.Id)
	s.Require().Error(err)

	// succesfully retreive proposal after migration
	err = keeper.MigrateProposals(ctx, storeKey, cdc)
	_, err = s.app.GovKeeper.Proposals.Get(ctx, prop.Id)
	s.Require().NoError(err)
}
