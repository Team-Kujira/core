package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/Team-Kujira/core/x/scheduler/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) HookAll(c context.Context, req *types.QueryAllHookRequest) (*types.QueryAllHookResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var hooks []types.Hook
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	hookStore := prefix.NewStore(store, types.KeyPrefix(types.HookKey))

	pageRes, err := query.Paginate(hookStore, req.Pagination, func(key []byte, value []byte) error {
		var hook types.Hook
		if err := k.cdc.Unmarshal(value, &hook); err != nil {
			return err
		}

		hooks = append(hooks, hook)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllHookResponse{Hook: hooks, Pagination: pageRes}, nil
}

func (k Keeper) Hook(c context.Context, req *types.QueryGetHookRequest) (*types.QueryGetHookResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	hook, found := k.GetHook(ctx, req.Id)
	if !found {
		return nil, sdkerrors.ErrKeyNotFound
	}

	return &types.QueryGetHookResponse{Hook: hook}, nil
}
