package keeper

import (
	"context"
	"fmt"

	"kujira/x/scheduler/types"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateHook(goCtx context.Context, msg *types.MsgCreateHook) (*types.MsgCreateHookResponse, error) {
	if k.authority != msg.Creator {
		return nil, errors.Wrapf(sdkerrors.ErrorInvalidSigner, "expected %s got %s", k.authority, msg.Creator)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	var hook = types.Hook{
		Creator:   msg.Creator,
		Executor:  msg.Executor,
		Contract:  msg.Contract,
		Msg:       msg.Msg,
		Frequency: msg.Frequency,
	}

	id := k.AppendHook(
		ctx,
		hook,
	)

	return &types.MsgCreateHookResponse{
		Id: id,
	}, nil
}

func (k msgServer) UpdateHook(goCtx context.Context, msg *types.MsgUpdateHook) (*types.MsgUpdateHookResponse, error) {
	if k.authority != msg.Creator {
		return nil, errors.Wrapf(sdkerrors.ErrorInvalidSigner, "expected %s got %s", k.authority, msg.Creator)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	var hook = types.Hook{
		Creator:   msg.Creator,
		Executor:  msg.Executor,
		Id:        msg.Id,
		Contract:  msg.Contract,
		Msg:       msg.Msg,
		Frequency: msg.Frequency,
	}

	// Checks that the element exists
	val, found := k.GetHook(ctx, msg.Id)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetHook(ctx, hook)

	return &types.MsgUpdateHookResponse{}, nil
}

func (k msgServer) DeleteHook(goCtx context.Context, msg *types.MsgDeleteHook) (*types.MsgDeleteHookResponse, error) {
	if k.authority != msg.Creator {
		return nil, errors.Wrapf(sdkerrors.ErrorInvalidSigner, "expected %s got %s", k.authority, msg.Creator)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Checks that the element exists
	val, found := k.GetHook(ctx, msg.Id)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveHook(ctx, msg.Id)

	return &types.MsgDeleteHookResponse{}, nil
}
