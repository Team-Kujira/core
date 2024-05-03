package keeper

import (
	"cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"
	"github.com/Team-Kujira/core/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetExchangeRate gets the consensus exchange rate of the denom asset from the store.
func (k Keeper) GetHistoricalExchangeRate(ctx sdk.Context, epoch, denom string, timestamp int64) (types.HistoricalExchangeRate, error) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetHistoricalExchangeRateKey(epoch, denom, timestamp))
	if b == nil {
		return types.HistoricalExchangeRate{}, errors.Wrap(types.ErrUnknownHistoricalExchangeRate, denom)
	}

	ph := types.HistoricalExchangeRate{}
	k.cdc.MustUnmarshal(b, &ph)
	return ph, nil
}

// SetExchangeRate sets the consensus exchange rate of the denom asset to the store.
func (k Keeper) SetHistoricalExchangeRate(ctx sdk.Context, history types.HistoricalExchangeRate) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&history)
	store.Set(types.GetHistoricalExchangeRateKey(history.Epoch, history.Denom, history.Timestamp), bz)
}

// DeleteExchangeRate deletes the consensus exchange rate of the denom asset from the store.
func (k Keeper) DeleteHistoricalExchangeRate(ctx sdk.Context, history types.HistoricalExchangeRate) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetHistoricalExchangeRateKey(history.Epoch, history.Denom, history.Timestamp))
}

// IterateHistoricalExchangeRate iterates over historical exchange rates in the store
func (k Keeper) IterateHistoricalExchangeRate(ctx sdk.Context, handler func(history types.HistoricalExchangeRate) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := storetypes.KVStorePrefixIterator(store, types.HistoricalExchangeRateKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		ph := types.HistoricalExchangeRate{}
		k.cdc.MustUnmarshal(iter.Value(), &ph)
		if handler(ph) {
			break
		}
	}
}

// IterateHistoricalExchangeRateByEpochDenom iterates over historical exchange rates by epoch and denom in the store
func (k Keeper) IterateHistoricalExchangeRateByEpochDenom(ctx sdk.Context, epoch, denom string, handler func(history types.HistoricalExchangeRate) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := storetypes.KVStorePrefixIterator(store, types.GetHistoricalExchangeRatePrefix(epoch, denom))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		historicalRate := types.HistoricalExchangeRate{}
		k.cdc.MustUnmarshal(iter.Value(), &historicalRate)
		if handler(historicalRate) {
			break
		}
	}
}

func (k Keeper) LatestHistoricalExchangeRateByEpochDenom(ctx sdk.Context, epoch, denom string) types.HistoricalExchangeRate {
	store := ctx.KVStore(k.storeKey)
	iter := storetypes.KVStoreReversePrefixIterator(store, types.GetHistoricalExchangeRatePrefix(epoch, denom))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		historicalRate := types.HistoricalExchangeRate{}
		k.cdc.MustUnmarshal(iter.Value(), &historicalRate)
		return historicalRate
	}
	return types.HistoricalExchangeRate{}
}

func (k Keeper) OldestHistoricalExchangeRateByEpochDenom(ctx sdk.Context, epoch, denom string) types.HistoricalExchangeRate {
	store := ctx.KVStore(k.storeKey)
	iter := storetypes.KVStorePrefixIterator(store, types.GetHistoricalExchangeRatePrefix(epoch, denom))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		historicalRate := types.HistoricalExchangeRate{}
		k.cdc.MustUnmarshal(iter.Value(), &historicalRate)
		return historicalRate
	}
	return types.HistoricalExchangeRate{}
}
