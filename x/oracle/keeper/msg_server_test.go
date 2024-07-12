package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Team-Kujira/core/x/oracle/types"

	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

var (
	stakingAmt = sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)
)

func setup(t *testing.T) (TestInput, types.MsgServer) {
	input := CreateTestInput(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.VotePeriod = 1
	params.SlashWindow = 100
	input.OracleKeeper.SetParams(input.Ctx, params)
	msgServer := NewMsgServerImpl(input.OracleKeeper)

	sh := stakingkeeper.NewMsgServerImpl(&input.StakingKeeper)

	// Validator created
	_, err := sh.CreateValidator(input.Ctx, NewTestMsgCreateValidator(ValAddrs[0], ValPubKeys[0], stakingAmt))
	require.NoError(t, err)
	_, err = sh.CreateValidator(input.Ctx, NewTestMsgCreateValidator(ValAddrs[1], ValPubKeys[1], stakingAmt))
	require.NoError(t, err)
	_, err = sh.CreateValidator(input.Ctx, NewTestMsgCreateValidator(ValAddrs[2], ValPubKeys[2], stakingAmt))
	require.NoError(t, err)

	input.StakingKeeper.EndBlocker(input.Ctx)

	return input, msgServer
}

func TestMsgUpdateParams(t *testing.T) {
	input, msgServer := setup(t)

	// Test default params setting
	input.OracleKeeper.SetParams(input.Ctx, types.DefaultParams())
	params := input.OracleKeeper.GetParams(input.Ctx)
	require.NotNil(t, params)

	// Test updating with invalid authority
	_, err := msgServer.UpdateParams(input.Ctx, &types.MsgUpdateParams{
		Authority: "invalid_authority",
		Params:    &params,
	})
	require.Error(t, err)

	// Test updating with correct authority
	_, err = msgServer.UpdateParams(input.Ctx, &types.MsgUpdateParams{
		Authority: "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
		Params:    &params,
	})
	require.NoError(t, err)

	// Test updating required denoms
	params.RequiredSymbols = []types.Symbol{{Symbol: "BTC", Id: 1}}
	_, err = msgServer.UpdateParams(input.Ctx, &types.MsgUpdateParams{
		Authority: "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
		Params:    &params,
	})
	require.Error(t, err)
}

func TestMsgAddRequiredSymbols(t *testing.T) {
	input, msgServer := setup(t)

	params := input.OracleKeeper.GetParams(input.Ctx)
	require.Len(t, params.RequiredSymbols, 3)
	require.Equal(t, params.LastSymbolId, uint32(3))

	// Test with invalid authority
	_, err := msgServer.AddRequiredSymbols(input.Ctx, &types.MsgAddRequiredSymbols{
		Authority: "invalid_authority",
		Symbols:   []string{"ATOM"},
	})
	require.Error(t, err)

	// Adding a single denom
	_, err = msgServer.AddRequiredSymbols(input.Ctx, &types.MsgAddRequiredSymbols{
		Authority: "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
		Symbols:   []string{"ATOM"},
	})
	require.NoError(t, err)
	params = input.OracleKeeper.GetParams(input.Ctx)
	require.Len(t, params.RequiredSymbols, 4)
	require.Equal(t, params.LastSymbolId, uint32(4))

	// Adding two denoms
	_, err = msgServer.AddRequiredSymbols(input.Ctx, &types.MsgAddRequiredSymbols{
		Authority: "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
		Symbols:   []string{"AKT", "ARB"},
	})
	require.NoError(t, err)
	params = input.OracleKeeper.GetParams(input.Ctx)
	require.Len(t, params.RequiredSymbols, 6)
	require.Equal(t, params.LastSymbolId, uint32(6))

	// Adding already existing denom
	_, err = msgServer.AddRequiredSymbols(input.Ctx, &types.MsgAddRequiredSymbols{
		Authority: "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
		Symbols:   []string{"ATOM"},
	})
	require.Error(t, err)
}

func TestMsgRemoveRequiredSymbols(t *testing.T) {
	input, msgServer := setup(t)

	params := input.OracleKeeper.GetParams(input.Ctx)
	require.Len(t, params.RequiredSymbols, 3)
	require.Equal(t, params.LastSymbolId, uint32(3))

	// Test with invalid authority
	_, err := msgServer.RemoveRequiredSymbols(input.Ctx, &types.MsgRemoveRequiredSymbols{
		Authority: "invalid_authority",
		Symbols:   []string{"BTC"},
	})
	require.Error(t, err)

	// Removing a single denom
	_, err = msgServer.RemoveRequiredSymbols(input.Ctx, &types.MsgRemoveRequiredSymbols{
		Authority: "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
		Symbols:   []string{"BTC"},
	})
	require.NoError(t, err)
	params = input.OracleKeeper.GetParams(input.Ctx)
	require.Len(t, params.RequiredSymbols, 2)
	require.Equal(t, params.LastSymbolId, uint32(3))

	// Removing two denoms
	_, err = msgServer.RemoveRequiredSymbols(input.Ctx, &types.MsgRemoveRequiredSymbols{
		Authority: "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
		Symbols:   []string{"ETH", "USDT"},
	})
	require.NoError(t, err)
	params = input.OracleKeeper.GetParams(input.Ctx)
	require.Len(t, params.RequiredSymbols, 0)
	require.Equal(t, params.LastSymbolId, uint32(3))

	// Removing not existing denom
	_, err = msgServer.RemoveRequiredSymbols(input.Ctx, &types.MsgRemoveRequiredSymbols{
		Authority: "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
		Symbols:   []string{"BTC"},
	})
	require.NoError(t, err)
}
