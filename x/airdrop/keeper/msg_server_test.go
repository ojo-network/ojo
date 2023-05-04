package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/x/airdrop/types"
)

func (s *IntegrationTestSuite) TestMsgServer_CreateAirdropAccount() {
	tokensToReceive := sdk.MustNewDecFromStr("1000")
	msg := types.NewMsgCreateAirdropAccount(
		addr.String(),
		&tokensToReceive,
		20,
	)

	s.msgServer.CreateAirdropAccount(s.ctx, msg)

	airdropAccount, err := s.app.AirdropKeeper.GetAirdropAccount(s.ctx, addr.String())
	s.Require().NoError(err)
	s.Require().Equal(addr.String(), airdropAccount.OriginAddress)
	s.Require().Equal(tokensToReceive, *airdropAccount.OriginAmount)
	s.Require().Equal(int64(20), airdropAccount.VestingEndTime)
}
