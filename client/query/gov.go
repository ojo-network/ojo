package query

import (
	"context"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
)

// GovQueryClient returns the govtypes.QueryClient
// initialized with the clients grpc connection
func (c *Client) GovQueryClient() govtypes.QueryClient {
	return govtypes.NewQueryClient(c.grpcConn)
}

// QueryProposal sends a grpc query with the given proposalID
// and returns the govtypes.Proposal object
func (c *Client) QueryProposal(proposalID uint64) (*govtypes.Proposal, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	queryResponse, err := c.GovQueryClient().Proposal(ctx, &govtypes.QueryProposalRequest{ProposalId: proposalID})
	if err != nil {
		return nil, err
	}
	return queryResponse.Proposal, nil
}

func (c *Client) QueryVotingParams() (*govtypes.VotingParams, error) {
	ctx, cancel := context.WithTimeout(context.Background(), queryTimeout)
	defer cancel()

	queryResponse, err := c.GovQueryClient().Params(ctx, &govtypes.QueryParamsRequest{ParamsType: govtypes.ParamVoting})
	if err != nil {
		return nil, err
	}
	return queryResponse.VotingParams, nil
}
