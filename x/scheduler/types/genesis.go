package types

import (
	"fmt"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		HookList: []Hook{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// Check for duplicated ID in hook
	hookIdMap := make(map[uint64]bool)
	hookCount := gs.GetHookCount()
	for _, elem := range gs.HookList {
		if _, ok := hookIdMap[elem.Id]; ok {
			return fmt.Errorf("duplicated id for hook")
		}
		if elem.Id >= hookCount {
			return fmt.Errorf("hook id should be lower or equal than the last id")
		}
		hookIdMap[elem.Id] = true
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
