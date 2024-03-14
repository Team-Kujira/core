package main

import (
	"os"

	"cosmossdk.io/log"
	"github.com/Team-Kujira/core/app"
	"github.com/Team-Kujira/core/cmd/kujirad/cmd"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()

	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		log.NewLogger(rootCmd.OutOrStderr()).Error("failure when running app", "err", err)
		os.Exit(1)
	}
}
