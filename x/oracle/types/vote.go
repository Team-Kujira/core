package types

import (
	"fmt"
	"strings"

	"cosmossdk.io/math"
	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewExchangeRateTuple creates a ExchangeRateTuple instance
func NewExchangeRateTuple(denom string, exchangeRate math.LegacyDec) ExchangeRateTuple {
	return ExchangeRateTuple{
		denom,
		exchangeRate,
	}
}

// String implement stringify
func (v ExchangeRateTuple) String() string {
	out, _ := yaml.Marshal(v)
	return string(out)
}

// ExchangeRateTuples - array of ExchangeRateTuple
type ExchangeRateTuples []ExchangeRateTuple

// String implements fmt.Stringer interface
func (tuples ExchangeRateTuples) String() string {
	out, _ := yaml.Marshal(tuples)
	return string(out)
}

// ParseExchangeRateTuples ExchangeRateTuple parser
func ParseExchangeRateTuples(tuplesStr string) (ExchangeRateTuples, error) {
	tuplesStr = strings.TrimSpace(tuplesStr)
	if len(tuplesStr) == 0 {
		return nil, nil
	}

	tupleStrs := strings.Split(tuplesStr, ",")
	tuples := make(ExchangeRateTuples, len(tupleStrs))
	duplicateCheckMap := make(map[string]bool)
	for i, tupleStr := range tupleStrs {
		decCoin, err := sdk.ParseDecCoin(tupleStr)
		if err != nil {
			return nil, err
		}

		tuples[i] = ExchangeRateTuple{
			Denom:        decCoin.Denom,
			ExchangeRate: decCoin.Amount,
		}

		if _, ok := duplicateCheckMap[decCoin.Denom]; ok {
			return nil, fmt.Errorf("duplicated denom %s", decCoin.Denom)
		}

		duplicateCheckMap[decCoin.Denom] = true
	}

	return tuples, nil
}
