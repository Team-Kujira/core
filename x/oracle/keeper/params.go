package keeper

import (
	"github.com/Team-Kujira/core/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// VotePeriod returns the number of blocks during which voting takes place.
func (k Keeper) VotePeriod(ctx sdk.Context) (res uint64) {
	return k.GetParams(ctx).VotePeriod
}

// VoteThreshold returns the minimum percentage of votes that must be received for a ballot to pass.
func (k Keeper) VoteThreshold(ctx sdk.Context) (res sdk.Dec) {
	return k.GetParams(ctx).VoteThreshold

}

// MaxDeviation returns the ratio of allowable exchange rate error that a validator can be rewared
func (k Keeper) MaxDeviation(ctx sdk.Context) (res sdk.Dec) {
	return k.GetParams(ctx).MaxDeviation
}


// RequiredDenoms returns the denom list that can be activated
func (k Keeper) RequiredDenoms(ctx sdk.Context) (res types.DenomList) {
	return k.GetParams(ctx).RequiredDenoms
}

// SetRequiredDenoms store new required denoms to param store
// this function is only for test purpose
func (k Keeper) SetRequiredDenoms(ctx sdk.Context, denoms types.DenomList) {
	params := k.GetParams(ctx)
	params.RequiredDenoms = denoms
	k.SetParams(ctx, params)
}

// SlashFraction returns oracle voting penalty rate
func (k Keeper) SlashFraction(ctx sdk.Context) (res sdk.Dec) {
	return k.GetParams(ctx).SlashFraction

}

// SlashWindow returns # of vote period for oracle slashing
func (k Keeper) SlashWindow(ctx sdk.Context) (res uint64) {
	return k.GetParams(ctx).SlashWindow

}

// MinValidPerWindow returns oracle slashing threshold
func (k Keeper) MinValidPerWindow(ctx sdk.Context) (res sdk.Dec) {
	return k.GetParams(ctx).MinValidPerWindow

}

// GetParams returns the total set of oracle parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return params
	}

	k.cdc.MustUnmarshal(bz, &params)
	return params
}

// SetParams sets the total set of oracle parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)
	return nil
}
