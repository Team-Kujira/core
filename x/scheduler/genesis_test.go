package scheduler_test

import (
	"testing"

	"github.com/Team-Kujira/core/x/scheduler"
	"github.com/Team-Kujira/core/x/scheduler/keeper"
	"github.com/Team-Kujira/core/x/scheduler/types"
	"github.com/stretchr/testify/require"
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
	}

	k, ctx := keeper.CreateTestKeeper(t)
	scheduler.InitGenesis(ctx, k, genesisState)
	got := scheduler.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	require.ElementsMatch(t, genesisState.HookList, got.HookList)
	require.Equal(t, genesisState.HookCount, got.HookCount)
}
