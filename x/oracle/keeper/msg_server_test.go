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
	params.RequiredDenoms = []types.Denom{{Denom: "BTC", Id: 1}}
	_, err = msgServer.UpdateParams(input.Ctx, &types.MsgUpdateParams{
		Authority: "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
		Params:    &params,
	})
	require.Error(t, err)
}

func TestMsgAddRequiredDenoms(t *testing.T) {
	input, msgServer := setup(t)

	params := input.OracleKeeper.GetParams(input.Ctx)
	require.Len(t, params.RequiredDenoms, 3)
	require.Equal(t, params.LastDenomId, uint32(3))

	// Test with invalid authority
	_, err := msgServer.AddRequiredDenoms(input.Ctx, &types.MsgAddRequiredDenoms{
		Authority: "invalid_authority",
		Symbols:   []string{"ATOM"},
	})
	require.Error(t, err)

	// Adding a single denom
	_, err = msgServer.AddRequiredDenoms(input.Ctx, &types.MsgAddRequiredDenoms{
		Authority: "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
		Symbols:   []string{"ATOM"},
	})
	require.NoError(t, err)
	params = input.OracleKeeper.GetParams(input.Ctx)
	require.Len(t, params.RequiredDenoms, 4)
	require.Equal(t, params.LastDenomId, uint32(4))

	// Adding two denoms
	_, err = msgServer.AddRequiredDenoms(input.Ctx, &types.MsgAddRequiredDenoms{
		Authority: "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
		Symbols:   []string{"AKT", "ARB"},
	})
	require.NoError(t, err)
	params = input.OracleKeeper.GetParams(input.Ctx)
	require.Len(t, params.RequiredDenoms, 6)
	require.Equal(t, params.LastDenomId, uint32(6))

	// Adding already existing denom
	_, err = msgServer.AddRequiredDenoms(input.Ctx, &types.MsgAddRequiredDenoms{
		Authority: "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
		Symbols:   []string{"ATOM"},
	})
	require.Error(t, err)
}

func TestMsgRemoveRequiredDenoms(t *testing.T) {
	input, msgServer := setup(t)

	params := input.OracleKeeper.GetParams(input.Ctx)
	require.Len(t, params.RequiredDenoms, 3)
	require.Equal(t, params.LastDenomId, uint32(3))

	// Test with invalid authority
	_, err := msgServer.RemoveRequiredDenoms(input.Ctx, &types.MsgRemoveRequiredDenoms{
		Authority: "invalid_authority",
		Symbols:   []string{"BTC"},
	})
	require.Error(t, err)

	// Removing a single denom
	_, err = msgServer.RemoveRequiredDenoms(input.Ctx, &types.MsgRemoveRequiredDenoms{
		Authority: "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
		Symbols:   []string{"BTC"},
	})
	require.NoError(t, err)
	params = input.OracleKeeper.GetParams(input.Ctx)
	require.Len(t, params.RequiredDenoms, 2)
	require.Equal(t, params.LastDenomId, uint32(3))

	// Removing two denoms
	_, err = msgServer.RemoveRequiredDenoms(input.Ctx, &types.MsgRemoveRequiredDenoms{
		Authority: "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
		Symbols:   []string{"ETH", "USDT"},
	})
	require.NoError(t, err)
	params = input.OracleKeeper.GetParams(input.Ctx)
	require.Len(t, params.RequiredDenoms, 0)
	require.Equal(t, params.LastDenomId, uint32(3))

	// Removing not existing denom
	_, err = msgServer.RemoveRequiredDenoms(input.Ctx, &types.MsgRemoveRequiredDenoms{
		Authority: "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
		Symbols:   []string{"BTC"},
	})
	require.NoError(t, err)
}
