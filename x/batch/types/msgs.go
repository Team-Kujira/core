package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// constants
const (
	TypeMsgWithdrawAllDelegatorRewards = "withdraw_all_rewards"
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

