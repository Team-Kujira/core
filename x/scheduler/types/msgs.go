package types

import (
	"cosmossdk.io/errors"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type ProposalType string

var (
	_ sdk.Msg = &MsgCreateHook{}
	_ sdk.Msg = &MsgUpdateHook{}
	_ sdk.Msg = &MsgDeleteHook{}
)

const (
	TypeMsgCreateHook = "create_hook"
	TypeMsgUpdateHook = "update_hook"
	TypeMsgDeleteHook = "delete_hook"
)

func NewMsgCreateHook(
	authority string,
	executor string,
	contract string,
	msg wasmtypes.RawContractMessage,
	frequency int64,
	funds sdk.Coins,
) *MsgCreateHook {
	return &MsgCreateHook{
		Authority: authority,
		Executor:  executor,
		Contract:  contract,
		Msg:       msg,
		Frequency: frequency,
		Funds:     funds,
	}
}

func (msg MsgCreateHook) Route() string {
	return RouterKey
}

func (msg MsgCreateHook) Type() string {
	return TypeMsgCreateHook
}

func (msg MsgCreateHook) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Contract); err != nil {
		return errors.Wrap(err, "contract")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Executor); err != nil {
		return errors.Wrap(err, "executor")
	}
	if !msg.Funds.IsValid() {
		return sdkerrors.ErrInvalidCoins
	}
	if err := msg.Msg.ValidateBasic(); err != nil {
		return errors.Wrap(err, "payload msg")
	}
	return nil
}

func (msg MsgCreateHook) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{authority}
}

func (msg MsgCreateHook) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func NewMsgUpdateHook(
	authority string,
	id uint64,
	executor string,
	contract string,
	msg wasmtypes.RawContractMessage,
	frequency int64,
	funds sdk.Coins,
) *MsgUpdateHook {
	return &MsgUpdateHook{
		Authority: authority,
		Id:        id,
		Executor:  executor,
		Contract:  contract,
		Msg:       msg,
		Frequency: frequency,
		Funds:     funds,
	}
}

func (msg MsgUpdateHook) Route() string {
	return RouterKey
}

func (msg MsgUpdateHook) Type() string {
	return TypeMsgUpdateHook
}

func (msg MsgUpdateHook) ValidateBasic() error {
	if msg.Id == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "ID is required")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Contract); err != nil {
		return errors.Wrap(err, "contract")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Executor); err != nil {
		return errors.Wrap(err, "executor")
	}
	if !msg.Funds.IsValid() {
		return sdkerrors.ErrInvalidCoins
	}
	if err := msg.Msg.ValidateBasic(); err != nil {
		return errors.Wrap(err, "payload msg")
	}
	return nil
}

func (msg MsgUpdateHook) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{authority}
}

func (msg MsgUpdateHook) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

func NewMsgDeleteHook(
	authority string,
	id uint64,
) *MsgDeleteHook {
	return &MsgDeleteHook{
		Authority: authority,
		Id:        id,
	}
}

func (msg MsgDeleteHook) Route() string {
	return RouterKey
}

func (msg MsgDeleteHook) Type() string {
	return TypeMsgDeleteHook
}

func (msg MsgDeleteHook) ValidateBasic() error {
	if msg.Id == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "ID is required")
	}
	return nil
}

func (msg MsgDeleteHook) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{authority}
}

func (msg MsgDeleteHook) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}
