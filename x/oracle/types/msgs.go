package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = &MsgAddRequiredDenom{}
	_ sdk.Msg = &MsgRemoveRequiredDenom{}
	_ sdk.Msg = &MsgUpdateParams{}
)

// oracle message types
const (
	TypeMsgAddRequiredDenom    = "add_price"
	TypeMsgRemoveRequiredDenom = "remove_price"
	TypeMsgUpdateParams        = "update_params"
)

//-------------------------------------------------
//-------------------------------------------------

// NewMsgAddRequiredDenom creates a MsgAddRequiredDenom instance
func NewMsgAddRequiredDenom(symbol string) *MsgAddRequiredDenom {
	return &MsgAddRequiredDenom{
		Symbol: symbol,
	}
}

// Route implements sdk.Msg
func (msg MsgAddRequiredDenom) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgAddRequiredDenom) Type() string { return TypeMsgAddRequiredDenom }

// GetSigners implements sdk.Msg
func (msg MsgAddRequiredDenom) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic implements sdk.Msg
func (msg MsgAddRequiredDenom) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid authority address (%s)", err)
	}

	return nil
}

// NewMsgRemoveRequiredDenom creates a MsgRemoveRequiredDenom instance
func NewMsgRemoveRequiredDenom(symbol string) *MsgRemoveRequiredDenom {
	return &MsgRemoveRequiredDenom{
		Symbol: symbol,
	}
}

// Route implements sdk.Msg
func (msg MsgRemoveRequiredDenom) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgRemoveRequiredDenom) Type() string { return TypeMsgRemoveRequiredDenom }

// GetSigners implements sdk.Msg
func (msg MsgRemoveRequiredDenom) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic implements sdk.Msg
func (msg MsgRemoveRequiredDenom) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid authority address (%s)", err)
	}

	return nil
}

// NewMsgUpdateParams creates a MsgUpdateParams instance
func NewMsgUpdateParams(params *Params) *MsgUpdateParams {
	return &MsgUpdateParams{
		Params: params,
	}
}

// Route implements sdk.Msg
func (msg MsgUpdateParams) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgUpdateParams) Type() string { return TypeMsgUpdateParams }

// GetSigners implements sdk.Msg
func (msg MsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic implements sdk.Msg
func (msg MsgUpdateParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid authority address (%s)", err)
	}

	return nil
}
