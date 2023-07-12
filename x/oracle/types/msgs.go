package types

import (
	"github.com/cometbft/cometbft/crypto/tmhash"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = &MsgDelegateFeedConsent{}
	_ sdk.Msg = &MsgAggregateExchangeRatePrevote{}
	_ sdk.Msg = &MsgAggregateExchangeRateVote{}
	_ sdk.Msg = &MsgAddPrice{}
	_ sdk.Msg = &MsgRemovePrice{}
	_ sdk.Msg = &MsgUpdateParams{}
)

// oracle message types
const (
	TypeMsgDelegateFeedConsent          = "delegate_feeder"
	TypeMsgAggregateExchangeRatePrevote = "aggregate_exchange_rate_prevote"
	TypeMsgAggregateExchangeRateVote    = "aggregate_exchange_rate_vote"
	TypeMsgAddPrice                     = "add_price"
	TypeMsgRemovePrice                  = "remove_price"
	TypeMsgUpdateParams                 = "update_params"
)

//-------------------------------------------------
//-------------------------------------------------

// NewMsgAggregateExchangeRatePrevote returns MsgAggregateExchangeRatePrevote instance
func NewMsgAggregateExchangeRatePrevote(hash AggregateVoteHash, feeder sdk.AccAddress, validator sdk.ValAddress) *MsgAggregateExchangeRatePrevote {
	return &MsgAggregateExchangeRatePrevote{
		Hash:      hash.String(),
		Feeder:    feeder.String(),
		Validator: validator.String(),
	}
}

// Route implements sdk.Msg
func (msg MsgAggregateExchangeRatePrevote) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgAggregateExchangeRatePrevote) Type() string { return TypeMsgAggregateExchangeRatePrevote }

// GetSignBytes implements sdk.Msg
func (msg MsgAggregateExchangeRatePrevote) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements sdk.Msg
func (msg MsgAggregateExchangeRatePrevote) GetSigners() []sdk.AccAddress {
	feeder, err := sdk.AccAddressFromBech32(msg.Feeder)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{feeder}
}

// ValidateBasic Implements sdk.Msg
func (msg MsgAggregateExchangeRatePrevote) ValidateBasic() error {
	_, err := AggregateVoteHashFromHexString(msg.Hash)
	if err != nil {
		return errors.Wrapf(ErrInvalidHash, "Invalid vote hash (%s)", err)
	}

	// HEX encoding doubles the hash length
	if len(msg.Hash) != tmhash.TruncatedSize*2 {
		return ErrInvalidHashLength
	}

	_, err = sdk.AccAddressFromBech32(msg.Feeder)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid feeder address (%s)", err)
	}

	_, err = sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid operator address (%s)", err)
	}

	return nil
}

// NewMsgAggregateExchangeRateVote returns MsgAggregateExchangeRateVote instance
func NewMsgAggregateExchangeRateVote(salt string, exchangeRates string, feeder sdk.AccAddress, validator sdk.ValAddress) *MsgAggregateExchangeRateVote {
	return &MsgAggregateExchangeRateVote{
		Salt:          salt,
		ExchangeRates: exchangeRates,
		Feeder:        feeder.String(),
		Validator:     validator.String(),
	}
}

// Route implements sdk.Msg
func (msg MsgAggregateExchangeRateVote) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgAggregateExchangeRateVote) Type() string { return TypeMsgAggregateExchangeRateVote }

// GetSignBytes implements sdk.Msg
func (msg MsgAggregateExchangeRateVote) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements sdk.Msg
func (msg MsgAggregateExchangeRateVote) GetSigners() []sdk.AccAddress {
	feeder, err := sdk.AccAddressFromBech32(msg.Feeder)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{feeder}
}

