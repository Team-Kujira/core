package wasmbinding

import (
	storetypes "cosmossdk.io/store/types"
	batchkeeper "github.com/Team-Kujira/core/x/batch/keeper"
	denomkeeper "github.com/Team-Kujira/core/x/denom/keeper"
	oraclekeeper "github.com/Team-Kujira/core/x/oracle/keeper"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	// bankkeeper "github.com/terra-money/alliance/custom/bank/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	cwicakeeper "github.com/Team-Kujira/core/x/cw-ica/keeper"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/controller/keeper"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
)

func RegisterCustomPlugins(
	bank bankkeeper.Keeper,
	oracle oraclekeeper.Keeper,
	denom denomkeeper.Keeper,
	batch batchkeeper.Keeper,
	ibc ibckeeper.Keeper,
	cwica cwicakeeper.Keeper,
	ica icacontrollerkeeper.Keeper,
	transfer ibctransferkeeper.Keeper,
	ibcStoreKey *storetypes.KVStoreKey,
) []wasmkeeper.Option {
	wasmQueryPlugin := NewQueryPlugin(bank, oracle, denom, ibc, cwica, ibcStoreKey)

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: CustomQuerier(wasmQueryPlugin),
	})

	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		CustomMessageDecorator(bank, denom, batch, cwica, ica, transfer),
	)

	return []wasmkeeper.Option{
		queryPluginOpt,
		messengerDecoratorOpt,
	}
}
