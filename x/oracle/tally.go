package oracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Team-Kujira/core/x/oracle/types"
)

// Tally calculates the median and returns it. Sets the set of voters to be rewarded, i.e. voted within
// a reasonable spread from the weighted median to the store
// CONTRACT: pb must be sorted
func Tally(_ sdk.Context,
	pb types.ExchangeRateBallot,
	rewardBand sdk.Dec,
	validatorClaimMap map[string]types.Claim,
	missMap map[string]sdk.ValAddress,
) (sdk.Dec, error) {
	weightedMedian, err := pb.WeightedMedian()
	if err != nil {
		return sdk.ZeroDec(), err
	}

	standardDeviation, err := pb.StandardDeviation()
	if err != nil {
		return sdk.ZeroDec(), err
	}

	rewardSpread := weightedMedian.Mul(rewardBand.QuoInt64(2))
	rewardSpread = sdk.MaxDec(rewardSpread, standardDeviation)

	for _, vote := range pb {
		key := vote.Voter.String()
		claim := validatorClaimMap[key]
		// Filter ballot winners & abstain voters
		if (vote.ExchangeRate.GTE(weightedMedian.Sub(rewardSpread)) &&
			vote.ExchangeRate.LTE(weightedMedian.Add(rewardSpread))) ||
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
