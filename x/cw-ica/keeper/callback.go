package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Team-Kujira/core/x/cw-ica/types"
)

// SetCallbackData set a specific callbackData in the store from its index
func (k Keeper) SetCallbackData(ctx sdk.Context, data types.CallbackData) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CallbackDataKeyPrefix))
	b := k.Codec.MustMarshal(&data)
	store.Set(types.CallbackDataKey(types.PacketID(data.PortId, data.ChannelId, data.Sequence)), b)
}

// GetCallbackData returns a callbackData from its index
func (k Keeper) GetCallbackData(
	ctx sdk.Context,
	callbackKey string,
) (val types.CallbackData, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CallbackDataKeyPrefix))

	b := store.Get(types.CallbackDataKey(callbackKey))
	if b == nil {
		return val, false
	}

	k.Codec.MustUnmarshal(b, &val)
	return val, true
}

// RemoveCallbackData removes a callbackData from the store
func (k Keeper) RemoveCallbackData(
	ctx sdk.Context,
	callbackKey string,
) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CallbackDataKeyPrefix))
	store.Delete(types.CallbackDataKey(callbackKey))
}

// GetAllCallbackData returns all callbackData
func (k Keeper) GetAllCallbackData(ctx sdk.Context) (list []types.CallbackData) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CallbackDataKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.CallbackData
		k.Codec.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
