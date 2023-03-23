package types

import (
	errorsmod "cosmossdk.io/errors"
)

// DONTCOVER

// x/scheduler module sentinel errors
var (
	ErrSample = errorsmod.Register(ModuleName, 1100, "sample error")
)
