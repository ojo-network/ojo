package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/x/airdrop/types"
)

func (s *IntegrationTestSuite) TestGetAirdropAccount() {
	app, ctx := s.app, s.ctx

	originAmount := sdk.MustNewDecFromStr("500")
	airdropAccount := types.AirdropAccount{
		OriginAddress: "test",
		OriginAmount:  &originAmount,
	}
	err := app.AirdropKeeper.SetAirdropAccount(ctx, airdropAccount)
	s.Require().NoError(err)

	airdropAccount2, err := app.AirdropKeeper.GetAirdropAccount(ctx, "test")
	s.Require().NoError(err)

	s.Require().Equal(airdropAccount2.OriginAddress, airdropAccount.OriginAddress)
	s.Require().Equal(airdropAccount2.OriginAmount, airdropAccount.OriginAmount)
}
