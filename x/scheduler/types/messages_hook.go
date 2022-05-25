package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateHook = "create_hook"
	TypeMsgUpdateHook = "update_hook"
	TypeMsgDeleteHook = "delete_hook"
)

var _ sdk.Msg = &MsgCreateHook{}

func NewMsgCreateHook(creator string, contract string, executor string, msg []byte, frequency int64) *MsgCreateHook {
	return &MsgCreateHook{
		Creator:   creator,
		Executor:  executor,
		Contract:  contract,
		Msg:       msg,
		Frequency: frequency,
	}
}

func (msg *MsgCreateHook) Route() string {
	return RouterKey
}

func (msg *MsgCreateHook) Type() string {
	return TypeMsgCreateHook
}

func (msg *MsgCreateHook) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateHook) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateHook) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateHook{}

func NewMsgUpdateHook(creator string, id uint64, contract string, executor string, msg []byte, frequency int64) *MsgUpdateHook {
	return &MsgUpdateHook{
		Id:        id,
		Creator:   creator,
		Executor:  executor,
		Contract:  contract,
		Msg:       msg,
		Frequency: frequency,
	}
}

func (msg *MsgUpdateHook) Route() string {
	return RouterKey
}

func (msg *MsgUpdateHook) Type() string {
	return TypeMsgUpdateHook
}

func (msg *MsgUpdateHook) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateHook) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateHook) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteHook{}

func NewMsgDeleteHook(creator string, id uint64) *MsgDeleteHook {
	return &MsgDeleteHook{
		Id:      id,
		Creator: creator,
	}
}
func (msg *MsgDeleteHook) Route() string {
	return RouterKey
}

func (msg *MsgDeleteHook) Type() string {
	return TypeMsgDeleteHook
}

func (msg *MsgDeleteHook) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteHook) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteHook) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
