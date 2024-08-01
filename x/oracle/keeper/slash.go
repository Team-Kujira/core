package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// SlashAndResetMissCounters do slash any operator who over criteria & clear all operators miss counter to zero
func (k Keeper) SlashAndResetMissCounters(ctx sdk.Context) {
	height := ctx.BlockHeight()
	distributionHeight := height - sdk.ValidatorUpdateDelay - 1

	slashWindow := k.SlashWindow(ctx)
	minValidPerWindow := k.MinValidPerWindow(ctx)
	slashFraction := k.SlashFraction(ctx)
	powerReduction := k.StakingKeeper.PowerReduction(ctx)

	k.IterateMissCounters(ctx, func(operator sdk.ValAddress, missCounter uint64) bool {
		// Calculate valid vote rate; (SlashWindow - MissCounter)/SlashWindow
		validVoteRate := math.LegacyNewDecFromInt(
			math.NewInt(int64(slashWindow - missCounter))).
			QuoInt64(int64(slashWindow))

		// Penalize the validator whose the valid vote rate is smaller than min threshold
		if validVoteRate.LT(minValidPerWindow) {
			validator, err := k.StakingKeeper.Validator(ctx, operator)
			if err != nil {
				panic(err)
			}

			if validator.IsBonded() && !validator.IsJailed() {
				consAddr, err := validator.GetConsAddr()
				if err != nil {
					panic(err)
				}

				err = k.SlashingKeeper.Slash(
					ctx, consAddr, slashFraction,
					validator.GetConsensusPower(powerReduction), distributionHeight,
				)
				if err != nil {
					panic(err)
				}
				err = k.SlashingKeeper.Jail(ctx, consAddr)
				if err != nil {
					panic(err)
				}
			}
		}

		k.DeleteMissCounter(ctx, operator)
		return false
	})
}
