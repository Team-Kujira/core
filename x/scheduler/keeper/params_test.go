package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "kujira/testutil/keeper"
	"kujira/x/scheduler/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.SchedulerKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
