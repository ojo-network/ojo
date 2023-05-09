package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/client/tx"
	"github.com/ojo-network/ojo/x/airdrop/types"
)

func (s *IntegrationTestSuite) TestMsgServer_SetParams() {
	expiryBlock := uint64(22000)
	delegationRequirement := sdk.MustNewDecFromStr("0.25")
	SetParams(s, expiryBlock, &delegationRequirement)

	params := types.DefaultParams()
	params.ExpiryBlock = expiryBlock
	params.DelegationRequirement = &delegationRequirement

	s.Require().Equal(params, s.app.AirdropKeeper.GetParams(s.ctx))
}

func (s *IntegrationTestSuite) TestMsgServer_CreateAirdropAccount() {
	tokensToReceive := uint64(1000)

	originAddress := CreateAccount(s)

	msg := types.NewMsgCreateAirdropAccount(
		originAddress.String(),
		tokensToReceive,
		20,
	)
	_, err := s.msgServer.CreateAirdropAccount(s.ctx, msg)
	s.Require().NoError(err)

	airdropAccount, err := s.app.AirdropKeeper.GetAirdropAccount(s.ctx, originAddress.String())
	s.Require().NoError(err)
	s.Require().Equal(originAddress.String(), airdropAccount.OriginAddress)
	s.Require().Equal(tokensToReceive, airdropAccount.OriginAmount)
	s.Require().Equal(msg.VestingEndTime, airdropAccount.VestingEndTime)

	balance := s.app.BankKeeper.GetBalance(s.ctx, originAddress, bondDenom)
	s.Require().Equal(tokensToReceive, balance.Amount.Uint64())
}

func (s *IntegrationTestSuite) TestMsgServer_ClaimAirdrop() {
	testCases := []struct {
		name                  string
		expiryBlock           uint64
		delegationRequirement sdk.Dec
		originAccount         sdk.AccAddress
		errMsg                string
	}{
		{
			name:                  "airdrop account doesn't exist",
			expiryBlock:           10000,
			delegationRequirement: sdk.MustNewDecFromStr("0"),
			originAccount:         CreateAccount(s),
			errMsg:                "no airdrop account found",
		},
		{
			name:                  "airdrop account already claimed",
			expiryBlock:           10000,
			delegationRequirement: sdk.MustNewDecFromStr("0"),
			originAccount:         CreateClaimedAccount(s),
			errMsg:                "already claimed by address",
		},
		{
			name:                  "past the expiry block",
			expiryBlock:           1,
			delegationRequirement: sdk.MustNewDecFromStr("0"),
			originAccount:         CreateAirdropAccount(s),
			errMsg:                "airdrop expired; chain is past the expire block",
		},
		{
			name:                  "delegation requirement not met",
			expiryBlock:           10000,
			delegationRequirement: sdk.MustNewDecFromStr("0.75"),
			originAccount:         CreateAirdropAccount(s),
			errMsg:                "delegation requirement not met",
		},
		{
			name:                  "claim successful",
			expiryBlock:           10000,
			delegationRequirement: sdk.MustNewDecFromStr("0"),
			originAccount:         CreateAirdropAccount(s),
			errMsg:                "",
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			SetParams(s, tc.expiryBlock, &tc.delegationRequirement)
			claimAddress := CreateAccount(s)

			msgClaimAirdrop := types.NewMsgClaimAirdrop(
				tc.originAccount.String(),
				claimAddress.String(),
			)

			_, err := s.msgServer.ClaimAirdrop(s.ctx, msgClaimAirdrop)
			if tc.errMsg != "" {
				s.Require().Error(err)
				s.Require().Contains(err.Error(), tc.errMsg)
				return
			}
			s.Require().NoError(err)

			airdropAccount, err := s.app.AirdropKeeper.GetAirdropAccount(s.ctx, tc.originAccount.String())
			s.Require().NoError(err)
			s.Require().Equal(claimAddress.String(), airdropAccount.ClaimAddress)

			balance := s.app.BankKeeper.GetBalance(s.ctx, claimAddress, bondDenom)
			s.Require().Equal(airdropAccount.ClaimAmount, balance.Amount.Uint64())
		})
	}
}

// Helper Functions \\

// CreateAccount creates a new account with a random mnemonic and returns the address
func CreateAccount(s *IntegrationTestSuite) sdk.AccAddress {
	mnemonic, err := tx.CreateMnemonic()
	s.Require().NoError(err)
	account, _, err := tx.CreateAccountFromMnemonic("test", mnemonic)
	s.Require().NoError(err)
	address, err := account.GetAddress()
	s.Require().NoError(err)
	return address
}

// CreateAirdropAccount uses the CreateAccount function to create an account and then
// creates an airdrop account using the new account as the origin address
func CreateAirdropAccount(s *IntegrationTestSuite) sdk.AccAddress {
	originAddress := CreateAccount(s)
	tokensToReceive := uint64(1000)
	msg := types.NewMsgCreateAirdropAccount(
		originAddress.String(),
		tokensToReceive,
		20,
	)
	_, err := s.msgServer.CreateAirdropAccount(s.ctx, msg)
	s.Require().NoError(err)
	return originAddress
}

func CreateClaimedAccount(s *IntegrationTestSuite) sdk.AccAddress {
	delegationRequirement := sdk.MustNewDecFromStr("0")
	SetParams(s, uint64(20000), &delegationRequirement)
	alreadyClaimedAcct := CreateAirdropAccount(s)
	claimAddress := CreateAccount(s)
	msgClaimAirdrop := types.NewMsgClaimAirdrop(
		alreadyClaimedAcct.String(),
		claimAddress.String(),
	)
	_, err := s.msgServer.ClaimAirdrop(s.ctx, msgClaimAirdrop)
	s.Require().NoError(err)
	return alreadyClaimedAcct
}

// SetParams sets the airdrop module params
func SetParams(
	s *IntegrationTestSuite,
	expiryBlock uint64,
	delegationRequirement *sdk.Dec,
) {
	params := types.DefaultParams()
	params.ExpiryBlock = expiryBlock
	params.DelegationRequirement = delegationRequirement
	authority := s.app.GovKeeper.GetGovernanceAccount(s.ctx).GetAddress().String()

	msg := types.NewMsgSetParams(
		params.ExpiryBlock,
		params.DelegationRequirement,
		params.AirdropFactor,
		authority,
	)

	_, err := s.msgServer.SetParams(s.ctx, msg)
	s.Require().NoError(err)
}
