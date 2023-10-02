package wasmbinding

import (
	denomkeeper "github.com/Team-Kujira/core/x/denom/keeper"
	oraclekeeper "github.com/Team-Kujira/core/x/oracle/keeper"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	ibckeeper "github.com/cosmos/ibc-go/v6/modules/core/keeper"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/terra-money/alliance/custom/bank/keeper"
)

func RegisterCustomPlugins(
	bank bankkeeper.Keeper,
	oracle oraclekeeper.Keeper,
	denom denomkeeper.Keeper,
	ibc ibckeeper.Keeper,
	auth authkeeper.AccountKeeper,
	ibcStoreKey *storetypes.KVStoreKey,
) []wasmkeeper.Option {
	wasmQueryPlugin := NewQueryPlugin(bank, oracle, denom, ibc, ibcStoreKey)

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: CustomQuerier(wasmQueryPlugin),
	})

	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		CustomMessageDecorator(bank, denom, auth),
	)

	return []wasm.Option{
		queryPluginOpt,
		messengerDecoratorOpt,
	}
}
