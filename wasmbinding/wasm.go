package wasmbinding

import (
	denomkeeper "github.com/Team-Kujira/core/x/denom/keeper"

	oraclekeeper "github.com/Team-Kujira/core/x/oracle/keeper"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	intertxkeeper "github.com/Team-Kujira/core/x/inter-tx/keeper"
	bankkeeper "github.com/terra-money/alliance/custom/bank/keeper"

	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
)

func RegisterCustomPlugins(
	bank bankkeeper.Keeper,
	oracle oraclekeeper.Keeper,
	denom denomkeeper.Keeper,
	intertx intertxkeeper.Keeper,
	ica icacontrollerkeeper.Keeper,
) []wasmkeeper.Option {
	wasmQueryPlugin := NewQueryPlugin(bank, oracle, denom, intertx)

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: CustomQuerier(wasmQueryPlugin),
	})

	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		CustomMessageDecorator(bank, denom, intertx, ica),
	)

	return []wasmkeeper.Option{
		queryPluginOpt,
		messengerDecoratorOpt,
	}
}
