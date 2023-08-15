package types

import (
	gov "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

var (
	proposalTypeMsgGovUpdateParams       = MsgGovUpdateParams{}.String()
	proposalTypeMsgGovCancelUpdateParams = MsgGovCancelUpdateParams{}.String()
)

func init() {
	gov.RegisterProposalType(proposalTypeMsgGovUpdateParams)
	gov.RegisterProposalType(proposalTypeMsgGovCancelUpdateParams)
}

// Implements Proposal Interface
var _ gov.Content = &MsgGovUpdateParams{}

// GetTitle returns the title of a community pool spend proposal.
func (msg *MsgGovUpdateParams) GetTitle() string { return msg.Title }

// GetDescription returns the description of a community pool spend proposal.
func (msg *MsgGovUpdateParams) GetDescription() string { return msg.Description }

// GetDescription returns the routing key of a community pool spend proposal.
func (msg *MsgGovUpdateParams) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a community pool spend proposal.
func (msg *MsgGovUpdateParams) ProposalType() string { return proposalTypeMsgGovUpdateParams }

// Implements Proposal Interface
var _ gov.Content = &MsgGovCancelUpdateParams{}

// GetTitle returns the title of a community pool spend proposal.
func (msg *MsgGovCancelUpdateParams) GetTitle() string { return msg.Title }

// GetDescription returns the description of a community pool spend proposal.
func (msg *MsgGovCancelUpdateParams) GetDescription() string { return msg.Description }

// GetDescription returns the routing key of a community pool spend proposal.
func (msg *MsgGovCancelUpdateParams) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a community pool spend proposal.
func (msg *MsgGovCancelUpdateParams) ProposalType() string {
	return proposalTypeMsgGovCancelUpdateParams
}
