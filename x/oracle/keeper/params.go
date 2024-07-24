package keeper

import (
	"cosmossdk.io/math"
	"github.com/Team-Kujira/core/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// VoteThreshold returns the minimum percentage of votes that must be received for a ballot to pass.
func (k Keeper) VoteThreshold(ctx sdk.Context) (res math.LegacyDec) {
	return k.GetParams(ctx).VoteThreshold
}

// MaxDeviation returns the ratio of allowable exchange rate error that a validator can be rewared
func (k Keeper) MaxDeviation(ctx sdk.Context) (res math.LegacyDec) {
	return k.GetParams(ctx).MaxDeviation
}

// RequiredSymbols returns the denom list that can be activated
func (k Keeper) RequiredSymbols(ctx sdk.Context) (res []types.Symbol) {
	return k.GetParams(ctx).RequiredSymbols
}

// // SetRequiredSymbols store new required denoms to param store
// // this function is only for test purpose
// func (k Keeper) SetRequiredSymbols(ctx sdk.Context, denoms []string) {
// 	params := k.GetParams(ctx)
// 	params.RequiredSymbols = denoms
// 	err := k.SetParams(ctx, params)
// 	if err != nil {
// 		return
// 	}
// }

// SlashFraction returns oracle voting penalty rate
func (k Keeper) SlashFraction(ctx sdk.Context) (res math.LegacyDec) {
	return k.GetParams(ctx).SlashFraction
}

// SlashWindow returns # of vote period for oracle slashing
func (k Keeper) SlashWindow(ctx sdk.Context) (res uint64) {
	return k.GetParams(ctx).SlashWindow
}

// MinValidPerWindow returns oracle slashing threshold
func (k Keeper) MinValidPerWindow(ctx sdk.Context) (res math.LegacyDec) {
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
