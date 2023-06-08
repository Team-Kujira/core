package main

import (
	"os"

	"github.com/Team-Kujira/core/app"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/ignite/cli/ignite/pkg/cosmoscmd" //nolint:staticcheck // TODO: remove
)

func main() {
	rootCmd, _ := cosmoscmd.NewRootCmd(
		app.Name,
		app.AccountAddressPrefix,
		app.DefaultNodeHome,
		app.Name,
		app.ModuleBasics,
		app.NewIgniteApp,
		// this line is used by starport scaffolding # root/arguments
	)
	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
