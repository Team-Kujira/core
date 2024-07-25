package oracle_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/Team-Kujira/core/x/oracle"
	"github.com/Team-Kujira/core/x/oracle/keeper"
	"github.com/Team-Kujira/core/x/oracle/types"
)

func TestExportInitGenesis(t *testing.T) {
	input, _ := setup(t)

	input.OracleKeeper.SetExchangeRate(input.Ctx, "denom", math.LegacyNewDec(123))
	input.OracleKeeper.SetMissCounter(input.Ctx, keeper.ValAddrs[0], 10)
	genesis := oracle.ExportGenesis(input.Ctx, input.OracleKeeper)

	newInput := keeper.CreateTestInput(t)
	oracle.InitGenesis(newInput.Ctx, newInput.OracleKeeper, genesis)
	newGenesis := oracle.ExportGenesis(newInput.Ctx, newInput.OracleKeeper)

	require.Equal(t, genesis, newGenesis)
}

func TestInitGenesis(t *testing.T) {
	input, _ := setup(t)
	genesis := types.DefaultGenesisState()
	require.NotPanics(t, func() {
		oracle.InitGenesis(input.Ctx, input.OracleKeeper, genesis)
	})

	genesis.MissCounters = []types.MissCounter{
		{
			ValidatorAddress: "invalid",
			MissCounter:      10,
		},
	}

	require.Panics(t, func() {
		oracle.InitGenesis(input.Ctx, input.OracleKeeper, genesis)
	})

	genesis.MissCounters = []types.MissCounter{
		{
			ValidatorAddress: keeper.ValAddrs[0].String(),
			MissCounter:      10,
		},
	}

	require.NotPanics(t, func() {
		oracle.InitGenesis(input.Ctx, input.OracleKeeper, genesis)
	})
}
