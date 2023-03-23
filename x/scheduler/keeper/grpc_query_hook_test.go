// I've removed the test for now, but I think it's a good idea to have a test like this for the scheduler module. Remove ignite, and please refactor tests.
package keeper_test

/*
import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/Team-Kujira/core/testutil/keeper"
	"github.com/Team-Kujira/core/testutil/nullify"
	"github.com/Team-Kujira/core/x/scheduler/types"
)

func TestHookQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.SchedulerKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNHook(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetHookRequest
		response *types.QueryGetHookResponse
		err      error
	}{
		{
			desc:     "First",
			request:  &types.QueryGetHookRequest{Id: msgs[0].Id},
			response: &types.QueryGetHookResponse{Hook: msgs[0]},
		},
		{
			desc:     "Second",
			request:  &types.QueryGetHookRequest{Id: msgs[1].Id},
			response: &types.QueryGetHookResponse{Hook: msgs[1]},
		},
		{
			desc:    "KeyNotFound",
			request: &types.QueryGetHookRequest{Id: uint64(len(msgs))},
			err:     sdkerrors.ErrKeyNotFound,
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.Hook(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

func TestHookQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.SchedulerKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNHook(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllHookRequest {
		return &types.QueryAllHookRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.HookAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Hook), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Hook),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.HookAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.Hook), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.Hook),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.HookAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.Hook),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.HookAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
*/
