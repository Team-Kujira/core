package types

import (
	"strings"
)

const (
	// ModuleName defines the module name
	ModuleName = "denom"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_denom"
)

// KeySeparator is used to combine parts of the keys in the store
const KeySeparator = "|"

var (
	DenomAuthorityMetadataKey = "authoritymetadata"
	DenomsPrefixKey           = "denoms"
	CreatorPrefixKey          = "creator"
	AdminPrefixKey            = "admin"
	NoFeeAccountPrefixKey     = "nofeeaccount"
)

// GetDenomPrefixStore returns the store prefix where all the data associated with a specific denom
// is stored
func GetDenomPrefixStore(denom string) []byte {
	return []byte(strings.Join([]string{DenomsPrefixKey, denom, ""}, KeySeparator))
}

// GetCreatorsPrefix returns the store prefix where the list of the denoms created by a specific
// creator are stored
func GetCreatorPrefix(creator string) []byte {
	return []byte(strings.Join([]string{CreatorPrefixKey, creator, ""}, KeySeparator))
}

// GetCreatorsPrefix returns the store prefix where a list of all creator addresses are stored
func GetCreatorsPrefix() []byte {
	return []byte(strings.Join([]string{CreatorPrefixKey, ""}, KeySeparator))
}

// GetNoFeeAccountPrefix returns the store prefix where a list of all no fee spending accounts in denom creation
func GetNoFeeAccountPrefix() []byte {
	return []byte(strings.Join([]string{NoFeeAccountPrefixKey, ""}, KeySeparator))
}
