package types

import (
	errorsmod "cosmossdk.io/errors"
)

// DONTCOVER

// x/denom module sentinel errors
var (
	ErrDenomExists              = errorsmod.Register(ModuleName, 2, "denom already exists")
	ErrUnauthorized             = errorsmod.Register(ModuleName, 3, "unauthorized account")
	ErrInvalidDenom             = errorsmod.Register(ModuleName, 4, "invalid denom")
	ErrInvalidCreator           = errorsmod.Register(ModuleName, 5, "invalid creator")
	ErrInvalidAuthorityMetadata = errorsmod.Register(ModuleName, 6, "invalid authority metadata")
	ErrInvalidGenesis           = errorsmod.Register(ModuleName, 7, "invalid genesis")
)
