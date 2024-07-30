package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// RegisterLegacyAminoCodec registers the necessary x/oracle interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateHook{}, "scheduler/MsgCreateHook", nil)
	cdc.RegisterConcrete(&MsgUpdateHook{}, "scheduler/MsgUpdateHook", nil)
	cdc.RegisterConcrete(&MsgDeleteHook{}, "scheduler/MsgDeleteHook", nil)
}

func RegisterCodec(cdc *codec.LegacyAmino) {
	// legacy proposals
	cdc.RegisterConcrete(&CreateHookProposal{}, "scheduler/CreateHookProposal", nil)
	cdc.RegisterConcrete(&UpdateHookProposal{}, "scheduler/UpdateHookProposal", nil)
	cdc.RegisterConcrete(&DeleteHookProposal{}, "scheduler/DeleteHookProposal", nil)

	cdc.RegisterConcrete(&MsgCreateHook{}, "scheduler/MsgCreateHook", nil)
	cdc.RegisterConcrete(&MsgUpdateHook{}, "scheduler/MsgUpdateHook", nil)
	cdc.RegisterConcrete(&MsgDeleteHook{}, "scheduler/MsgDeleteHook", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	// legacy proposals
	registry.RegisterImplementations((*govtypes.Content)(nil),
		&CreateHookProposal{},
		&UpdateHookProposal{},
		&DeleteHookProposal{},
	)

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateHook{},
		&MsgUpdateHook{},
		&MsgDeleteHook{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
