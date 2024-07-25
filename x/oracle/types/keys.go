package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName is the name of the oracle module
	ModuleName = "oracle"

	// StoreKey is the string store representation
	StoreKey = ModuleName

	// RouterKey is the msg router key for the oracle module
	RouterKey = ModuleName

	// QuerierRoute is the query router key for the oracle module
	QuerierRoute = ModuleName
)

// Keys for oracle store
// Items are stored with the following key: values
//
// - 0x01<denom_Bytes>: sdk.Dec
//
// - 0x03<valAddress_Bytes>: int64
//
// - 0x06: Params
var (
	// Keys for store prefixes
	ExchangeRateKey = []byte{0x01} // prefix for each key to a rate
	MissCounterKey  = []byte{0x03} // prefix for each key to a miss counter
	ParamsKey       = []byte{0x06}
)

// GetExchangeRateKey - stored by *denom*
func GetExchangeRateKey(denom string) []byte {
	return append(ExchangeRateKey, []byte(denom)...)
}

// GetMissCounterKey - stored by *Validator* address
func GetMissCounterKey(v sdk.ValAddress) []byte {
	return append(MissCounterKey, address.MustLengthPrefix(v)...)
}
