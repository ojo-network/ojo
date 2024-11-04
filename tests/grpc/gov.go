package grpc

import (
	"fmt"
	"strconv"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	proposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/ojo-network/ojo/client"
	"github.com/rs/zerolog/log"
)

const extraWaitTime = 5 * time.Second // at least one full block

// VerifyProposalPassed returns a non-nil error if the proposal did not pass
func VerifyProposalPassed(ojoClient *client.OjoClient, proposalID uint64) error {
	prop, err := ojoClient.QueryClient.QueryProposal(proposalID)
	if err != nil {
		return err
	}
	status := prop.Status.String()
	if status != "PROPOSAL_STATUS_PASSED" {
		return fmt.Errorf("proposal %d failed to pass with status: %s", proposalID, status)
	}
	return nil
}

// SleepUntilProposalEndTime sleeps until the end of the voting period + 1 block
func SleepUntilProposalEndTime(ojoClient *client.OjoClient, proposalID uint64) error {
	prop, err := ojoClient.QueryClient.QueryProposal(proposalID)
	if err != nil {
		return err
	}

	now := time.Now()
	sleepDuration := prop.VotingEndTime.Sub(now) + extraWaitTime
	log.Info().Msgf("sleeping %s until end of voting period + 1 block", sleepDuration)
	time.Sleep(sleepDuration)
	return nil
}

// ParseProposalID parses the proposalID from a tx response
func ParseProposalID(response *sdk.TxResponse) (uint64, error) {
	for _, event := range response.Events {
		if event.Type == "submit_proposal" {
			for _, attribute := range event.Attributes {
				if attribute.Key == "proposal_id" {
					return strconv.ParseUint(attribute.Value, 10, 64)
				}
			}
		}
	}
	return 0, fmt.Errorf("unable to find proposalID in tx response")
}

func SubmitAndPassProposal(ojoClient *client.OjoClient, msgs []sdk.Msg, title, summary string) error {
	deposit := sdk.NewCoins(sdk.NewCoin("uojo", math.NewInt(10000000)))
	resp, err := ojoClient.TxClient.TxSubmitProposal(msgs, deposit, title, summary)
	if err != nil {
		return err
	}

	// retry
	resp, err = retryTx(ojoClient, resp)
	if err != nil {
		return err
	}

	proposalID, err := ParseProposalID(resp)
	if err != nil {
		return err
	}

	_, err = ojoClient.TxClient.TxVoteYes(proposalID)
	if err != nil {
		return err
	}

	err = SleepUntilProposalEndTime(ojoClient, proposalID)
	if err != nil {
		return err
	}

	return VerifyProposalPassed(ojoClient, proposalID)
}

// SubmitAndPassProposal submits a proposal and votes yes on it
func SubmitAndPassLegacyProposal(ojoClient *client.OjoClient, changes []proposal.ParamChange) error {
	resp, err := ojoClient.TxClient.TxSubmitLegacyProposal(changes)
	if err != nil {
		return err
	}

	// retry
	resp, err = retryTx(ojoClient, resp)
	if err != nil {
		return err
	}

	proposalID, err := ParseProposalID(resp)
	if err != nil {
		return err
	}

	_, err = ojoClient.TxClient.TxVoteYes(proposalID)
	if err != nil {
		return err
	}

	err = SleepUntilProposalEndTime(ojoClient, proposalID)
	if err != nil {
		return err
	}

	return VerifyProposalPassed(ojoClient, proposalID)
}

func OracleParamChanges(
	historicStampPeriod uint64,
	maximumPriceStamps uint64,
	medianStampPeriod uint64,
) []proposal.ParamChange {
	return []proposal.ParamChange{
		{
			Subspace: "oracle",
			Key:      "HistoricStampPeriod",
			Value:    fmt.Sprintf("\"%d\"", historicStampPeriod),
		},
		{
			Subspace: "oracle",
			Key:      "MaximumPriceStamps",
			Value:    fmt.Sprintf("\"%d\"", maximumPriceStamps),
		},
		{
			Subspace: "oracle",
			Key:      "MedianStampPeriod",
			Value:    fmt.Sprintf("\"%d\"", medianStampPeriod),
		},
	}
}

func retryTx(ojoClient *client.OjoClient, resp *sdk.TxResponse) (*sdk.TxResponse, error) {
	for i := 0; i < 5; i++ {
		newResp, err := ojoClient.QueryTxHash(resp.TxHash)
		if err != nil && i == 4 {
			return nil, err
		}
		if err == nil {
			resp = newResp
			break
		}

		time.Sleep(time.Second * (1 + time.Duration(i)))
	}

	return resp, nil
}
