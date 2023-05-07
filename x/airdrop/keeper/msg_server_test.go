package keeper_test

import (
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
	msg := types.NewMsgCreateAirdropAccount(
		addr.String(),
		tokensToReceive,
		20,
	)

	s.msgServer.CreateAirdropAccount(s.ctx, msg)

	airdropAccount, err := s.app.AirdropKeeper.GetAirdropAccount(s.ctx, addr.String())
	s.Require().NoError(err)
	s.Require().Equal(addr.String(), airdropAccount.OriginAddress)
	s.Require().Equal(tokensToReceive, airdropAccount.OriginAmount)
	s.Require().Equal(msg.VestingEndTime, airdropAccount.VestingEndTime)
}

func (s *IntegrationTestSuite) TestMsgServer_ClaimAirdrop() {

}
