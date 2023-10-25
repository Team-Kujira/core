package keeper

import (
	"testing"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/Team-Kujira/core/x/scheduler/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestCreateHook(t *testing.T) {
	keeper, ctx := SetupKeeper(t)
	ms := NewMsgServerImpl(keeper)

	authority := "kujira10d07y265gmmuvt4z0w9aw880jnsr700jt23ame"
	executor := "cosmos1fl48vsnmsdzcv85q5d2q4z5ajdha8yu34mf0eh"
	contract := "cosmos1fl48vsnmsdzcv85q5d2q4z5ajdha8yu34mf0eh"
	jsonMsg := wasmtypes.RawContractMessage(`{"foo": 123}`)
	frequency := int64(1000)
	funds := sdk.NewCoins(sdk.NewInt64Coin("atom", 100))

	_, err := ms.CreateHook(ctx, &types.MsgCreateHook{
		Authority: authority,
		Executor:  executor,
		Contract:  contract,
		Msg:       jsonMsg,
		Frequency: frequency,
		Funds:     funds,
	})
	require.NoError(t, err)

	authority = "cosmos1fl48vsnmsdzcv85q5d2q4z5ajdha8yu34mf0eh"

	_, err = ms.CreateHook(ctx, &types.MsgCreateHook{
		Authority: authority,
		Executor:  executor,
		Contract:  contract,
		Msg:       jsonMsg,
		Frequency: frequency,
		Funds:     funds,
	})
	require.Error(t, err)
}

func TestUpdateHook(t *testing.T) {
	keeper, ctx := SetupKeeper(t)
	ms := NewMsgServerImpl(keeper)

	authority := "kujira10d07y265gmmuvt4z0w9aw880jnsr700jt23ame"
	executor := "cosmos1fl48vsnmsdzcv85q5d2q4z5ajdha8yu34mf0eh"
	contract := "cosmos1fl48vsnmsdzcv85q5d2q4z5ajdha8yu34mf0eh"
	jsonMsg := wasmtypes.RawContractMessage(`{"foo": 123}`)
	frequency := int64(1000)
	funds := sdk.NewCoins(sdk.NewInt64Coin("atom", 100))

	hook := types.Hook{
		Executor:  executor,
		Contract:  contract,
		Msg:       jsonMsg,
		Frequency: frequency,
		Funds:     funds,
	}
	id := keeper.AppendHook(ctx, hook)

	_, err := ms.UpdateHook(ctx, &types.MsgUpdateHook{
		Authority: authority,
		Id:        id,
		Executor:  executor,
		Contract:  contract,
		Msg:       jsonMsg,
		Frequency: frequency,
		Funds:     funds,
	})
	require.NoError(t, err)

	authority = "cosmos1fl48vsnmsdzcv85q5d2q4z5ajdha8yu34mf0eh"

	_, err = ms.UpdateHook(ctx, &types.MsgUpdateHook{
		Authority: authority,
		Id:        id,
		Executor:  executor,
		Contract:  contract,
		Msg:       jsonMsg,
		Frequency: frequency,
		Funds:     funds,
	})
	require.Error(t, err)
}

func TestDeleteHook(t *testing.T) {
	keeper, ctx := SetupKeeper(t)
	ms := NewMsgServerImpl(keeper)

	authority := "kujira10d07y265gmmuvt4z0w9aw880jnsr700jt23ame"
	executor := "cosmos1fl48vsnmsdzcv85q5d2q4z5ajdha8yu34mf0eh"
	contract := "cosmos1fl48vsnmsdzcv85q5d2q4z5ajdha8yu34mf0eh"
	jsonMsg := wasmtypes.RawContractMessage(`{"foo": 123}`)
	frequency := int64(1000)
	funds := sdk.NewCoins(sdk.NewInt64Coin("atom", 100))

	hook := types.Hook{
		Executor:  executor,
		Contract:  contract,
		Msg:       jsonMsg,
		Frequency: frequency,
		Funds:     funds,
	}
	id := keeper.AppendHook(ctx, hook)

	_, err := ms.DeleteHook(ctx, &types.MsgDeleteHook{
		Authority: authority,
		Id:        id,
	})
	require.NoError(t, err)

	authority = "cosmos1fl48vsnmsdzcv85q5d2q4z5ajdha8yu34mf0eh"
	_, err = ms.DeleteHook(ctx, &types.MsgDeleteHook{
		Authority: authority,
		Id:        id,
	})
	require.Error(t, err)
}

// Setup the testing environment
func SetupKeeper(t *testing.T) (Keeper, sdk.Context) {
	return CreateTestKeeper(t)
}
