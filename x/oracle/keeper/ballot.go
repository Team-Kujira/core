package keeper

import (
	"sort"

	"github.com/Team-Kujira/core/x/oracle/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// OrganizeBallotByDenom collects all oracle votes for the period, categorized by the votes' denom parameter
func (k Keeper) OrganizeBallotByDenom(ctx sdk.Context, validatorClaimMap map[string]types.Claim) (votes map[string]types.ExchangeRateBallot) {
	votes = map[string]types.ExchangeRateBallot{}

	// Organize aggregate votes
	aggregateHandler := func(voterAddr sdk.ValAddress, vote types.AggregateExchangeRateVote) (stop bool) {
		// organize ballot only for the active validators
		claim, ok := validatorClaimMap[vote.Voter]

		if ok {
			power := claim.Power
			for _, tuple := range vote.ExchangeRateTuples {
				tmpPower := power
				if !tuple.ExchangeRate.IsPositive() {
					// Make the power of abstain vote zero
					tmpPower = 0
				}

				votes[tuple.Denom] = append(votes[tuple.Denom],
					types.NewVoteForTally(
						tuple.ExchangeRate,
						tuple.Denom,
						voterAddr,
						tmpPower,
					),
				)
			}

		}

		return false
	}

	k.IterateAggregateExchangeRateVotes(ctx, aggregateHandler)

	// sort created ballot
	for denom, ballot := range votes {
		sort.Sort(ballot)
		votes[denom] = ballot
	}

	return
}

// ClearBallots clears all tallied prevotes and votes from the store
func (k Keeper) ClearBallots(ctx sdk.Context, votePeriod uint64) {
	// Clear all aggregate prevotes
	k.IterateAggregateExchangeRatePrevotes(ctx, func(voterAddr sdk.ValAddress, aggregatePrevote types.AggregateExchangeRatePrevote) (stop bool) {
		if ctx.BlockHeight() > int64(aggregatePrevote.SubmitBlock+votePeriod) {
			k.DeleteAggregateExchangeRatePrevote(ctx, voterAddr)
		}

		return false
	})

	// Clear all aggregate votes
	k.IterateAggregateExchangeRateVotes(ctx, func(voterAddr sdk.ValAddress, aggregateVote types.AggregateExchangeRateVote) (stop bool) {
		k.DeleteAggregateExchangeRateVote(ctx, voterAddr)
		return false
	})
}
