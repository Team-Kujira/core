package types

import (
	"testing"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateHook_ValidateBasic(t *testing.T) {
	jsonMsg := wasmtypes.RawContractMessage(`{"foo": 123}`)

	msg := NewMsgCreateHook(
		"cosmos1fl48vsnmsdzcv85q5d2q4z5ajdha8yu34mf0eh",
		"cosmos1fl48vsnmsdzcv85q5d2q4z5ajdha8yu34mf0eh",
		"cosmos1fl48vsnmsdzcv85q5d2q4z5ajdha8yu34mf0eh",
		jsonMsg,
		1000,
		sdk.NewCoins(sdk.NewInt64Coin("atom", 100)),
	)
	require.NoError(t, msg.ValidateBasic())
}

func TestMsgCreateHook_ValidateBasic_Error(t *testing.T) {
	msg := NewMsgCreateHook(
		"invalidAddress",
		"invalidAddress",
		"invalidAddress",
		wasmtypes.RawContractMessage{},
		1000,
		sdk.NewCoins(sdk.NewInt64Coin("atom", 100)),
	)
	require.Error(t, msg.ValidateBasic())
}

func TestMsgUpdateHook_ValidateBasic(t *testing.T) {
	jsonMsg := wasmtypes.RawContractMessage(`{"foo": 123}`)

	msg := NewMsgUpdateHook(
		"cosmos1fl48vsnmsdzcv85q5d2q4z5ajdha8yu34mf0eh",
		1,
		"cosmos1fl48vsnmsdzcv85q5d2q4z5ajdha8yu34mf0eh",
		"cosmos1fl48vsnmsdzcv85q5d2q4z5ajdha8yu34mf0eh",
		jsonMsg,
		1000,
		sdk.NewCoins(sdk.NewInt64Coin("atom", 100)),
	)
	require.NoError(t, msg.ValidateBasic())
}

func TestMsgUpdateHook_ValidateBasic_Error(t *testing.T) {
	msg := NewMsgUpdateHook(
		"invalidAddress",
		0,
		"invalidAddress",
		"invalidAddress",
		wasmtypes.RawContractMessage{},
		1000,
		sdk.NewCoins(sdk.NewInt64Coin("atom", 100)),
	)
	require.Error(t, msg.ValidateBasic())
}

func TestMsgDeleteHook_ValidateBasic(t *testing.T) {
	msg := NewMsgDeleteHook(
		"cosmos1fl48vsnmsdzcv85q5d2q4z5ajdha8yu34mf0eh",
		1,
	)
	require.NoError(t, msg.ValidateBasic())
}

func TestMsgDeleteHook_ValidateBasic_Error(t *testing.T) {
	msg := NewMsgDeleteHook(
		"invalidAddress",
		0,
	)
	require.Error(t, msg.ValidateBasic())
}
