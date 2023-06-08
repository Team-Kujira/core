package types

// DONTCOVER
import errorsmod "cosmossdk.io/errors"

// x/scheduler module sentinel errors
var (
	ErrSample = errorsmod.Register(ModuleName, 1100, "sample error")
)
