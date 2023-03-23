// I've disabled the scheduler test for now, but I think it's a good idea to have a test like this for the scheduler module. Remove ignite, and please refactor tests.

package keeper_test

/*
import (
	"testing"

	testkeeper "github.com/Team-Kujira/core/testutil/keeper"
	"github.com/Team-Kujira/core/x/scheduler/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestParamsQuery(t *testing.T) {
	keeper, ctx := testkeeper.SchedulerKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	params := types.DefaultParams()
	keeper.SetParams(ctx, params)

	response, err := keeper.Params(wctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, &types.QueryParamsResponse{Params: params}, response)
}
*/
