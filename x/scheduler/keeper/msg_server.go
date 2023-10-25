package keeper

import (
	"context"
	"fmt"

	"github.com/Team-Kujira/core/x/scheduler/types"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the scheduler MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (ms msgServer) CreateHook(
	goCtx context.Context,
	msg *types.MsgCreateHook,
) (*types.MsgCreateHookResponse, error) {
	if ms.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner,
			"invalid authority; expected %s, got %s",
			ms.authority,
			msg.Authority,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	hook := types.Hook{
		Executor:  msg.Executor,
		Contract:  msg.Contract,
		Msg:       msg.Msg,
		Frequency: msg.Frequency,
		Funds:     msg.Funds,
	}

	ms.AppendHook(ctx, hook)

	return &types.MsgCreateHookResponse{}, nil
}

func (ms msgServer) UpdateHook(
	goCtx context.Context,
	msg *types.MsgUpdateHook,
) (*types.MsgUpdateHookResponse, error) {
	if ms.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner,
			"invalid authority; expected %s, got %s",
			ms.authority,
			msg.Authority,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	hook := types.Hook{
		Executor:  msg.Executor,
		Contract:  msg.Contract,
		Msg:       msg.Msg,
		Frequency: msg.Frequency,
		Funds:     msg.Funds,
	}

	// Checks that the element exists
	_, found := ms.GetHook(ctx, msg.Id)
	if !found {
		return nil, errors.Wrap(
			sdkerrors.ErrKeyNotFound,
			fmt.Sprintf("key %d doesn't exist", msg.Id),
		)
	}

	ms.SetHook(ctx, hook)

	return &types.MsgUpdateHookResponse{}, nil
}

func (ms msgServer) DeleteHook(
	goCtx context.Context,
	msg *types.MsgDeleteHook,
) (*types.MsgDeleteHookResponse, error) {
	if ms.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner,
			"invalid authority; expected %s, got %s",
			ms.authority,
			msg.Authority,
		)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	_, found := ms.GetHook(ctx, msg.Id)
	if !found {
		return nil, errors.Wrap(
			sdkerrors.ErrKeyNotFound,
			fmt.Sprintf("key %d doesn't exist", msg.Id),
		)
	}

	ms.RemoveHook(ctx, msg.Id)

	return &types.MsgDeleteHookResponse{}, nil
}
