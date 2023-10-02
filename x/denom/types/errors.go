package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/denom module sentinel errors
var (
	ErrDenomExists              = errors.Register(ModuleName, 2, "denom already exists")
	ErrUnauthorized             = errors.Register(ModuleName, 3, "unauthorized account")
	ErrInvalidDenom             = errors.Register(ModuleName, 4, "invalid denom")
	ErrInvalidCreator           = errors.Register(ModuleName, 5, "invalid creator")
	ErrInvalidAuthorityMetadata = errors.Register(ModuleName, 6, "invalid authority metadata")
	ErrInvalidGenesis           = errors.Register(ModuleName, 7, "invalid genesis")
)
