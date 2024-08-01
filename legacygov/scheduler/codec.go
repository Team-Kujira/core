package scheduler

import (
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// legacy proposals
	registry.RegisterImplementations((*govtypes.Content)(nil),
		&CreateHookProposal{},
		&UpdateHookProposal{},
		&DeleteHookProposal{},
	)
}
