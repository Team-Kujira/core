package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&CreateHookProposal{}, "scheduler/CreateHookProposal", nil)
	cdc.RegisterConcrete(&UpdateHookProposal{}, "scheduler/UpdateHookProposal", nil)
	cdc.RegisterConcrete(&DeleteHookProposal{}, "scheduler/DeleteHookProposal", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*govtypes.Content)(nil),
		&CreateHookProposal{},
		&UpdateHookProposal{},
		&DeleteHookProposal{},
	)
	// this line is used by starport scaffolding # 3

}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
