package keeper

import (
	"encoding/binary"

	"github.com/Team-Kujira/core/x/scheduler/types"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetHookCount get the total number of hook
func (k Keeper) GetHookCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.HookCountKey)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// SetHookCount set the total number of hook
func (k Keeper) SetHookCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(types.HookCountKey)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendHook appends a hook in the store with a new id and update the count
func (k Keeper) AppendHook(
	ctx sdk.Context,
	hook types.Hook,
) uint64 {
	// Create the hook
	count := k.GetHookCount(ctx)

	// Set the ID of the appended value
	hook.Id = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.HookKey))
	appendedValue := k.cdc.MustMarshal(&hook)
	store.Set(GetHookIDBytes(hook.Id), appendedValue)

	// Update hook count
	k.SetHookCount(ctx, count+1)

	return count
}

// SetHook set a specific hook in the store
func (k Keeper) SetHook(ctx sdk.Context, hook types.Hook) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.HookKey))
	b := k.cdc.MustMarshal(&hook)
	store.Set(GetHookIDBytes(hook.Id), b)
}

// GetHook returns a hook from its id
func (k Keeper) GetHook(ctx sdk.Context, id uint64) (val types.Hook, found bool) {
	parent := ctx.KVStore(k.storeKey)
	store := prefix.NewStore(
		parent,
		types.KeyPrefix(types.HookKey),
	)
	b := store.Get(GetHookIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveHook removes a hook from the store
func (k Keeper) RemoveHook(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.HookKey))
	store.Delete(GetHookIDBytes(id))
}

// GetAllHook returns all hook
func (k Keeper) GetAllHook(ctx sdk.Context) (list []types.Hook) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.HookKey))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Hook
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetHookIDBytes returns the byte representation of the ID
func GetHookIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetHookIDFromBytes returns ID in uint64 format from a byte array
func GetHookIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
