package keeper

import (
	"context"

	"github.com/Team-Kujira/core/x/scheduler/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(_ context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	return &types.QueryParamsResponse{}, nil
}
