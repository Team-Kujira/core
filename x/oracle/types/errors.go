package types

import (
	"fmt"

	"github.com/cometbft/cometbft/crypto/tmhash"

	"cosmossdk.io/errors"
)

// Oracle Errors
var (
	ErrInvalidExchangeRate           = errors.Register(ModuleName, 1, "invalid exchange rate")
	ErrNoPrevote                     = errors.Register(ModuleName, 2, "no prevote")
	ErrNoVote                        = errors.Register(ModuleName, 3, "no vote")
	ErrNoVotingPermission            = errors.Register(ModuleName, 4, "unauthorized voter")
	ErrInvalidHash                   = errors.Register(ModuleName, 5, "invalid hash")
	ErrInvalidHashLength             = errors.Register(ModuleName, 6, fmt.Sprintf("invalid hash length; should equal %d", tmhash.TruncatedSize))
	ErrVerificationFailed            = errors.Register(ModuleName, 7, "hash verification failed")
	ErrRevealPeriodMissMatch         = errors.Register(ModuleName, 8, "reveal period of submitted vote do not match with registered prevote")
	ErrInvalidSaltLength             = errors.Register(ModuleName, 9, "invalid salt length; must be 64")
	ErrInvalidSaltFormat             = errors.Register(ModuleName, 10, "invalid salt format")
	ErrNoAggregatePrevote            = errors.Register(ModuleName, 11, "no aggregate prevote")
	ErrNoAggregateVote               = errors.Register(ModuleName, 12, "no aggregate vote")
	ErrUnknownDenom                  = errors.Register(ModuleName, 13, "unknown denom")
	ErrBallotNotSorted               = errors.Register(ModuleName, 14, "ballot not sorted")
	ErrSetParams                     = errors.Register(ModuleName, 15, "could not set params")
	ErrUnknownHistoricalExchangeRate = errors.Register(ModuleName, 16, "unknown historical exchange rate")
)
