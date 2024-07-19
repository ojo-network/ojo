package app

import (
	"context"
	"fmt"

	storetypes "cosmossdk.io/store/types"
	circuittypes "cosmossdk.io/x/circuit/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	tenderminttypes "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	gmptypes "github.com/ojo-network/ojo/x/gmp/types"

	oraclekeeper "github.com/ojo-network/ojo/x/oracle/keeper"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
)

// RegisterUpgradeHandlersregisters upgrade handlers.
func (app App) RegisterUpgradeHandlers() {
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	app.registerUpgrade0_1_4(upgradeInfo)
	app.registerUpgrade0_2_0(upgradeInfo)
	app.registerUpgrade0_2_1(upgradeInfo)
	app.registerUpgrade0_2_2(upgradeInfo)
	app.registerUpgrade0_3_0(upgradeInfo)
	app.registerUpgrade0_3_0Rc8(upgradeInfo)
	app.registerUpgrade0_3_1Rc1(upgradeInfo)
	app.registerUpgrade0_3_1Rc2(upgradeInfo)
	app.registerUpgrade0_3_1(upgradeInfo)
	app.registerUpgrade0_3_2(upgradeInfo)
	app.registerUpgrade0_4_0(upgradeInfo)
}

// performs upgrade from v0.1.3 to v0.1.4
func (app *App) registerUpgrade0_1_4(_ upgradetypes.Plan) {
	const planName = "v0.1.4"
	app.UpgradeKeeper.SetUpgradeHandler(planName,
		func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			sdkCtx := sdk.UnwrapSDKContext(ctx)
			sdkCtx.Logger().Info("Upgrade handler execution", "name", planName)
			upgrader := oraclekeeper.NewMigrator(&app.OracleKeeper)
			err := upgrader.MigrateValidatorSet(sdkCtx)
			if err != nil {
				panic(err)
			}
			return app.mm.RunMigrations(ctx, app.configurator, fromVM)
		},
	)
}

//nolint: all
func (app *App) registerUpgrade0_2_0(upgradeInfo upgradetypes.Plan) {
	const planName = "v0.2.0"

	// Set param key table for params module migration
	for _, subspace := range app.ParamsKeeper.GetSubspaces() {
		subspace := subspace

		found := true
		var keyTable paramstypes.KeyTable
		switch subspace.Name() {
		case authtypes.ModuleName:
			keyTable = authtypes.ParamKeyTable()
		case banktypes.ModuleName:
			keyTable = banktypes.ParamKeyTable()
		case stakingtypes.ModuleName:
			keyTable = stakingtypes.ParamKeyTable()
		case minttypes.ModuleName:
			keyTable = minttypes.ParamKeyTable()
		case distrtypes.ModuleName:
			keyTable = distrtypes.ParamKeyTable()
		case slashingtypes.ModuleName:
			keyTable = slashingtypes.ParamKeyTable()
		case govtypes.ModuleName:
			keyTable = govv1.ParamKeyTable()
		case crisistypes.ModuleName:
			keyTable = crisistypes.ParamKeyTable()
		case oracletypes.ModuleName:
			keyTable = oracletypes.ParamKeyTable()
		default:
			// subspace not handled
			found = false
		}

		if found && !subspace.HasKeyTable() {
			subspace.WithKeyTable(keyTable)
		}
	}

	baseAppLegacySS := app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(
		paramstypes.ConsensusParamsKeyTable(),
	)

	app.UpgradeKeeper.SetUpgradeHandler(planName,
		func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			sdkCtx := sdk.UnwrapSDKContext(ctx)
			// Migrate CometBFT consensus parameters from x/params module to a dedicated x/consensus module.
			err := baseapp.MigrateParams(sdkCtx, baseAppLegacySS, &app.ConsensusParamsKeeper.ParamsStore)
			if err != nil {
				return nil, nil
			}

			return app.mm.RunMigrations(ctx, app.configurator, fromVM)
		},
	)

	app.storeUpgrade(planName, upgradeInfo, storetypes.StoreUpgrades{
		Added: []string{
			consensustypes.ModuleName,
			crisistypes.ModuleName,
		},
	})
}

func (app *App) registerUpgrade0_2_1(_ upgradetypes.Plan) {
	const planName = "v0.2.1"
	app.UpgradeKeeper.SetUpgradeHandler(planName,
		func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			sdkCtx := sdk.UnwrapSDKContext(ctx)
			sdkCtx.Logger().Info("Upgrade handler execution", "name", planName)
			return app.mm.RunMigrations(ctx, app.configurator, fromVM)
		},
	)
}

// performs upgrade from v0.2.1 to v0.2.2
func (app *App) registerUpgrade0_2_2(_ upgradetypes.Plan) {
	const planName = "v0.2.2"
	app.UpgradeKeeper.SetUpgradeHandler(planName,
		func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			sdkCtx := sdk.UnwrapSDKContext(ctx)
			sdkCtx.Logger().Info("Upgrade handler execution", "name", planName)
			upgrader := oraclekeeper.NewMigrator(&app.OracleKeeper)
			upgrader.MigrateCurrencyPairProviders(sdkCtx)
			upgrader.MigrateCurrencyDeviationThresholds(sdkCtx)
			return app.mm.RunMigrations(ctx, app.configurator, fromVM)
		},
	)
}

