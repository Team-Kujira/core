package batch

import (
	"github.com/Team-Kujira/core/x/batch/keeper"
	"github.com/Team-Kujira/core/x/batch/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(_ sdk.Context, _ keeper.Keeper, _ types.GenesisState) {}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(_ sdk.Context, _ keeper.Keeper) *types.GenesisState {
	return types.DefaultGenesis()
}
