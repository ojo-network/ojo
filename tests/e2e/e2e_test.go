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
	err := grpc.MedianCheck(s.orchestrator.OjoClient)
	s.Require().NoError(err)
}

// TestUpdateOracleParams updates the oracle params with a gov prop
// and then verifies the new params are returned by the params query.
func (s *IntegrationTestSuite) TestUpdateOracleParams() {
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
	delegationRequirement := sdk.MustNewDecFromStr("0.1")
	airdropFactor := sdk.MustNewDecFromStr("0.2")

	params := airdroptypes.Params{
		ExpiryBlock:           expiryBlock,
		DelegationRequirement: &delegationRequirement,
		AirdropFactor:         &airdropFactor,
	}

	ojoClient := s.orchestrator.OjoClient

	govAddress, err := ojoClient.QueryClient.QueryGovAccount()
	s.Require().NoError(err)

	resp, err := ojoClient.TxClient.TxSubmitAirdropProposal(&params, govAddress.Address)
	s.Require().NoError(err)

	proposalID, err := grpc.ParseProposalID(resp)
	s.Require().NoError(err)

	_, err = ojoClient.TxClient.TxVoteYes(proposalID)
	s.Require().NoError(err)

	err = grpc.SleepUntilProposalEndTime(ojoClient, proposalID)
	s.Require().NoError(err)

	err = grpc.VerifyProposalPassed(ojoClient, proposalID)
	s.Require().NoError(err)
}
