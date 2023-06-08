package bindings

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
)

type AuthMsg struct {
	CreateVestingAccount *CreateVestingAccount `json:"create_vesting_account,omitempty"`
}

type CreateVestingAccount struct {
	ToAddress string    `json:"to_address"`
	Amount    sdk.Coins `json:"amount"`
	EndTime   sdk.Int   `json:"end_time"`
	Delayed   bool      `json:"delayed,omitempty"`
}

type AuthHandler struct {
	keeper.AccountKeeper
	types.BankKeeper
}

func (s AuthHandler) CreateVestingAccount(ctx sdk.Context, msg *types.MsgCreateVestingAccount) ([]sdk.Event, [][]byte, error) {
	ak := s.AccountKeeper
	bk := s.BankKeeper

	if err := bk.IsSendEnabledCoins(ctx, msg.Amount...); err != nil {
		return nil, nil, err
	}

	from, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return nil, nil, err
	}
	to, err := sdk.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return nil, nil, err
	}

	if bk.BlockedAddr(to) {
		return nil, nil, errorsmod.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive funds", msg.ToAddress)
	}

	if acc := ak.GetAccount(ctx, to); acc != nil {
		return nil, nil, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "account %s already exists", msg.ToAddress)
	}

	baseAccount := authtypes.NewBaseAccountWithAddress(to)
	baseAccount = ak.NewAccount(ctx, baseAccount).(*authtypes.BaseAccount) //nolint:forcetypeassert
	baseVestingAccount := types.NewBaseVestingAccount(baseAccount, msg.Amount.Sort(), msg.EndTime)

	var vestingAccount authtypes.AccountI
	if msg.Delayed {
		vestingAccount = types.NewDelayedVestingAccountRaw(baseVestingAccount)
	} else {
		vestingAccount = types.NewContinuousVestingAccountRaw(baseVestingAccount, ctx.BlockTime().Unix())
	}

	ak.SetAccount(ctx, vestingAccount)

	defer func() {
		telemetry.IncrCounter(1, "new", "account")

		for _, a := range msg.Amount {
			if a.Amount.IsInt64() {
				telemetry.SetGaugeWithLabels(
					[]string{"tx", "msg", "create_vesting_account"},
					float32(a.Amount.Int64()),
					[]metrics.Label{telemetry.NewLabel("denom", a.Denom)},
				)
			}
		}
	}()

	err = bk.SendCoins(ctx, from, to, msg.Amount)
	if err != nil {
		return nil, nil, err
	}

	event := sdk.NewEvent(
		sdk.EventTypeMessage,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
	)

	return sdk.Events{event}, nil, nil
}