func (app *App) registerUpgrade0_3_0(upgradeInfo upgradetypes.Plan) {
	const planName = "v0.3.0"
	app.UpgradeKeeper.SetUpgradeHandler(planName,
		func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			sdkCtx := sdk.UnwrapSDKContext(ctx)
			sdkCtx.Logger().Info("Upgrade handler execution", "name", planName)
			return app.mm.RunMigrations(ctx, app.configurator, fromVM)
		},
	)

	app.storeUpgrade(planName, upgradeInfo, storetypes.StoreUpgrades{
		Added: []string{
			gmptypes.ModuleName,
		},
	})
}

func (app *App) registerUpgrade0_3_0Rc8(_ upgradetypes.Plan) {
	const planName = "v0.3.0-rc8"
	app.UpgradeKeeper.SetUpgradeHandler(planName,
		func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			sdkCtx := sdk.UnwrapSDKContext(ctx)
			sdkCtx.Logger().Info("Upgrade handler execution", "name", planName)
			return app.mm.RunMigrations(ctx, app.configurator, fromVM)
		},
	)
}

func (app *App) registerUpgrade0_3_1Rc1(_ upgradetypes.Plan) {
	const planName = "v0.3.1-rc1"
	app.UpgradeKeeper.SetUpgradeHandler(planName,
		func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			sdkCtx := sdk.UnwrapSDKContext(ctx)
			sdkCtx.Logger().Info("Upgrade handler execution", "name", planName)
			return app.mm.RunMigrations(ctx, app.configurator, fromVM)
		},
	)
}

func (app *App) registerUpgrade0_3_1Rc2(_ upgradetypes.Plan) {
	const planName = "v0.3.1-rc2"
	app.UpgradeKeeper.SetUpgradeHandler(planName,
		func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			sdkCtx := sdk.UnwrapSDKContext(ctx)
			sdkCtx.Logger().Info("Upgrade handler execution", "name", planName)
			return app.mm.RunMigrations(ctx, app.configurator, fromVM)
		},
	)
}

func (app *App) registerUpgrade0_3_1(_ upgradetypes.Plan) {
	const planName = "v0.3.1"

	app.UpgradeKeeper.SetUpgradeHandler(planName,
		func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			sdkCtx := sdk.UnwrapSDKContext(ctx)
			sdkCtx.Logger().Info("Upgrade handler execution", "name", planName)
			return app.mm.RunMigrations(ctx, app.configurator, fromVM)
		},
	)
}

func (app *App) registerUpgrade0_3_2(_ upgradetypes.Plan) {
	const planName = "v0.3.2"

	app.UpgradeKeeper.SetUpgradeHandler(planName,
		func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			sdkCtx := sdk.UnwrapSDKContext(ctx)

			// migrate old proposals
			govMigrator := govkeeper.NewMigrator(&app.GovKeeper, app.GetSubspace(govtypes.ModuleName))
			err := govMigrator.Migrate2to3(sdkCtx)
			if err != nil {
				panic("failed to migrate governance module")
			}

			sdkCtx.Logger().Info("Upgrade handler execution", "name", planName)
			return app.mm.RunMigrations(ctx, app.configurator, fromVM)
		},
	)
}

func (app *App) registerUpgrade0_4_0(upgradeInfo upgradetypes.Plan) {
	const planName = "v0.4.0"
	app.UpgradeKeeper.SetUpgradeHandler(planName,
		func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			sdkCtx := sdk.UnwrapSDKContext(ctx)
			sdkCtx.Logger().Info("Upgrade handler execution", "name", planName)

			// enable vote extensions after upgrade
			consensusParamsKeeper := app.ConsensusParamsKeeper
			currentParams, err := consensusParamsKeeper.Params(ctx, &consensustypes.QueryParamsRequest{})
			if err != nil || currentParams == nil || currentParams.Params == nil {
				panic(fmt.Sprintf("failed to retrieve existing consensus params in upgrade handler: %s", err))
			}
			currentParams.Params.Abci = &tenderminttypes.ABCIParams{
				VoteExtensionsEnableHeight: sdkCtx.BlockHeight() + int64(4), // enable vote extensions 4 blocks after upgrade
			}
			_, err = consensusParamsKeeper.UpdateParams(ctx, &consensustypes.MsgUpdateParams{
				Authority: consensusParamsKeeper.GetAuthority(),
				Block:     currentParams.Params.Block,
				Evidence:  currentParams.Params.Evidence,
				Validator: currentParams.Params.Validator,
				Abci:      currentParams.Params.Abci,
			})
			if err != nil {
				panic(fmt.Sprintf("failed to update consensus params : %s", err))
			}
			sdkCtx.Logger().Info(
				"Successfully set VoteExtensionsEnableHeight",
				"consensus_params",
				currentParams.Params.String(),
			)

			// update vote period to 1 block
			oracleKeeper := app.OracleKeeper
			oracleKeeper.SetVotePeriod(sdkCtx, 1)

			return app.mm.RunMigrations(ctx, app.configurator, fromVM)
		},
	)

	// REF: https://github.com/cosmos/cosmos-sdk/blob/a32186608aab0bd436049377ddb34f90006fcbf7/simapp/upgrades.go
	app.storeUpgrade(planName, upgradeInfo, storetypes.StoreUpgrades{
		Added: []string{
			circuittypes.ModuleName,
		},
	})
}

// helper function to check if the store loader should be upgraded
func (app *App) storeUpgrade(planName string, ui upgradetypes.Plan, stores storetypes.StoreUpgrades) {
	if ui.Name == planName && !app.UpgradeKeeper.IsSkipHeight(ui.Height) {
		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(
			upgradetypes.UpgradeStoreLoader(ui.Height, &stores))
	}
}
