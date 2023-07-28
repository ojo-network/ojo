package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	oraclekeeper "github.com/ojo-network/ojo/x/oracle/keeper"
)

// RegisterUpgradeHandlersregisters upgrade handlers.
func (app App) RegisterUpgradeHandlers() {
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	app.registerUpgrade0_1_4(upgradeInfo)
	app.registerUpgrade1_0_0(upgradeInfo)
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

func (app *App) registerUpgrade1_0_0(_ upgradetypes.Plan) {
	planName := "v1.0"
	app.UpgradeKeeper.SetUpgradeHandler(planName,
		func(ctx sdk.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			ctx.Logger().Info("Upgrade handler execution", "name", planName)
			return app.mm.RunMigrations(ctx, app.configurator, fromVM)
		},
	)
}
