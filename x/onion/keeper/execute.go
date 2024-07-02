package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/group/errors"
)

func (k Keeper) ExecuteTxMsgs(ctx sdk.Context, tx sdk.Tx) ([]sdk.Result, error) {
	msgs := tx.GetMsgs()
	results := make([]sdk.Result, len(msgs))
	for i, msg := range msgs {
		handler := k.router.Handler(msg)
		if handler == nil {
			return nil, errorsmod.Wrapf(errors.ErrInvalid, "no message handler found for %q", sdk.MsgTypeURL(msg))
		}
		r, err := handler(ctx, msg)
		if err != nil {
			return nil, errorsmod.Wrapf(err, "message %s at position %d", sdk.MsgTypeURL(msg), i)
		}
		// Handler should always return non-nil sdk.Result.
		if r == nil {
			return nil, fmt.Errorf("got nil sdk.Result for message %q at position %d", msg, i)
		}

		results[i] = *r
	}
	return results, nil
}
