package e2e

import (
	"time"

	"github.com/ojo-network/ojo/tests/grpc"
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
	params, err := s.orchestrator.OjoClient.QueryClient.QueryParams()
	s.Require().NoError(err)

	s.Require().Equal(uint64(5), params.HistoricStampPeriod)
	s.Require().Equal(uint64(4), params.MaximumPriceStamps)
	s.Require().Equal(uint64(20), params.MedianStampPeriod)

	err = grpc.SubmitAndPassProposal(
		s.orchestrator.OjoClient,
		grpc.OracleParamChanges(10, 2, 20),
	)
	s.Require().NoError(err)

	params, err = s.orchestrator.OjoClient.QueryClient.QueryParams()
	s.Require().NoError(err)

	s.Require().Equal(uint64(10), params.HistoricStampPeriod)
	s.Require().Equal(uint64(2), params.MaximumPriceStamps)
	s.Require().Equal(uint64(20), params.MedianStampPeriod)

	s.Require().NoError(err)
}

// TestUpdateVotingPeriod updates the voting period with a gov prop
// and then verifies the new voting period is returned by the params query.
func (s *IntegrationTestSuite) TestUpdateVotingPeriod() {
	oldDuration := 5 * time.Second
	newDuration := 10 * time.Second

	params, err := s.orchestrator.OjoClient.QueryClient.QueryVotingParams()
	s.Require().NoError(err)

	s.Require().Equal(&oldDuration, params.VotingPeriod)

	// s.orchestrator.SubmitProposal("/root/proposals/voting_period.json")
	s.orchestrator.SubmitLegacyParamChangeProposal("/root/proposals/voting_period_legacy.json")

	params, err = s.orchestrator.OjoClient.QueryClient.QueryVotingParams()
	s.Require().NoError(err)

	s.Require().Equal(&newDuration, params.VotingPeriod)

	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TestLegacyProposal() {
	s.orchestrator.SubmitLegacyParamChangeProposal("/root/proposals/max_vals.json")
}
