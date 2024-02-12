package integration

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtestutil "github.com/cosmos/cosmos-sdk/x/staking/testutil"
	"github.com/stretchr/testify/require"

	ojoapp "github.com/ojo-network/ojo/app"
	appparams "github.com/ojo-network/ojo/app/params"
	airdropkeeper "github.com/ojo-network/ojo/x/airdrop/keeper"
	airdroptypes "github.com/ojo-network/ojo/x/airdrop/types"
	oraclekeeper "github.com/ojo-network/ojo/x/oracle/keeper"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
)

const (
	bondDenom    = appparams.BondDenom
	initialPower = int64(1000)
	isCheckTx    = false
)

var validatorPowers = []int64{599, 398, 2}

type TestValidatorKey struct {
	PubKey     cryptotypes.PubKey
	ValAddress sdk.ValAddress
	AccAddress sdk.AccAddress
	Power      int64
}

func CreateTestValidatorKeys(numValidators int) []TestValidatorKey {
	var validatorKeys []TestValidatorKey

	for i := 0; i < numValidators; i++ {
		pubKey := secp256k1.GenPrivKey().PubKey()
		valInfo := TestValidatorKey{
			PubKey:     pubKey,
			ValAddress: sdk.ValAddress(pubKey.Address()),
			AccAddress: sdk.AccAddress(pubKey.Address()),
			Power:      validatorPowers[i],
		}
		validatorKeys = append(validatorKeys, valInfo)
	}

	return validatorKeys
}

func SetupAppWithContext(
	t *testing.T,
) (
	*ojoapp.App,
	sdk.Context,
	[]TestValidatorKey,
) {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(appparams.AccountAddressPrefix, appparams.AccountPubKeyPrefix)
	config.SetBech32PrefixForValidator(appparams.ValidatorAddressPrefix, appparams.ValidatorPubKeyPrefix)
	config.SetBech32PrefixForConsensusNode(appparams.ConsNodeAddressPrefix, appparams.ConsNodePubKeyPrefix)

	app := ojoapp.Setup(t)
	ctx := app.BaseApp.NewContext(isCheckTx)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	oracletypes.RegisterQueryServer(queryHelper, oraclekeeper.NewQuerier(app.OracleKeeper))
	airdroptypes.RegisterQueryServer(queryHelper, airdropkeeper.NewQuerier(app.AirdropKeeper))

	sh := stakingtestutil.NewHelper(t, ctx, app.StakingKeeper)
	sh.Denom = bondDenom

	initTokens := sdk.TokensFromConsensusPower(initialPower, sdk.DefaultPowerReduction)
	initCoins := sdk.NewCoins(sdk.NewCoin(appparams.BondDenom, initTokens))

	validatorKeys := CreateTestValidatorKeys(3)

	// mint and send coins to validators
	for _, val := range validatorKeys {
		require.NoError(t, app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initCoins))
		require.NoError(t, app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, val.AccAddress, initCoins))
		sh.CreateValidatorWithValPower(val.ValAddress, val.PubKey, val.Power, true)
	}

	// mint and send coins to oracle module to fill up reward pool
	require.NoError(t, app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initCoins))
	require.NoError(t,
		app.BankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, oracletypes.ModuleName, initCoins),
	)

	_, err := app.StakingKeeper.EndBlocker(ctx)
	require.NoError(t, err)

	return app, ctx, validatorKeys
}
