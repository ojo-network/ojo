package keeper_test

import (
	"github.com/ojo-network/ojo/x/airdrop/types"
)

func (s *IntegrationTestSuite) TestSetAndGetAirdropAccount() {
	app, ctx := s.app, s.ctx

	originAmount := uint64(500)
	airdropAccount := types.NewAirdropAccount("test", originAmount, 0)
	err := app.AirdropKeeper.SetAirdropAccount(ctx, airdropAccount)
	s.Require().NoError(err)

	airdropAccount2, err := app.AirdropKeeper.GetAirdropAccount(ctx, "test", airdropAccount.State)
	s.Require().NoError(err)

	s.Require().Equal(airdropAccount2.OriginAddress, airdropAccount.OriginAddress)
	s.Require().Equal(airdropAccount2.OriginAmount, airdropAccount.OriginAmount)
}

func (s *IntegrationTestSuite) TestGetAllAirdropAccounts() {
	app, ctx := s.app, s.ctx

	accounts := app.AirdropKeeper.GetAllAirdropAccounts(ctx)
	prevAcctsLen := len(accounts)

	originAmount := uint64(500)
	airdropAccount := types.NewAirdropAccount("test", originAmount, 0)
	err := app.AirdropKeeper.SetAirdropAccount(ctx, airdropAccount)
	s.Require().NoError(err)

	airdropAccount2 := types.NewAirdropAccount("test2", originAmount, 0)
	err = app.AirdropKeeper.SetAirdropAccount(ctx, airdropAccount2)
	s.Require().NoError(err)

	accounts = app.AirdropKeeper.GetAllAirdropAccounts(ctx)
	s.Require().Equal(prevAcctsLen+2, len(accounts))
}

func (s *IntegrationTestSuite) TestPaginatedAirdropAccounts() {
	app, ctx := s.app, s.ctx

	accounts := app.AirdropKeeper.GetAllAirdropAccounts(ctx)
	prevAcctsLen := len(accounts)

	originAmount := uint64(500)
	airdropAccount := types.NewAirdropAccount("testpaginate", originAmount, 0)
	err := app.AirdropKeeper.SetAirdropAccount(ctx, airdropAccount)
	s.Require().NoError(err)

	airdropAccount2 := types.NewAirdropAccount("testpaginate2", originAmount, 0)
	err = app.AirdropKeeper.SetAirdropAccount(ctx, airdropAccount2)
	s.Require().NoError(err)

	accounts = app.AirdropKeeper.GetAllAirdropAccounts(ctx)
	s.Require().Equal(prevAcctsLen+2, len(accounts))

	accounts = app.AirdropKeeper.PaginatedAirdropAccounts(ctx, types.StateCreated, 1)
	s.Require().Equal(1, len(accounts))

	accounts2 := app.AirdropKeeper.PaginatedAirdropAccounts(ctx, types.StateCreated, 1)
	s.Require().Equal(1, len(accounts2))

	s.Require().Equal(accounts[0].OriginAddress, accounts2[0].OriginAddress)
}

func (s *IntegrationTestSuite) TestCreateAirdropAccount() {
	app, ctx := s.app, s.ctx

	tokensToReceive := uint64(1000)
	originAddress := CreateAccount(s)
	vestingEndTime := ctx.BlockTime().Unix() + 20
	airdropAccount := types.NewAirdropAccount(originAddress.String(), tokensToReceive, vestingEndTime)

	app.AirdropKeeper.CreateAirdropAccount(ctx, airdropAccount)

	airdropAccount, err := s.app.AirdropKeeper.GetAirdropAccount(s.ctx, originAddress.String(), airdropAccount.State)
	s.Require().NoError(err)
	s.Require().Equal(originAddress.String(), airdropAccount.OriginAddress)
	s.Require().Equal(tokensToReceive, airdropAccount.OriginAmount)
	s.Require().Equal(vestingEndTime, airdropAccount.VestingEndTime)

	balance := s.app.BankKeeper.GetBalance(s.ctx, originAddress, bondDenom)
	s.Require().Equal(tokensToReceive, balance.Amount.Uint64())
}
