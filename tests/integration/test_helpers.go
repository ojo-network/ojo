package integration

import (
	"fmt"
	"testing"

	tmrand "github.com/cometbft/cometbft/libs/rand"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
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
		}
		validatorKeys = append(validatorKeys, valInfo)
	}

	return validatorKeys
}

func SetupAppWithContext(
	t *testing.T,
	numVals int,
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
	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{
		ChainID: fmt.Sprintf("test-chain-%s", tmrand.Str(4)),
		Height:  9,
	})

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	oracletypes.RegisterQueryServer(queryHelper, oraclekeeper.NewQuerier(app.OracleKeeper))
	airdroptypes.RegisterQueryServer(queryHelper, airdropkeeper.NewQuerier(app.AirdropKeeper))

	sh := stakingtestutil.NewHelper(t, ctx, app.StakingKeeper)
	sh.Denom = bondDenom

	// amt := sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction)

	// make validators 60%, 29%, 1% of total power

	initTokens := sdk.TokensFromConsensusPower(initialPower, sdk.DefaultPowerReduction)
	initCoins := sdk.NewCoins(sdk.NewCoin(appparams.BondDenom, initTokens))

	validatorKeys := CreateTestValidatorKeys(3)

	// mint and send coins to validators
	for _, val := range validatorKeys {
		require.NoError(t, app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initCoins))
		require.NoError(t, app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, val.AccAddress, initCoins))
		// sh.CreateValidator(val.ValAddress, val.PubKey, amt, true)
	}

	sh.CreateValidatorWithValPower(validatorKeys[0].ValAddress, validatorKeys[0].PubKey, 599, true)
	sh.CreateValidatorWithValPower(validatorKeys[1].ValAddress, validatorKeys[1].PubKey, 398, true)
	sh.CreateValidatorWithValPower(validatorKeys[2].ValAddress, validatorKeys[2].PubKey, 2, true)

	// mint and send coins to oracle module to fill up reward pool
	require.NoError(t, app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, initCoins))
	require.NoError(t, app.BankKeeper.SendCoinsFromModuleToModule(ctx, minttypes.ModuleName, oracletypes.ModuleName, initCoins))

	staking.EndBlocker(ctx, app.StakingKeeper)

	return app, ctx, validatorKeys
}
