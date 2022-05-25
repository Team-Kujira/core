package scheduler

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"kujira/x/scheduler/keeper"
	"kujira/x/scheduler/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// Set all the hook
	for _, elem := range genState.HookList {
		k.SetHook(ctx, elem)
	}

	// Set hook count
	k.SetHookCount(ctx, genState.HookCount)
	// this line is used by starport scaffolding # genesis/module/init
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.HookList = k.GetAllHook(ctx)
	genesis.HookCount = k.GetHookCount(ctx)
	// this line is used by starport scaffolding # genesis/module/export

	return genesis
}
