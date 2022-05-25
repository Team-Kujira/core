package scheduler_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "kujira/testutil/keeper"
	"kujira/testutil/nullify"
	"kujira/x/scheduler"
	"kujira/x/scheduler/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		HookList: []types.Hook{
			{
				Id: 0,
			},
			{
				Id: 1,
			},
		},
		HookCount: 2,
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.SchedulerKeeper(t)
	scheduler.InitGenesis(ctx, *k, genesisState)
	got := scheduler.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.ElementsMatch(t, genesisState.HookList, got.HookList)
	require.Equal(t, genesisState.HookCount, got.HookCount)
	// this line is used by starport scaffolding # genesis/test/assert
}
