package types

import (
	"encoding/json"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(
	params Params, rates []ExchangeRateTuple,
	missCounters []MissCounter,
) *GenesisState {
	return &GenesisState{
		Params:        params,
		ExchangeRates: rates,
		MissCounters:  missCounters,
	}
}

// DefaultGenesisState - default GenesisState used by columbus-2
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(DefaultParams(),
		[]ExchangeRateTuple{},
		[]MissCounter{})
}

// ValidateGenesis validates the oracle genesis state
func ValidateGenesis(data *GenesisState) error {
	return data.Params.Validate()
}

// GetGenesisStateFromAppState returns x/oracle GenesisState given raw application
// genesis state.
func GetGenesisStateFromAppState(appState map[string]json.RawMessage) *GenesisState {
	var genesisState GenesisState

	if appState[ModuleName] != nil {
		err := json.Unmarshal(appState[ModuleName], &genesisState)
		if err != nil {
			panic(err)
		}
	}

	return &genesisState
}
