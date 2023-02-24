package query

import (
	"context"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

// GovQueryClient returns the govtypes.QueryClient
// initialized with the clients grpc connection
func (qc *QueryClient) GovQueryClient() govtypes.QueryClient {
	return govtypes.NewQueryClient(qc.grpcConn)
}

// QueryProposal sends a grpc query with the given proposalID
// and returns the govtypes.Proposal object
func (qc *QueryClient) QueryProposal(proposalID uint64) (*govtypes.Proposal, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	queryResponse, err := qc.GovQueryClient().Proposal(ctx, &govtypes.QueryProposalRequest{ProposalId: proposalID})
	if err != nil {
		return nil, err
	}
	return queryResponse.Proposal, nil
}
