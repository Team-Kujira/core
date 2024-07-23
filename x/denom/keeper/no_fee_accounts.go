package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Team-Kujira/core/x/denom/types"
)

// IsNoFeeAccount returns if an address is no fee account
func (k Keeper) IsNoFeeAccount(ctx sdk.Context, address string) bool {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.GetNoFeeAccountPrefix())
	bz := prefixStore.Get([]byte(address))
	return bz != nil
}

// SetNoFeeAccount sets an address as no fee account
func (k Keeper) SetNoFeeAccount(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.GetNoFeeAccountPrefix())
	prefixStore.Set([]byte(address), []byte(address))
}

// RemoveNoFeeAccount removes an address from no fee account list
func (k Keeper) RemoveNoFeeAccount(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	prefixStore := prefix.NewStore(store, types.GetNoFeeAccountPrefix())
	prefixStore.Delete([]byte(address))
}

// GetNoFeeAccounts returns all no fee accounts
func (k Keeper) GetNoFeeAccounts(ctx sdk.Context) []string {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetNoFeeAccountPrefix())
	defer iterator.Close()

	accounts := []string{}
	for ; iterator.Valid(); iterator.Next() {
		accounts = append(accounts, string(iterator.Value()))
	}
	return accounts
}
