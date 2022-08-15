package keeper

import (
	"fmt"

	"github.com/Team-Kujira/core/x/scheduler/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// NewSchedulerProposalHandler defines the 02-client proposal handler
func NewSchedulerProposalHandler(k Keeper) govtypes.Handler {
	return func(ctx sdk.Context, content govtypes.Content) error {
		switch c := content.(type) {
		case *types.CreateHookProposal:

			var hook = types.Hook{
				Executor:  c.Executor,
				Contract:  c.Contract,
				Msg:       c.Msg,
				Frequency: c.Frequency,
				Funds:     c.Funds,
			}

			k.AppendHook(ctx, hook)

			return nil

		case *types.UpdateHookProposal:

			var hook = types.Hook{
				Executor:  c.Executor,
				Id:        c.Id,
				Contract:  c.Contract,
				Msg:       c.Msg,
				Frequency: c.Frequency,
				Funds:     c.Funds,
			}

			// Checks that the element exists
			_, found := k.GetHook(ctx, c.Id)
			if !found {
				return sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", c.Id))
			}

			k.SetHook(ctx, hook)

			return nil

		case *types.DeleteHookProposal:

			// Checks that the element exists
			_, found := k.GetHook(ctx, c.Id)
			if !found {
				return sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", c.Id))
			}

			k.RemoveHook(ctx, c.Id)

			return nil
		default:
			return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized scheduler proposal content type: %T", c)
		}
	}
}
