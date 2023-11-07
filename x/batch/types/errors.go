package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/batch module sentinel errors
var (
	ErrValidatorsAndAmountsMismatch = errorsmod.Register(ModuleName, 1, "validators and amounts length mismatch")
	ErrInvalidAmount                = errorsmod.Register(ModuleName, 2, "invalid amount")
	ErrNegativeDelegationAmount     = errorsmod.Register(ModuleName, 3, "negative delegation amount")
)
