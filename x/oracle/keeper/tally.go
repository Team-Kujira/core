package keeper

import (
	"github.com/hashicorp/go-metrics"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Team-Kujira/core/x/oracle/types"
)

// Tally calculates the median and returns it. Sets the set of voters to be rewarded, i.e. voted within
// a reasonable spread from the weighted median to the store
// CONTRACT: pb must be sorted
func Tally(_ sdk.Context,
	pb types.ExchangeRateBallot,
	maxDeviation math.LegacyDec,
	validatorClaimMap map[string]types.Claim,
	missMap map[string]sdk.ValAddress,
) (math.LegacyDec, error) {
	weightedMedian, err := pb.WeightedMedian()
	if err != nil {
		return math.LegacyZeroDec(), err
	}

	standardDeviation, err := pb.StandardDeviation()
	if err != nil {
		return math.LegacyZeroDec(), err
	}

	spread := weightedMedian.Mul(maxDeviation)
	spread = math.LegacyMaxDec(spread, standardDeviation)

	if pb.Len() == 0 {
		return weightedMedian, nil
	}

	labels := []metrics.Label{
		telemetry.NewLabel("denom", pb[0].Denom),
	}

	telemetry.SetGaugeWithLabels(
		[]string{"oracle", "median"},
		float32(weightedMedian.MustFloat64()),
		labels,
	)

	telemetry.SetGaugeWithLabels(
		[]string{"oracle", "stddev"},
		float32(standardDeviation.MustFloat64()),
		labels,
	)

	telemetry.SetGaugeWithLabels(
		[]string{"oracle", "spread"},
		float32(spread.MustFloat64()),
		labels,
	)

	for _, vote := range pb {
		key := vote.Voter.String()
		claim := validatorClaimMap[key]

		telemetry.SetGaugeWithLabels(
			[]string{"oracle", "price"},
			float32(vote.ExchangeRate.MustFloat64()),
			append(labels, telemetry.NewLabel("validator", key)),
		)

		// Filter ballot winners & abstain voters
		if (vote.ExchangeRate.GTE(weightedMedian.Sub(spread)) &&
			vote.ExchangeRate.LTE(weightedMedian.Add(spread))) ||
			!vote.ExchangeRate.IsPositive() {
			claim := validatorClaimMap[key]
			claim.Weight += vote.Power
			claim.WinCount++
			validatorClaimMap[key] = claim
		} else {
			missMap[claim.Recipient.String()] = claim.Recipient
		}
	}

	return weightedMedian, nil
}
