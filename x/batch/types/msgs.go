package types

import (
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// constants
const (
	TypeMsgWithdrawAllDelegatorRewards = "withdraw_all_rewards"
	TypeMsgBatchResetDelegation        = "batch_reset_delegation"
)

var _ sdk.Msg = &MsgWithdrawAllDelegatorRewards{}

// NewMsgWithdrawAllDelegatorRewards creates a msg to create a new distrib
func NewMsgWithdrawAllDelegatorRewards(delegator sdk.AccAddress) *MsgWithdrawAllDelegatorRewards {
	return &MsgWithdrawAllDelegatorRewards{
		DelegatorAddress: delegator.String(),
	}
}

func (m MsgWithdrawAllDelegatorRewards) Route() string { return RouterKey }
func (m MsgWithdrawAllDelegatorRewards) Type() string  { return TypeMsgWithdrawAllDelegatorRewards }
func (m MsgWithdrawAllDelegatorRewards) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid delegator address (%s)", err)
	}
	return nil
}

func (m MsgWithdrawAllDelegatorRewards) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgWithdrawAllDelegatorRewards) GetSigners() []sdk.AccAddress {
	sender, _ := sdk.AccAddressFromBech32(m.DelegatorAddress)
	return []sdk.AccAddress{sender}
}

var _ sdk.Msg = &MsgBatchResetDelegation{}

// NewMsgWithdrawAllDelegatorRewards creates a msg to create a new distrib
func NewMsgBatchResetDelegation(delegator sdk.AccAddress, validators []string, amounts []math.Int) *MsgBatchResetDelegation {
	return &MsgBatchResetDelegation{
		DelegatorAddress: delegator.String(),
		Validators:       validators,
		Amounts:          amounts,
	}
}

func (m MsgBatchResetDelegation) Route() string { return RouterKey }
func (m MsgBatchResetDelegation) Type() string  { return TypeMsgBatchResetDelegation }
func (m MsgBatchResetDelegation) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
	if err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid delegator address (%s)", err)
	}

	for _, valStr := range m.Validators {
		_, err := sdk.ValAddressFromBech32(valStr)
		if err != nil {
			return err
		}
	}

	for _, amount := range m.Amounts {
		if amount.IsNegative() {
			return ErrNegativeDelegationAmount
		}
	}

	if len(m.Validators) != len(m.Amounts) {
		return ErrValidatorsAndAmountsMismatch
	}

	return nil
}

func (m MsgBatchResetDelegation) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}

func (m MsgBatchResetDelegation) GetSigners() []sdk.AccAddress {
	sender, _ := sdk.AccAddressFromBech32(m.DelegatorAddress)
	return []sdk.AccAddress{sender}
}
