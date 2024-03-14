package keeper

import (
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Team-Kujira/core/x/oracle/types"
)

// RewardBallotWinners implements
// at the end of every VotePeriod, give out a portion of spread fees collected in the oracle reward pool
//
//	to the oracle voters that voted faithfully.
func (k Keeper) RewardBallotWinners(
	ctx sdk.Context,
	votePeriod int64,
	rewardDistributionWindow int64,
	voteTargets []string,
	ballotWinners map[string]types.Claim,
) error {
	rewardDenoms := voteTargets

	// Sum weight of the claims
	ballotPowerSum := int64(0)
	for _, winner := range ballotWinners {
		ballotPowerSum += winner.Weight
	}

	// Exit if the ballot is empty
	if ballotPowerSum == 0 {
		return nil
	}

	// The Reward distributionRatio = votePeriod/rewardDistributionWindow
	distributionRatio := math.LegacyNewDec(votePeriod).QuoInt64(rewardDistributionWindow)

	var periodRewards sdk.DecCoins
	for _, denom := range rewardDenoms {
		rewardPool := k.GetRewardPool(ctx, denom)

		// return if there's no rewards to give out
		if rewardPool.IsZero() {
			continue
		}

		periodRewards = periodRewards.Add(sdk.NewDecCoinFromDec(
			denom,
			math.LegacyNewDecFromInt(rewardPool.Amount).Mul(distributionRatio),
		))
	}

	// Dole out rewards
	var distributedReward sdk.Coins
	for _, winner := range ballotWinners {
		receiverVal, err := k.StakingKeeper.Validator(ctx, winner.Recipient)
		if err != nil {
			return err
		}

		// Reflects contribution
		rewardCoins, _ := periodRewards.MulDec(math.LegacyNewDec(winner.Weight).QuoInt64(ballotPowerSum)).TruncateDecimal()

		// In case absence of the validator, we just skip distribution
		if receiverVal != nil && !rewardCoins.IsZero() {
			err = k.distrKeeper.AllocateTokensToValidator(ctx, receiverVal, sdk.NewDecCoinsFromCoins(rewardCoins...))
			if err != nil {
				return err
			}
			distributedReward = distributedReward.Add(rewardCoins...)
		}
	}

	// Move distributed reward to distribution module
	err := k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.distrName, distributedReward)
	if err != nil {
		return errors.Wrap(err, "[oracle] Failed to send coins to distribution module")
	}

	return nil
}
