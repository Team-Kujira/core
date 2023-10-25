package v1

import (
	denomtypes "github.com/Team-Kujira/core/x/denom/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

func MigrateParams(
	ctx sdk.Context,
	store sdk.KVStore,
	subspace paramtypes.Subspace,
	cdc codec.BinaryCodec,
) error {
	var denomParams denomtypes.Params

	subspace.Get(ctx, denomtypes.ParamsKey, &denomParams)

	bz := cdc.MustMarshal(&denomParams)
	store.Set(denomtypes.ParamsKey, bz)

	return nil
}
