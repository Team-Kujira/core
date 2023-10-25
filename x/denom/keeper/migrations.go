package keeper

import (
	"github.com/CosmWasm/wasmd/x/wasm/exported"
	v1 "github.com/Team-Kujira/core/x/denom/migrations/v1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Migrator is a struct for handling in-place store migrations.
type Migrator struct {
	keeper         Keeper
	legacySubspace exported.Subspace
	// queryServer    grpc.Server
}

// NewMigrator returns a new Migrator.
func NewMigrator(keeper Keeper, subspace paramstypes.Subspace) Migrator {
	return Migrator{keeper: keeper, legacySubspace: subspace}
}

// Migrate1to2 migrates from version 1 to 2.
func (m Migrator) Migrate1to2(ctx sdk.Context) error {
	store := ctx.KVStore(m.keeper.storeKey)
	return v1.MigrateParams(ctx, store, m.keeper.paramSpace, m.keeper.cdc)
}
