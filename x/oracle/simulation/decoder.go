package simulation

import (
	"bytes"
	"fmt"

	gogotypes "github.com/cosmos/gogoproto/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/Team-Kujira/core/x/oracle/types"
)

// NewDecodeStore returns a decoder function closure that unmarshals the KVPair's
// Value to the corresponding oracle type.
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], types.ExchangeRateKey):
			var exchangeRateA, exchangeRateB sdk.DecProto
			cdc.MustUnmarshal(kvA.Value, &exchangeRateA)
			cdc.MustUnmarshal(kvB.Value, &exchangeRateB)
			return fmt.Sprintf("%v\n%v", exchangeRateA, exchangeRateB)
		case bytes.Equal(kvA.Key[:1], types.FeederDelegationKey):
			return fmt.Sprintf("%v\n%v", sdk.AccAddress(kvA.Value), sdk.AccAddress(kvB.Value))
		case bytes.Equal(kvA.Key[:1], types.MissCounterKey):
			var counterA, counterB gogotypes.UInt64Value
			cdc.MustUnmarshal(kvA.Value, &counterA)
			cdc.MustUnmarshal(kvB.Value, &counterB)
			return fmt.Sprintf("%v\n%v", counterA.Value, counterB.Value)
		default:
			panic(fmt.Sprintf("invalid oracle key prefix %X", kvA.Key[:1]))
		}
	}
}
