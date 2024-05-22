package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
)

const IbcAcknowledgementErrorType = "ibc-acknowledgement-error"

// NewSuccessAckRepresentingAnError creates a new success acknowledgement that represents an error.
// This is useful for notifying the sender that an error has occurred in a way that does not allow
// the received tokens to be reverted (which means they shouldn't be released by the sender's ics20 escrow)
func NewSuccessAckRepresentingAnError(ctx sdk.Context, err error, errorContent []byte, errorContexts ...string) channeltypes.Acknowledgement {
	EmitIBCErrorEvents(ctx, err, errorContexts)

	return channeltypes.NewResultAcknowledgement(errorContent)
}

// EmitIBCErrorEvents Emit and Log errors
func EmitIBCErrorEvents(ctx sdk.Context, err error, errorContexts []string) {
	attributes := make([]sdk.Attribute, len(errorContexts)+1)
	attributes[0] = sdk.NewAttribute("error", err.Error())
	for i, s := range errorContexts {
		attributes[i+1] = sdk.NewAttribute("error-context", s)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			IbcAcknowledgementErrorType,
			attributes...,
		),
	})
}
