package scheduler

import (
	"github.com/Team-Kujira/core/x/scheduler/keeper"
	"github.com/Team-Kujira/core/x/scheduler/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	k.SetParams(ctx, genState.Params)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	genesis.HookList = k.GetAllHook(ctx)
	genesis.HookCount = k.GetHookCount(ctx)

	return genesis
}
