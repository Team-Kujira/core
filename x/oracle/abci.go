package oracle

import (
	"time"

	"github.com/Team-Kujira/core/x/oracle/keeper"
	"github.com/Team-Kujira/core/x/oracle/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker is called at the end of every block
func EndBlocker(ctx sdk.Context, k keeper.Keeper) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)
	params := k.GetParams(ctx)
	if IsPeriodLastBlock(ctx, params.VotePeriod) {
		// Build claim map over all validators in active set
		validatorClaimMap := make(map[string]types.Claim)

		maxValidators := k.StakingKeeper.MaxValidators(ctx)
		iterator := k.StakingKeeper.ValidatorsPowerStoreIterator(ctx)
		defer iterator.Close()

		powerReduction := k.StakingKeeper.PowerReduction(ctx)

		i := 0
		for ; iterator.Valid() && i < int(maxValidators); iterator.Next() {
			validator := k.StakingKeeper.Validator(ctx, iterator.Value())

			// Exclude not bonded validator
			if validator.IsBonded() {
				valAddr := validator.GetOperator()
				validatorClaimMap[valAddr.String()] = types.NewClaim(validator.GetConsensusPower(powerReduction), 0, 0, valAddr)
				i++
			}
		}

		// voteTargets defines the symbol (ticker) denoms that we require votes on
		var voteTargets []string
		for _, v := range params.Whitelist {
			voteTargets = append(voteTargets, v.Name)
		}

		// Clear all exchange rates
		k.IterateExchangeRates(ctx, func(denom string, _ sdk.Dec) (stop bool) {
			k.DeleteExchangeRate(ctx, denom)
			return false
		})

		// Organize votes to ballot by denom
		voteMap := k.OrganizeBallotByDenom(ctx, validatorClaimMap)

		// Keep track, if a voter submitted a price deviating too much
		missMap := map[string]sdk.ValAddress{}

		// Iterate through ballots and update exchange rates; drop if not enough votes have been achieved.
		for denom, ballot := range voteMap {
			totalBondedPower := sdk.TokensToConsensusPower(k.StakingKeeper.TotalBondedTokens(ctx), k.StakingKeeper.PowerReduction(ctx))
			voteThreshold := k.VoteThreshold(ctx)
			thresholdVotes := voteThreshold.MulInt64(totalBondedPower).RoundInt()
			ballotPower := sdk.NewInt(ballot.Power())

			if !ballotPower.IsZero() && ballotPower.GTE(thresholdVotes) {
				exchangeRate, err := Tally(
					ctx, ballot, params.RewardBand, validatorClaimMap, missMap,
				)
				if err != nil {
					return err
				}

				// Set the exchange rate, emit ABCI event
				k.SetExchangeRateWithEvent(ctx, denom, exchangeRate)
			}
		}

		//---------------------------
		// Do miss counting & slashing
		denomMap := map[string]map[string]struct{}{}

		for _, denom := range voteTargets {
			denomMap[denom] = map[string]struct{}{}
		}

		for denom, votes := range voteMap {
			for _, vote := range votes {
				// ignore denoms, not requested in voteTargets
				_, ok := denomMap[denom]
				if !ok {
					continue
				}

				denomMap[denom][vote.Voter.String()] = struct{}{}
			}
		}

		// Check if each validator is missing a required denom price
		for _, claim := range validatorClaimMap {
			for _, denom := range voteTargets {
				_, ok := denomMap[denom][claim.Recipient.String()]
				if !ok {
					missMap[claim.Recipient.String()] = claim.Recipient
					break
				}
			}
		}

		for _, valAddr := range missMap {
			k.SetMissCounter(ctx, valAddr, k.GetMissCounter(ctx, valAddr)+1)
		}

		// // Distribute rewards to ballot winners
		// k.RewardBallotWinners(
		// 	ctx,
		// 	(int64)(params.VotePeriod),
		// 	(int64)(params.RewardDistributionWindow),
		// 	voteTargets,
		// 	validatorClaimMap,
		// )

		// Clear the ballot
		k.ClearBallots(ctx, params.VotePeriod)
	}

	// Do slash who did miss voting over threshold and
	// reset miss counters of all validators at the last block of slash window
	if IsPeriodLastBlock(ctx, params.SlashWindow) {
		k.SlashAndResetMissCounters(ctx)
	}

	return nil
}

func IsPeriodLastBlock(ctx sdk.Context, blocksPerPeriod uint64) bool {
	return (uint64(ctx.BlockHeight())+1)%blocksPerPeriod == 0
}
