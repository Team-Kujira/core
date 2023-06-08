package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

// Oracle Errors
var (
	ErrInvalidExchangeRate   = errorsmod.Register(ModuleName, 1, "invalid exchange rate")
	ErrNoPrevote             = errorsmod.Register(ModuleName, 2, "no prevote")
	ErrNoVote                = errorsmod.Register(ModuleName, 3, "no vote")
	ErrNoVotingPermission    = errorsmod.Register(ModuleName, 4, "unauthorized voter")
	ErrInvalidHash           = errorsmod.Register(ModuleName, 5, "invalid hash")
	ErrInvalidHashLength     = errorsmod.Register(ModuleName, 6, fmt.Sprintf("invalid hash length; should equal %d", tmhash.TruncatedSize))
	ErrVerificationFailed    = errorsmod.Register(ModuleName, 7, "hash verification failed")
	ErrRevealPeriodMissMatch = errorsmod.Register(ModuleName, 8, "reveal period of submitted vote do not match with registered prevote")
	ErrInvalidSaltLength     = errorsmod.Register(ModuleName, 9, "invalid salt length; must be 64")
	ErrInvalidSaltFormat     = errorsmod.Register(ModuleName, 10, "invalid salt format")
	ErrNoAggregatePrevote    = errorsmod.Register(ModuleName, 11, "no aggregate prevote")
	ErrNoAggregateVote       = errorsmod.Register(ModuleName, 12, "no aggregate vote")
	ErrUnknownDenom          = errorsmod.Register(ModuleName, 13, "unknown denom")
	ErrBallotNotSorted       = errorsmod.Register(ModuleName, 14, "ballot not sorted")
)
