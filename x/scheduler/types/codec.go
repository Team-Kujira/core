package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the necessary x/oracle interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateHook{}, "scheduler/MsgCreateHook", nil)
	cdc.RegisterConcrete(&MsgUpdateHook{}, "scheduler/MsgUpdateHook", nil)
	cdc.RegisterConcrete(&MsgDeleteHook{}, "scheduler/MsgDeleteHook", nil)
}

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateHook{}, "scheduler/MsgCreateHook", nil)
	cdc.RegisterConcrete(&MsgUpdateHook{}, "scheduler/MsgUpdateHook", nil)
	cdc.RegisterConcrete(&MsgDeleteHook{}, "scheduler/MsgDeleteHook", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
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
