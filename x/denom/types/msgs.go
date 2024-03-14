package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// constants
const (
	TypeMsgCreateDenom   = "create_denom"
	TypeMsgMint          = "mint"
	TypeMsgBurn          = "burn"
	TypeMsgForceTransfer = "force_transfer"
	TypeMsgChangeAdmin   = "change_admin"
	TypeMsgUpdateParams  = "update_params"
)

var (
	_ sdk.Msg = &MsgCreateDenom{}
	_ sdk.Msg = &MsgUpdateParams{}
)

// NewMsgCreateDenom creates a msg to create a new denom
func NewMsgCreateDenom(sender, nonce string) *MsgCreateDenom {
	return &MsgCreateDenom{
		Sender: sender,
		Nonce:  nonce,
	}
}

func (m MsgCreateDenom) Route() string { return RouterKey }
func (m MsgCreateDenom) Type() string  { return TypeMsgCreateDenom }
func (m MsgCreateDenom) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	_, err = GetTokenDenom(m.Sender, m.Nonce)
	if err != nil {
		return errors.Wrap(ErrInvalidDenom, err.Error())
	}

	return nil
}

func (m MsgCreateDenom) GetSigners() []sdk.AccAddress {
	sender, _ := sdk.AccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgMint{}

// NewMsgMint creates a message to mint tokens
func NewMsgMint(sender string, amount sdk.Coin, recipient string) *MsgMint {
	return &MsgMint{
		Sender:    sender,
		Amount:    amount,
		Recipient: recipient,
	}
}

func (m MsgMint) Route() string { return RouterKey }
func (m MsgMint) Type() string  { return TypeMsgMint }
func (m MsgMint) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	if !m.Amount.IsValid() {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, m.Amount.String())
	}

	return nil
}

func (m MsgMint) GetSigners() []sdk.AccAddress {
	sender, _ := sdk.AccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgBurn{}

// NewMsgBurn creates a message to burn tokens
func NewMsgBurn(sender string, amount sdk.Coin) *MsgBurn {
	return &MsgBurn{
		Sender: sender,
		Amount: amount,
	}
}

func (m MsgBurn) Route() string { return RouterKey }
func (m MsgBurn) Type() string  { return TypeMsgBurn }
func (m MsgBurn) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	if !m.Amount.IsValid() {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, m.Amount.String())
	}

	return nil
}

func (m MsgBurn) GetSigners() []sdk.AccAddress {
	sender, _ := sdk.AccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{sender}
}

// var _ sdk.Msg = &MsgForceTransfer{}

// // NewMsgForceTransfer creates a transfer funds from one account to another
// func NewMsgForceTransfer(sender string, amount sdk.Coin, fromAddr, toAddr string) *MsgForceTransfer {
// 	return &MsgForceTransfer{
// 		Sender:              sender,
// 		Amount:              amount,
// 		TransferFromAddress: fromAddr,
// 		TransferToAddress:   toAddr,
// 	}
// }

// func (m MsgForceTransfer) Route() string { return RouterKey }
// func (m MsgForceTransfer) Type() string  { return TypeMsgForceTransfer }
// func (m MsgForceTransfer) ValidateBasic() error {
// 	_, err := sdk.AccAddressFromBech32(m.Sender)
// 	if err != nil {
// 		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
// 	}

// 	_, err = sdk.AccAddressFromBech32(m.TransferFromAddress)
// 	if err != nil {
// 		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid address (%s)", err)
// 	}
// 	_, err = sdk.AccAddressFromBech32(m.TransferToAddress)
// 	if err != nil {
// 		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid address (%s)", err)
// 	}

// 	if !m.Amount.IsValid() {
// 		return errors.Wrap(sdkerrors.ErrInvalidCoins, m.Amount.String())
// 	}

// 	return nil
// }

// func (m MsgForceTransfer) GetSigners() []sdk.AccAddress {
// 	sender, _ := sdk.AccAddressFromBech32(m.Sender)
// 	return []sdk.AccAddress{sender}
// }

var _ sdk.Msg = &MsgChangeAdmin{}

// NewMsgChangeAdmin creates a message to burn tokens
func NewMsgChangeAdmin(sender, denom, newAdmin string) *MsgChangeAdmin {
	return &MsgChangeAdmin{
		Sender:   sender,
		Denom:    denom,
		NewAdmin: newAdmin,
	}
}

func (m MsgChangeAdmin) Route() string { return RouterKey }
func (m MsgChangeAdmin) Type() string  { return TypeMsgChangeAdmin }
func (m MsgChangeAdmin) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.Sender)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid sender address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(m.NewAdmin)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid address (%s)", err)
	}

	_, _, err = DeconstructDenom(m.Denom)
	if err != nil {
		return err
	}

	return nil
}

func (m MsgChangeAdmin) GetSigners() []sdk.AccAddress {
	sender, _ := sdk.AccAddressFromBech32(m.Sender)
	return []sdk.AccAddress{sender}
}

// NewMsgUpdateParams creates a MsgUpdateParams instance
func NewMsgUpdateParams(params *Params) *MsgUpdateParams {
	return &MsgUpdateParams{
		Params: params,
	}
}

// Route implements sdk.Msg
func (m MsgUpdateParams) Route() string { return RouterKey }

// Type implements sdk.Msg
func (m MsgUpdateParams) Type() string { return TypeMsgUpdateParams }

// GetSigners implements sdk.Msg
func (m MsgUpdateParams) GetSigners() []sdk.AccAddress {
	operator, err := sdk.ValAddressFromBech32(m.Authority)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{sdk.AccAddress(operator)}
}

// ValidateBasic implements sdk.Msg
func (m MsgUpdateParams) ValidateBasic() error {
	_, err := sdk.ValAddressFromBech32(m.Authority)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid authority address (%s)", err)
	}

	return nil
}
