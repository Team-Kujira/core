package wasmbinding

import (
	denomkeeper "github.com/Team-Kujira/core/x/denom/keeper"
	oraclekeeper "github.com/Team-Kujira/core/x/oracle/keeper"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	cwicakeeper "github.com/Team-Kujira/core/x/cw-ica/keeper"
	bankkeeper "github.com/terra-money/alliance/custom/bank/keeper"
	alliancekeeper "github.com/terra-money/alliance/x/alliance/keeper"

	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
)

func RegisterCustomPlugins(
	bank bankkeeper.Keeper,
	oracle oraclekeeper.Keeper,
	denom denomkeeper.Keeper,
	ibc ibckeeper.Keeper,
	cwica cwicakeeper.Keeper,
	ica icacontrollerkeeper.Keeper,
	allianceKeeper alliancekeeper.Keeper,
	ibcStoreKey *storetypes.KVStoreKey,
) []wasmkeeper.Option {
	wasmQueryPlugin := NewQueryPlugin(bank, oracle, denom, ibc, cwica, allianceKeeper, ibcStoreKey)

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: CustomQuerier(wasmQueryPlugin),
	})

	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		CustomMessageDecorator(bank, denom, cwica, ica),
	)

	return []wasmkeeper.Option{
		queryPluginOpt,
		messengerDecoratorOpt,
	}
}
