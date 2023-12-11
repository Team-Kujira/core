package types_test

import (
	"encoding/json"
	"testing"

	"github.com/Team-Kujira/core/x/oracle/types"

	"github.com/stretchr/testify/require"
)

func TestGenesisValidation(t *testing.T) {
	genState := types.DefaultGenesisState()
	require.NoError(t, types.ValidateGenesis(genState))

	genState.Params.VotePeriod = 0
	require.Error(t, types.ValidateGenesis(genState))
}

func TestGetGenesisStateFromAppState(t *testing.T) {
	defaultGenesisState := types.DefaultGenesisState()
	bz, err := json.Marshal(defaultGenesisState)
	require.Nil(t, err)

	require.NotNil(t, types.GetGenesisStateFromAppState(map[string]json.RawMessage{
		types.ModuleName: bz,
	}))
	require.NotNil(t, types.GetGenesisStateFromAppState(map[string]json.RawMessage{}))
}
