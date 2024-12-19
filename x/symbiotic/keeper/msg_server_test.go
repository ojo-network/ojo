package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	appparams "github.com/ojo-network/ojo/app/params"
	"github.com/ojo-network/ojo/x/symbiotic/types"
)

var (
	pubKey    = secp256k1.GenPrivKey().PubKey()
	addr      = sdk.AccAddress(pubKey.Address())
	initCoins = sdk.NewCoins(sdk.NewCoin(appparams.BondDenom, math.NewInt(1000000000000000000)))
)

func (s *IntegrationTestSuite) TestMsgServer_SetParams() {
	SetParams(s, "0x0000000000000000000000000000000000000000", int64(10), uint64(10))

	params := types.DefaultParams()

	s.Require().Equal(params, s.app.GasEstimateKeeper.GetParams(s.ctx))
}

// SetParams sets the gasestimate module params
func SetParams(
	s *IntegrationTestSuite,
	middlewareAddress string,
	symbioticSyncPeriod int64,
	maximumCachedBlockHashes uint64,
) {
	authority := s.app.GovKeeper.GetGovernanceAccount(s.ctx).GetAddress().String()

	msg := types.NewMsgSetParams(
		authority,
		middlewareAddress,
		symbioticSyncPeriod,
		maximumCachedBlockHashes,
	)

	_, err := s.msgServer.SetParams(s.ctx, msg)
	s.Require().NoError(err)
}
