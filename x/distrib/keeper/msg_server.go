package keeper

import (
	"context"

	"github.com/armon/go-metrics"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	"github.com/Team-Kujira/core/x/distrib/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the distrib MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (k msgServer) WithdrawAllDelegatorRewards(goCtx context.Context, msg *types.MsgWithdrawAllDelegatorRewards) (*types.MsgWithdrawAllDelegatorRewardsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return nil, err
	}
	amount, err := k.WithdrawAllDelegationRewards(ctx, delAddr)
	if err != nil {
		return nil, err
	}
	defer func() {
		for _, a := range amount {
			if a.Amount.IsInt64() {
				telemetry.SetGaugeWithLabels(
					[]string{"tx", "msg", "withdraw_reward"},
					float32(a.Amount.Int64()),
					[]metrics.Label{telemetry.NewLabel("denom", a.Denom)},
				)
			}
		}
	}()
	return &types.MsgWithdrawAllDelegatorRewardsResponse{Amount: amount}, nil
}