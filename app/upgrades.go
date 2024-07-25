package app

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	v1migrations "github.com/cosmos/cosmos-sdk/x/gov/migrations/v1"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
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
}

// performs upgrade from v0.1.3 to v0.1.4
func (app *App) registerUpgrade0_1_4(_ upgradetypes.Plan) {
	const planName = "v0.1.4"
	app.UpgradeKeeper.SetUpgradeHandler(planName,
		func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			ctx.Logger().Info("Upgrade handler execution", "name", planName)
			upgrader := oraclekeeper.NewMigrator(&app.OracleKeeper)
			upgrader.MigrateValidatorSet(ctx)
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
		func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			// Migrate CometBFT consensus parameters from x/params module to a dedicated x/consensus module.
			baseapp.MigrateParams(ctx, baseAppLegacySS, &app.ConsensusParamsKeeper)

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
		func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			ctx.Logger().Info("Upgrade handler execution", "name", planName)
			return app.mm.RunMigrations(ctx, app.configurator, fromVM)
		},
	)
}

// performs upgrade from v0.2.1 to v0.2.2
func (app *App) registerUpgrade0_2_2(_ upgradetypes.Plan) {
	const planName = "v0.2.2"
	app.UpgradeKeeper.SetUpgradeHandler(planName,
		func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			ctx.Logger().Info("Upgrade handler execution", "name", planName)
			upgrader := oraclekeeper.NewMigrator(&app.OracleKeeper)
			upgrader.MigrateCurrencyPairProviders(ctx)
			upgrader.MigrateCurrencyDeviationThresholds(ctx)
			return app.mm.RunMigrations(ctx, app.configurator, fromVM)
		},
	)
}

func (app *App) registerUpgrade0_3_0(upgradeInfo upgradetypes.Plan) {
	const planName = "v0.3.0"
	app.UpgradeKeeper.SetUpgradeHandler(planName,
		func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			ctx.Logger().Info("Upgrade handler execution", "name", planName)
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
		func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			ctx.Logger().Info("Upgrade handler execution", "name", planName)
			return app.mm.RunMigrations(ctx, app.configurator, fromVM)
		},
	)
}

func (app *App) registerUpgrade0_3_1Rc1(_ upgradetypes.Plan) {
	const planName = "v0.3.1-rc1"
	app.UpgradeKeeper.SetUpgradeHandler(planName,
		func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			ctx.Logger().Info("Upgrade handler execution", "name", planName)
			return app.mm.RunMigrations(ctx, app.configurator, fromVM)
		},
	)
}

func (app *App) registerUpgrade0_3_1Rc2(_ upgradetypes.Plan) {
	const planName = "v0.3.1-rc2"
	app.UpgradeKeeper.SetUpgradeHandler(planName,
		func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			ctx.Logger().Info("Upgrade handler execution", "name", planName)
			return app.mm.RunMigrations(ctx, app.configurator, fromVM)
		},
	)
}

func (app *App) registerUpgrade0_3_1(_ upgradetypes.Plan) {
	const planName = "v0.3.1"

	app.UpgradeKeeper.SetUpgradeHandler(planName,
		func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			ctx.Logger().Info("Upgrade handler execution", "name", planName)
			return app.mm.RunMigrations(ctx, app.configurator, fromVM)
		},
	)
}

func (app *App) registerUpgrade0_3_2(_ upgradetypes.Plan) {
	const planName = "v0.3.2"

	app.UpgradeKeeper.SetUpgradeHandler(planName,
		func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			// migrate old proposals
			err := migrateProposals(ctx, app.keys[govtypes.StoreKey], app.appCodec)
			if err != nil {
				ctx.Logger().Error("failed to migrate governance proposals", "err", err)
			}

			ctx.Logger().Info("Upgrade handler execution", "name", planName)
			return app.mm.RunMigrations(ctx, app.configurator, fromVM)
		},
	)
}

// helper function to check if the store loader should be upgraded
func (app *App) storeUpgrade(planName string, ui upgradetypes.Plan, stores storetypes.StoreUpgrades) {
	if ui.Name == planName && !app.UpgradeKeeper.IsSkipHeight(ui.Height) {
		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(
			upgradetypes.UpgradeStoreLoader(ui.Height, &stores))
	}
}

// migrateProposals migrates all legacy MsgUpgateGovParam proposals into non legacy param update versions.
func migrateProposals(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)
	propStore := prefix.NewStore(store, v1migrations.ProposalsKeyPrefix)

	iter := propStore.Iterator(nil, nil)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var prop govv1.Proposal
		err := cdc.Unmarshal(iter.Value(), &prop)
		// if error unmarshaling prop, convert to non legacy prop
		if err != nil {
			newProp, err := convertProposal(prop, cdc)
			if err != nil {
				return err
			}
			bz, err := cdc.Marshal(&newProp)
			if err != nil {
				return err
			}
			// Set new value on store.
			propStore.Set(iter.Key(), bz)
		}
	}

	return nil
}

func convertProposal(prop govv1.Proposal, cdc codec.BinaryCodec) (govv1.Proposal, error) {
	msgs := prop.Messages

	for _, msg := range msgs {
		var oldUpdateParamMsg oracletypes.MsgLegacyGovUpdateParams
		err := cdc.Unmarshal(msg.GetValue(), &oldUpdateParamMsg)

		// if able to unmarshal into MsgLegacyGovUpdateParams, update to non legacy version
		if err != nil {
			newUpdateParamMsg := oracletypes.MsgGovUpdateParams{
				Authority:   oldUpdateParamMsg.Authority,
				Title:       oldUpdateParamMsg.Title,
				Description: oldUpdateParamMsg.Description,
				Plan: oracletypes.ParamUpdatePlan{
					Keys:    oldUpdateParamMsg.Keys,
					Height:  0, // placeholder value for height
					Changes: oldUpdateParamMsg.Changes,
				},
			}

			msg.Value, err = newUpdateParamMsg.Marshal()
			if err != nil {
				return govv1.Proposal{}, err
			}
		}
	}

	prop.Messages = msgs
	return prop, nil
}
