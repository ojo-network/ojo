package e2e

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ojo-network/ojo/client/tx"
	"github.com/ojo-network/ojo/tests/grpc"
	airdroptypes "github.com/ojo-network/ojo/x/airdrop/types"

	appparams "github.com/ojo-network/ojo/app/params"
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
	title := "Update Airdrop Params"
	summary := "Update Airdrop Params expiry block, delegation requirement, and airdrop factor"

	err = grpc.SubmitAndPassProposal(ojoClient, []sdk.Msg{msg}, title, summary)
	s.Require().NoError(err)

	queriedParams, err := ojoClient.QueryClient.QueryAirdropParams()
	s.Require().NoError(err)

	s.Require().True(delegationRequirement.Equal(*queriedParams.DelegationRequirement))
}

func (s *IntegrationTestSuite) TestClaimAirdrop() {
	ojoClient := s.orchestrator.AirdropClient
	originAddress, err := ojoClient.TxClient.Address()
	s.Require().NoError(err)

	airdropAccount, err := ojoClient.QueryClient.QueryAirdropAccount(originAddress.String())
	s.Require().NoError(err)

	// Delegate tokens to qualify for claiming the airdrop
	originAccAddress, err := sdk.AccAddressFromBech32(originAddress.String())
	s.Require().NoError(err)
	val1Address, err := s.orchestrator.OjoClient.TxClient.Address()
	s.Require().NoError(err)

	val1ValAddress := sdk.ValAddress(val1Address)

	s.Require().NoError(err)
	_, err = ojoClient.TxClient.TxDelegate(
		originAccAddress,
		val1ValAddress,
		sdk.NewCoin(appparams.BondDenom, sdk.NewInt(int64(airdropAccount.OriginAmount))),
	)
	s.Require().NoError(err)

	// prevent account sequence mismatch
	time.Sleep(time.Second * 2)

	// Claim the airdrop
	claimAccount, err := tx.NewOjoAccount("claim_account")
	s.Require().NoError(err)
	claimAddress, err := claimAccount.KeyInfo.GetAddress()
	s.Require().NoError(err)

	_, err = ojoClient.TxClient.TxClaimAirdrop(originAddress.String(), claimAddress.String())
	s.Require().NoError(err)

	// Verify the new address has the claimed amount in it
	airdropAccount, err = ojoClient.QueryClient.QueryAirdropAccount(originAddress.String())
	s.Require().NoError(err)

	amount, err := ojoClient.QueryClient.QueryBalance(claimAddress.String(), appparams.BondDenom)
	s.Require().NoError(err)

	s.Require().Equal(airdropAccount.ClaimAmount, amount.Uint64())
}
