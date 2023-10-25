package types

import (
	"sort"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NOTE: we don't need to implement proto interface on this file
//       these are not used in store or rpc response

// VoteForTally is a convenience wrapper to reduce redundant lookup cost
type VoteForTally struct {
	Denom        string
	ExchangeRate math.LegacyDec
	Voter        sdk.ValAddress
	Power        int64
}

// NewVoteForTally returns a new VoteForTally instance
func NewVoteForTally(rate math.LegacyDec, denom string, voter sdk.ValAddress, power int64) VoteForTally {
	return VoteForTally{
		ExchangeRate: rate,
		Denom:        denom,
		Voter:        voter,
		Power:        power,
	}
}

// ExchangeRateBallot is a convenience wrapper around a ExchangeRateVote slice
type ExchangeRateBallot []VoteForTally

// ToMap return organized exchange rate map by validator
func (pb ExchangeRateBallot) ToMap() map[string]math.LegacyDec {
	exchangeRateMap := make(map[string]math.LegacyDec)
	for _, vote := range pb {
		if vote.ExchangeRate.IsPositive() {
			exchangeRateMap[string(vote.Voter)] = vote.ExchangeRate
		}
	}

	return exchangeRateMap
}

// Power returns the total amount of voting power in the ballot
func (pb ExchangeRateBallot) Power() int64 {
	totalPower := int64(0)
	for _, vote := range pb {
		totalPower += vote.Power
	}

	return totalPower
}

// WeightedMedian returns the median weighted by the power of the ExchangeRateVote.
// CONTRACT: ballot must be sorted
func (pb ExchangeRateBallot) WeightedMedian() (math.LegacyDec, error) {
	if !sort.IsSorted(pb) {
		return math.LegacyZeroDec(), ErrBallotNotSorted
	}

	totalPower := pb.Power()
	if pb.Len() > 0 {
		pivot := int64(0)
		for _, v := range pb {
			votePower := v.Power

			pivot += votePower
			if pivot >= (totalPower / 2) {
				return v.ExchangeRate, nil
			}
		}
	}
	return math.LegacyZeroDec(), nil
}

// StandardDeviation returns the standard deviation by the power of the ExchangeRateVote.
func (pb ExchangeRateBallot) StandardDeviation() (math.LegacyDec, error) {
	if len(pb) == 0 {
		return math.LegacyZeroDec(), nil
	}

	median, err := pb.WeightedMedian()
	if err != nil {
		return math.LegacyZeroDec(), err
	}

	sum := math.LegacyZeroDec()
	ballotLength := int64(len(pb))
	for _, v := range pb {
		func() {
			defer func() {
				if e := recover(); e != nil {
					ballotLength--
				}
			}()
			deviation := v.ExchangeRate.Sub(median)
			sum = sum.Add(deviation.Mul(deviation))
		}()
	}

	variance := sum.QuoInt64(ballotLength)

	standardDeviation, err := variance.ApproxSqrt()
	if err != nil {
		return math.LegacyZeroDec(), err
	}

	return standardDeviation, nil
}

// Len implements sort.Interface
func (pb ExchangeRateBallot) Len() int {
	return len(pb)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (pb ExchangeRateBallot) Less(i, j int) bool {
	return pb[i].ExchangeRate.LT(pb[j].ExchangeRate)
}

// Swap implements sort.Interface.
func (pb ExchangeRateBallot) Swap(i, j int) {
	pb[i], pb[j] = pb[j], pb[i]
}

// Claim is an interface that directs its rewards to an attached bank account.
type Claim struct {
	Power     int64
	Weight    int64
	WinCount  int64
	Recipient sdk.ValAddress
}

// NewClaim generates a Claim instance.
func NewClaim(power, weight, winCount int64, recipient sdk.ValAddress) Claim {
	return Claim{
		Power:     power,
		Weight:    weight,
		WinCount:  winCount,
		Recipient: recipient,
	}
}
