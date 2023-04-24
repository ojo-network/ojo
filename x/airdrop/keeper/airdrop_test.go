package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ojo-network/ojo/x/airdrop/types"
)

func (s *IntegrationTestSuite) TestGetAirdropAccount(t *testing.T) {
	app, ctx := s.app, s.ctx

	originAmount := sdk.MustNewDecFromStr("500")
	airdropAccount := types.AirdropAccount{
		OriginAddress: "test",
		OriginAmount:  &originAmount,
	}
	err := app.AirdropKeeper.SetAirdropAccount(ctx, airdropAccount)
	s.Require().Error(err)

	airdropAccount2, err := app.AirdropKeeper.GetAirdropAccount(ctx, "test")
	s.Require().Error(err)

	s.Require().Equal(airdropAccount2.OriginAddress, airdropAccount.OriginAddress)
}
