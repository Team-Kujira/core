package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// increment the reference count for a historical rewards value
func (k Keeper) incrementReferenceCount(ctx sdk.Context, valAddr sdk.ValAddress, period uint64) {
	historical := k.distrKeeper.GetValidatorHistoricalRewards(ctx, valAddr, period)
	if historical.ReferenceCount > 2 {
		panic("reference count should never exceed 2")
	}
	historical.ReferenceCount++
	k.distrKeeper.SetValidatorHistoricalRewards(ctx, valAddr, period, historical)
}

// decrement the reference count for a historical rewards value, and delete if zero references remain
func (k Keeper) decrementReferenceCount(ctx sdk.Context, valAddr sdk.ValAddress, period uint64) {
	historical := k.distrKeeper.GetValidatorHistoricalRewards(ctx, valAddr, period)
	if historical.ReferenceCount == 0 {
		panic("cannot set negative reference count")
	}
	historical.ReferenceCount--
	if historical.ReferenceCount == 0 {
		k.distrKeeper.DeleteValidatorHistoricalReward(ctx, valAddr, period)
	} else {
		k.distrKeeper.SetValidatorHistoricalRewards(ctx, valAddr, period, historical)
	}
}

// initialize starting info for a new delegation
func (k Keeper) initializeDelegation(ctx sdk.Context, val sdk.ValAddress, del sdk.AccAddress) {
	// period has already been incremented - we want to store the period ended by this delegation action
	previousPeriod := k.distrKeeper.GetValidatorCurrentRewards(ctx, val).Period - 1

	// increment reference count for the period we're going to track
	k.incrementReferenceCount(ctx, val, previousPeriod)

	validator := k.stakingKeeper.Validator(ctx, val)
	delegation := k.stakingKeeper.Delegation(ctx, del, val)

	// calculate delegation stake in tokens
	// we don't store directly, so multiply delegation shares * (tokens per share)
	// note: necessary to truncate so we don't allow withdrawing more rewards than owed
	stake := validator.TokensFromSharesTruncated(delegation.GetShares())
	k.distrKeeper.SetDelegatorStartingInfo(ctx, val, del, types.NewDelegatorStartingInfo(previousPeriod, stake, uint64(ctx.BlockHeight())))
}

func (k Keeper) withdrawAllDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress) (sdk.Coins, error) {
	rewardsTotal := sdk.Coins{}
	k.stakingKeeper.IterateDelegations(ctx, delAddr, func(_ int64, del stakingtypes.DelegationI) (stop bool) {
		valAddr := del.GetValidatorAddr()
		val := k.stakingKeeper.Validator(ctx, valAddr)
		// end current period and calculate rewards
		endingPeriod := k.distrKeeper.IncrementValidatorPeriod(ctx, val)
		rewardsRaw := k.distrKeeper.CalculateDelegationRewards(ctx, val, del, endingPeriod)
		outstanding := k.distrKeeper.GetValidatorOutstandingRewardsCoins(ctx, del.GetValidatorAddr())
		// defensive edge case may happen on the very final digits
		// of the decCoins due to operation order of the distribution mechanism.
		rewards := rewardsRaw.Intersect(outstanding)
		if !rewards.IsEqual(rewardsRaw) {
			logger := k.Logger(ctx)
			logger.Info(
				"rounding error withdrawing rewards from validator",
				"delegator", del.GetDelegatorAddr().String(),
				"validator", val.GetOperator().String(),
				"got", rewards.String(),
				"expected", rewardsRaw.String(),
			)
		}
		// truncate reward dec coins, return remainder to community pool
		finalRewards, remainder := rewards.TruncateDecimal()
		k.distrKeeper.SetValidatorOutstandingRewards(ctx, del.GetValidatorAddr(), types.ValidatorOutstandingRewards{Rewards: outstanding.Sub(rewards)})
		feePool := k.distrKeeper.GetFeePool(ctx)
		feePool.CommunityPool = feePool.CommunityPool.Add(remainder...)
		k.distrKeeper.SetFeePool(ctx, feePool)
		// decrement reference count of starting period
		startingInfo := k.distrKeeper.GetDelegatorStartingInfo(ctx, del.GetValidatorAddr(), del.GetDelegatorAddr())
		startingPeriod := startingInfo.PreviousPeriod
		k.decrementReferenceCount(ctx, del.GetValidatorAddr(), startingPeriod)
		k.initializeDelegation(ctx, valAddr, delAddr)
		rewardsTotal = rewardsTotal.Add(finalRewards...)
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeWithdrawRewards,
				sdk.NewAttribute(sdk.AttributeKeyAmount, finalRewards.String()),
				sdk.NewAttribute(types.AttributeKeyValidator, val.GetOperator().String()),
			),
		)
		return false
	})
	// add coins to user account
	if !rewardsTotal.IsZero() {
		withdrawAddr := k.distrKeeper.GetDelegatorWithdrawAddr(ctx, delAddr)
		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawAddr, rewardsTotal)
		if err != nil {
			return nil, err
		}
	} else {
		baseDenom, _ := sdk.GetBaseDenom()
		rewardsTotal = sdk.Coins{sdk.Coin{
			Denom:  baseDenom,
			Amount: math.ZeroInt(),
		}}
	}
	return rewardsTotal, nil
}