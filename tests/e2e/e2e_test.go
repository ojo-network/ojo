package e2e

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ojo-network/ojo/tests/grpc"
	airdroptypes "github.com/ojo-network/ojo/x/airdrop/types"
)

// TestMedians queries for the oracle params, collects historical
// prices based on those params, checks that the stored medians and
// medians deviations are correct, updates the oracle params with
// a gov prop, then checks the medians and median deviations again.
func (s *IntegrationTestSuite) TestMedians() {
	s.T().SkipNow()
	err := grpc.MedianCheck(s.orchestrator.OjoClient)
	s.Require().NoError(err)
}

// TestUpdateOracleParams updates the oracle params with a gov prop
// and then verifies the new params are returned by the params query.
func (s *IntegrationTestSuite) TestUpdateOracleParams() {
	s.T().SkipNow()
	err := grpc.SubmitAndPassLegacyProposal(
		s.orchestrator.OjoClient,
		grpc.OracleParamChanges(10, 2, 20),
	)
	s.Require().NoError(err)

	params, err := s.orchestrator.OjoClient.QueryClient.QueryParams()
	s.Require().NoError(err)

	s.Require().Equal(uint64(10), params.HistoricStampPeriod)
	s.Require().Equal(uint64(2), params.MaximumPriceStamps)
	s.Require().Equal(uint64(20), params.MedianStampPeriod)
}

// TestUpdateAirdropParams updates the airdrop params with a gov prop
// and then verifies the new params are returned by the params query.
func (s *IntegrationTestSuite) TestUpdateAirdropParams() {
	expiryBlock := uint64(100)
	delegationRequirement := sdk.MustNewDecFromStr("8")
	airdropFactor := sdk.MustNewDecFromStr("7")

	params := airdroptypes.Params{
		ExpiryBlock:           expiryBlock,
		DelegationRequirement: &delegationRequirement,
		AirdropFactor:         &airdropFactor,
	}

	ojoClient := s.orchestrator.OjoClient

	govAddress, err := ojoClient.QueryClient.QueryGovAccount()
	s.Require().NoError(err)

	msg := airdroptypes.NewMsgSetParams(
		params.ExpiryBlock,
		params.DelegationRequirement,
		params.AirdropFactor,
		govAddress.Address,
	)

	err = grpc.SubmitAndPassProposal(ojoClient, []sdk.Msg{msg})
	s.Require().NoError(err)

	queriedParams, err := ojoClient.QueryClient.QueryAirdropParams()
	s.Require().NoError(err)

	s.Require().True(delegationRequirement.Equal(*queriedParams.DelegationRequirement))
}
