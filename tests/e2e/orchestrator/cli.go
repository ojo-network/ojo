package orchestrator

import "fmt"

func (o *Orchestrator) SubmitProposal(filePath string, govAddress string) error {
	err := o.ExecCommand([]string{
		"sed", "-i", fmt.Sprintf("s/$GOV_AUTHORITY_ADDRESS/%s/g", govAddress), filePath,
	})
	if err != nil {
		return err
	}
	return o.ExecCommand([]string{
		"ojod", "tx", "gov", "submit-proposal",
		filePath,
		"--from", "val",
		"--keyring-backend", "test",
		"--chain-id", o.chain.id,
		"--gas", "auto",
		"-b", "block",
	})
}

func (o *Orchestrator) SubmitLegacyParamChangeProposal(filePath string) error {
	return o.ExecCommand([]string{
		"ojod", "tx", "gov", "submit-legacy-proposal", "param-change",
		filePath,
		"--from", "val",
		"--keyring-backend", "test",
		"--chain-id", o.chain.id,
		"--gas", "auto",
		"-b", "block",
	})
}

func (o *Orchestrator) ProposalStatus(proposalID uint64) error {
	return o.ExecCommand([]string{
		"ojod", "query", "gov", "proposal", fmt.Sprint(proposalID),
		"--chain-id", o.chain.id,
		"--output", "json",
	})
}
