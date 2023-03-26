package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateDenom{}, "github.com/Team-Kujira/core/denom/create-denom", nil)
	cdc.RegisterConcrete(&MsgMint{}, "github.com/Team-Kujira/core/denom/mint", nil)
	cdc.RegisterConcrete(&MsgBurn{}, "github.com/Team-Kujira/core/denom/burn", nil)
	// cdc.RegisterConcrete(&MsgForceTransfer{}, "github.com/Team-Kujira/core/denom/force-transfer", nil)
	cdc.RegisterConcrete(&MsgChangeAdmin{}, "github.com/Team-Kujira/core/denom/change-admin", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateDenom{},
		&MsgMint{},
		&MsgBurn{},
		// &MsgForceTransfer{},
		&MsgChangeAdmin{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
