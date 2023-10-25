package keeper

import (
	"context"

	"github.com/Team-Kujira/core/x/scheduler/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(_c context.Context, _req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	return &types.QueryParamsResponse{}, nil
}
