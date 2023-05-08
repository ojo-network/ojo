package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/client/tx"
	"github.com/ojo-network/ojo/x/airdrop/types"
)

func (s *IntegrationTestSuite) TestMsgServer_SetParams() {
	params := types.DefaultParams()
	params.ExpiryBlock = 22000
	authority := s.app.GovKeeper.GetGovernanceAccount(s.ctx).GetAddress().String()

	msg := types.NewMsgSetParams(
		params.ExpiryBlock,
		params.DelegationRequirement,
		params.AirdropFactor,
		authority,
	)

	s.msgServer.SetParams(s.ctx, msg)

	s.Require().Equal(params, s.app.AirdropKeeper.GetParams(s.ctx))
}

func (s *IntegrationTestSuite) TestMsgServer_CreateAirdropAccount() {
	tokensToReceive := uint64(1000)

	originAddress, err := CreateAccount()
	s.Require().NoError(err)

	msg := types.NewMsgCreateAirdropAccount(
		originAddress.String(),
		tokensToReceive,
		20,
	)
	_, err = s.msgServer.CreateAirdropAccount(s.ctx, msg)
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
	tokensToReceive := uint64(1000)

	originAddress, err := CreateAccount()
	s.Require().NoError(err)

	claimAddress, err := CreateAccount()
	s.Require().NoError(err)

	msg := types.NewMsgCreateAirdropAccount(
		originAddress.String(),
		tokensToReceive,
		20,
	)
	_, err = s.msgServer.CreateAirdropAccount(s.ctx, msg)
	s.Require().NoError(err)

	msgClaimAirdrop := types.NewMsgClaimAirdrop(
		originAddress.String(),
		claimAddress.String(),
	)

	_, err = s.msgServer.ClaimAirdrop(s.ctx, msgClaimAirdrop)
	s.Require().NoError(err)

	airdropAccount, err := s.app.AirdropKeeper.GetAirdropAccount(s.ctx, originAddress.String())
	s.Require().NoError(err)
	s.Require().Equal(claimAddress.String(), airdropAccount.ClaimAddress)

	balance := s.app.BankKeeper.GetBalance(s.ctx, claimAddress, bondDenom)
	s.Require().Equal(airdropAccount.ClaimAmount, balance.Amount.Uint64())
}

func CreateAccount() (sdk.AccAddress, error) {
	mnemonic, err := tx.CreateMnemonic()
	if err != nil {
		return nil, err
	}
	account, _, err := tx.CreateAccountFromMnemonic("test", mnemonic)
	if err != nil {
		return nil, err
	}
	address, err := account.GetAddress()
	if err != nil {
		return nil, err
	}
	return address, nil
}
