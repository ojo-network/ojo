package types

import (
	gov "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

var proposalTypeMsgGovUpdateParams = MsgGovUpdateParams{}.Type()

func init() {
	gov.RegisterProposalType(proposalTypeMsgGovUpdateParams)
}

// Implements Proposal Interface
var _ gov.Content = &MsgGovUpdateParams{}

// GetTitle returns the title of a community pool spend proposal.
func (msg *MsgGovUpdateParams) GetTitle() string { return msg.Plan.Title }

// GetDescription returns the description of a community pool spend proposal.
func (msg *MsgGovUpdateParams) GetDescription() string { return msg.Plan.Description }

// GetDescription returns the routing key of a community pool spend proposal.
func (msg *MsgGovUpdateParams) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a community pool spend proposal.
func (msg *MsgGovUpdateParams) ProposalType() string { return proposalTypeMsgGovUpdateParams }
