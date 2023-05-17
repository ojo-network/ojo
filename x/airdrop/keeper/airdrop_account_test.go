package keeper_test

import (
	"github.com/ojo-network/ojo/x/airdrop/types"
)

func (s *IntegrationTestSuite) TestSetAndGetAirdropAccount() {
	app, ctx := s.app, s.ctx

	originAmount := uint64(500)
	airdropAccount := &types.AirdropAccount{
		OriginAddress: "test",
		OriginAmount:  originAmount,
	}
	err := app.AirdropKeeper.SetAirdropAccount(ctx, airdropAccount)
	s.Require().NoError(err)

	airdropAccount2, err := app.AirdropKeeper.GetAirdropAccount(ctx, "test")
	s.Require().NoError(err)

	s.Require().Equal(airdropAccount2.OriginAddress, airdropAccount.OriginAddress)
	s.Require().Equal(airdropAccount2.OriginAmount, airdropAccount.OriginAmount)
}

func (s *IntegrationTestSuite) TestGetAllAirdropAccounts() {
	app, ctx := s.app, s.ctx

	originAmount := uint64(500)
	airdropAccount := &types.AirdropAccount{
		OriginAddress: "test",
		OriginAmount:  originAmount,
	}
	err := app.AirdropKeeper.SetAirdropAccount(ctx, airdropAccount)
	s.Require().NoError(err)

	airdropAccount2 := &types.AirdropAccount{
		OriginAddress: "test2",
		OriginAmount:  originAmount,
	}
	err = app.AirdropKeeper.SetAirdropAccount(ctx, airdropAccount2)
	s.Require().NoError(err)

	accounts := app.AirdropKeeper.GetAllAirdropAccounts(ctx)
	s.Require().Equal(2, len(accounts))
}