// ValidateBasic implements sdk.Msg
func (msg MsgAggregateExchangeRateVote) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Feeder)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid feeder address (%s)", err)
	}

	_, err = sdk.ValAddressFromBech32(msg.Validator)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid operator address (%s)", err)
	}

	if l := len(msg.ExchangeRates); l == 0 {
		return errors.Wrap(sdkerrors.ErrUnknownRequest, "must provide at least one oracle exchange rate")
	} else if l > 4096 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "exchange rates string can not exceed 4096 characters")
	}

	exchangeRates, err := ParseExchangeRateTuples(msg.ExchangeRates)
	if err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, "failed to parse exchange rates string cause: "+err.Error())
	}

	for _, exchangeRate := range exchangeRates {
		// Check overflow bit length
		if exchangeRate.ExchangeRate.BigInt().BitLen() > 255+sdk.DecimalPrecisionBits {
			return errors.Wrap(ErrInvalidExchangeRate, "overflow")
		}
	}

	if len(msg.Salt) != 64 {
		return ErrInvalidSaltLength
	}
	_, err = AggregateVoteHashFromHexString(msg.Salt)
	if err != nil {
		return errors.Wrap(ErrInvalidSaltFormat, "salt must be a valid hex string")
	}

	return nil
}

// NewMsgDelegateFeedConsent creates a MsgDelegateFeedConsent instance
func NewMsgDelegateFeedConsent(operatorAddress sdk.ValAddress, feederAddress sdk.AccAddress) *MsgDelegateFeedConsent {
	return &MsgDelegateFeedConsent{
		Operator: operatorAddress.String(),
		Delegate: feederAddress.String(),
	}
}

// Route implements sdk.Msg
func (msg MsgDelegateFeedConsent) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgDelegateFeedConsent) Type() string { return TypeMsgDelegateFeedConsent }

// GetSignBytes implements sdk.Msg
func (msg MsgDelegateFeedConsent) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements sdk.Msg
func (msg MsgDelegateFeedConsent) GetSigners() []sdk.AccAddress {
	operator, err := sdk.ValAddressFromBech32(msg.Operator)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{sdk.AccAddress(operator)}
}

// ValidateBasic implements sdk.Msg
func (msg MsgDelegateFeedConsent) ValidateBasic() error {
	_, err := sdk.ValAddressFromBech32(msg.Operator)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid operator address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Delegate)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid delegate address (%s)", err)
	}

	return nil
}

// NewMsgAddPrice creates a MsgAddPrice instance
func NewMsgAddPrice(symbol string) *MsgAddPrice {
	return &MsgAddPrice{
		Symbol: symbol,
	}
}

// Route implements sdk.Msg
func (msg MsgAddPrice) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgAddPrice) Type() string { return TypeMsgAddPrice }

// GetSignBytes implements sdk.Msg
func (msg MsgAddPrice) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements sdk.Msg
func (msg MsgAddPrice) GetSigners() []sdk.AccAddress {
	operator, err := sdk.ValAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{sdk.AccAddress(operator)}
}

// ValidateBasic implements sdk.Msg
func (msg MsgAddPrice) ValidateBasic() error {
	_, err := sdk.ValAddressFromBech32(msg.Authority)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid authority address (%s)", err)
	}

	return nil
}

// NewMsgRemovePrice creates a MsgRemovePrice instance
func NewMsgRemovePrice(symbol string) *MsgRemovePrice {
	return &MsgRemovePrice{
		Symbol: symbol,
	}
}

// Route implements sdk.Msg
func (msg MsgRemovePrice) Route() string { return RouterKey }

// Type implements sdk.Msg
func (msg MsgRemovePrice) Type() string { return TypeMsgRemovePrice }

// GetSignBytes implements sdk.Msg
func (msg MsgRemovePrice) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements sdk.Msg
func (msg MsgRemovePrice) GetSigners() []sdk.AccAddress {
	operator, err := sdk.ValAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{sdk.AccAddress(operator)}
}

// ValidateBasic implements sdk.Msg
func (msg MsgRemovePrice) ValidateBasic() error {
	_, err := sdk.ValAddressFromBech32(msg.Authority)
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

// GetSignBytes implements sdk.Msg
func (msg MsgUpdateParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements sdk.Msg
func (msg MsgUpdateParams) GetSigners() []sdk.AccAddress {
	operator, err := sdk.ValAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{sdk.AccAddress(operator)}
}

// ValidateBasic implements sdk.Msg
func (msg MsgUpdateParams) ValidateBasic() error {
	_, err := sdk.ValAddressFromBech32(msg.Authority)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid authority address (%s)", err)
	}

	return nil
}
