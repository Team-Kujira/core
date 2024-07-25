package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = &MsgAddRequiredSymbols{}
	_ sdk.Msg = &MsgRemoveRequiredSymbols{}
	_ sdk.Msg = &MsgUpdateParams{}
)

// oracle message types
const (
	TypeMsgAddRequiredSymbol    = "add_price"
	TypeMsgRemoveRequiredSymbol = "remove_price"
	TypeMsgUpdateParams         = "update_params"
)

//-------------------------------------------------
//-------------------------------------------------

// NewMsgAddRequiredSymbol creates a MsgAddRequiredSymbol instance
func NewMsgAddRequiredSymbol(symbols []string) *MsgAddRequiredSymbols {
	return &MsgAddRequiredSymbols{
		Symbols: symbols,
	}
}

// Route implements sdk.Msg
func (msg MsgAddRequiredSymbols) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgAddRequiredSymbols) Type() string { return TypeMsgAddRequiredSymbol }

// GetSigners implements sdk.Msg
func (msg MsgAddRequiredSymbols) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic implements sdk.Msg
func (msg MsgAddRequiredSymbols) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid authority address (%s)", err)
	}

	return nil
}

// NewMsgRemoveRequiredSymbol creates a MsgRemoveRequiredSymbol instance
func NewMsgRemoveRequiredSymbol(symbols []string) *MsgRemoveRequiredSymbols {
	return &MsgRemoveRequiredSymbols{
		Symbols: symbols,
	}
}

// Route implements sdk.Msg
func (msg MsgRemoveRequiredSymbols) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgRemoveRequiredSymbols) Type() string { return TypeMsgRemoveRequiredSymbol }

// GetSigners implements sdk.Msg
func (msg MsgRemoveRequiredSymbols) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic implements sdk.Msg
func (msg MsgRemoveRequiredSymbols) ValidateBasic() error {
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

	return msg.Params.Validate()
}
