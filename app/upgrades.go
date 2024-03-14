package app

import (
	"context"

	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	batchtypes "github.com/Team-Kujira/core/x/batch/types"
	cwicatypes "github.com/Team-Kujira/core/x/cw-ica/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	ibcwasmtypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/types"
)

const UpgradeName = "v1.0.0"

func (app App) RegisterUpgradeHandlers() {
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(err)
	}

	if upgradeInfo.Name == UpgradeName && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{
				ibcwasmtypes.ModuleName,
				batchtypes.ModuleName,
				cwicatypes.ModuleName,
			},
		}

		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}

	app.UpgradeKeeper.SetUpgradeHandler(
		UpgradeName,
		func(ctx context.Context,
			_ upgradetypes.Plan,
			fromVM module.VersionMap,
		) (module.VersionMap, error) {
			return app.ModuleManager.RunMigrations(ctx, app.Configurator(), fromVM)
		},
	)
}
