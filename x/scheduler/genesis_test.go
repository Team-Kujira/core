// I've disabled the scheduler test for now, but I think it's a good idea to have a test like this for the scheduler module. Remove ignite, and please refactor tests.
package scheduler_test

/*


import (
	"testing"

	keepertest "github.com/Team-Kujira/core/testutil/keeper"
	"github.com/Team-Kujira/core/testutil/nullify"
	"github.com/Team-Kujira/core/x/scheduler"
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
*/
