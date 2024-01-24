package keeper

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Team-Kujira/core/x/oracle/types"

	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

var (
	stakingAmt         = sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction)
	randomExchangeRate = math.LegacyNewDec(1700)
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
