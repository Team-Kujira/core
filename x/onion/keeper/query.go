package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Team-Kujira/core/x/onion/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Sequence(c context.Context, req *types.QuerySequenceRequest) (*types.QuerySequenceResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	seq, err := k.GetSequence(ctx, req.Address)
	if err != nil {
		return nil, err
	}
	return &types.QuerySequenceResponse{Seq: seq}, nil
}
