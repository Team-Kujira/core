package oracle_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/Team-Kujira/core/x/oracle/keeper"
	"github.com/Team-Kujira/core/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

var (
	uSDRAmt    = math.NewInt(1005 * types.MicroUnit)
	stakingAmt = sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)

	randomExchangeRate        = math.LegacyNewDec(1700)
	anotherRandomExchangeRate = math.LegacyNewDecWithPrec(4882, 2) // swap rate
)

func setupWithSmallVotingPower(t *testing.T) (keeper.TestInput, types.MsgServer) {
	input := keeper.CreateTestInput(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.VotePeriod = 1
	params.SlashWindow = 100
	input.OracleKeeper.SetParams(input.Ctx, params)
	h := keeper.NewMsgServerImpl(input.OracleKeeper)

	sh := stakingkeeper.NewMsgServerImpl(&input.StakingKeeper)
	_, err := sh.CreateValidator(input.Ctx,
		keeper.NewTestMsgCreateValidator(keeper.ValAddrs[0],
			keeper.ValPubKeys[0],
			sdk.TokensFromConsensusPower(1, sdk.DefaultPowerReduction)),
	)

	require.NoError(t, err)

	input.StakingKeeper.EndBlocker(input.Ctx)

	return input, h
}

func setup(t *testing.T) (keeper.TestInput, types.MsgServer) {
	input := keeper.CreateTestInput(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.VotePeriod = 1
	params.SlashWindow = 100
	params.RequiredDenoms = []string{types.TestDenomA, types.TestDenomC, types.TestDenomD}
	input.OracleKeeper.SetParams(input.Ctx, params)
	h := keeper.NewMsgServerImpl(input.OracleKeeper)

	sh := stakingkeeper.NewMsgServerImpl(&input.StakingKeeper)
	// Validator created
	_, err := sh.CreateValidator(input.Ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[0], keeper.ValPubKeys[0], stakingAmt))
	require.NoError(t, err)
	_, err = sh.CreateValidator(input.Ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[1], keeper.ValPubKeys[1], stakingAmt))
	require.NoError(t, err)
	_, err = sh.CreateValidator(input.Ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[2], keeper.ValPubKeys[2], stakingAmt))
	require.NoError(t, err)
	input.StakingKeeper.EndBlocker(input.Ctx)

	return input, h
}

func setupVal5(t *testing.T) (keeper.TestInput, types.MsgServer) {
	input := keeper.CreateTestInput(t)
	params := input.OracleKeeper.GetParams(input.Ctx)
	params.VotePeriod = 1
	params.SlashWindow = 100
	input.OracleKeeper.SetParams(input.Ctx, params)
	h := keeper.NewMsgServerImpl(input.OracleKeeper)

	sh := stakingkeeper.NewMsgServerImpl(&input.StakingKeeper)

	// Validator created
	_, err := sh.CreateValidator(input.Ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[0], keeper.ValPubKeys[0], stakingAmt))
	require.NoError(t, err)
	_, err = sh.CreateValidator(input.Ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[1], keeper.ValPubKeys[1], stakingAmt))
	require.NoError(t, err)
	_, err = sh.CreateValidator(input.Ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[2], keeper.ValPubKeys[2], stakingAmt))
	require.NoError(t, err)
	_, err = sh.CreateValidator(input.Ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[3], keeper.ValPubKeys[3], stakingAmt))
	require.NoError(t, err)
	_, err = sh.CreateValidator(input.Ctx, keeper.NewTestMsgCreateValidator(keeper.ValAddrs[4], keeper.ValPubKeys[4], stakingAmt))
	require.NoError(t, err)
	input.StakingKeeper.EndBlocker(input.Ctx)

	return input, h
}
