package keeper

import (
	"time"

	"github.com/Team-Kujira/core/x/batch/types"
	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// increment the reference count for a historical rewards value
// this func was copied from
// https://github.com/cosmos/cosmos-sdk/blob/main/x/distribution/keeper/validator.go
func (k Keeper) incrementReferenceCount(ctx sdk.Context, valAddr sdk.ValAddress, period uint64) {
	historical := k.distrKeeper.GetValidatorHistoricalRewards(ctx, valAddr, period)
	if historical.ReferenceCount > 2 {
		panic("reference count should never exceed 2")
	}
	historical.ReferenceCount++
	k.distrKeeper.SetValidatorHistoricalRewards(ctx, valAddr, period, historical)
}

// decrement the reference count for a historical rewards value, and delete if zero references remain
// this func was copied from
// https://github.com/cosmos/cosmos-sdk/blob/main/x/distribution/keeper/validator.go
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
// this func was copied from
// https://github.com/cosmos/cosmos-sdk/blob/main/x/distribution/keeper/delegation.go
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
	k.distrKeeper.SetDelegatorStartingInfo(ctx, val, del, distrtypes.NewDelegatorStartingInfo(previousPeriod, stake, uint64(ctx.BlockHeight())))
}

func (k Keeper) withdrawAllDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress) (sdk.Coins, error) {
	rewardsTotal := sdk.Coins{}
	remainderTotal := sdk.DecCoins{}
	// callback func was referenced from withdrawDelegationRewards func in
	// https://github.com/cosmos/cosmos-sdk/blob/main/x/distribution/keeper/delegation.go
	k.stakingKeeper.IterateDelegations(ctx, delAddr, func(_ int64, del stakingtypes.DelegationI) (stop bool) {
		valAddr := del.GetValidatorAddr()
		val := k.stakingKeeper.Validator(ctx, valAddr)

		// check existence of delegator starting info
		if !k.distrKeeper.HasDelegatorStartingInfo(ctx, del.GetValidatorAddr(), del.GetDelegatorAddr()) {
			return false
		}

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
		remainderTotal = remainderTotal.Add(remainder...)
		rewardsTotal = rewardsTotal.Add(finalRewards...)

		// update the outstanding rewards and the community pool only if the
		// transaction was successful
		k.distrKeeper.SetValidatorOutstandingRewards(ctx, del.GetValidatorAddr(), distrtypes.ValidatorOutstandingRewards{Rewards: outstanding.Sub(rewards)})

		// decrement reference count of starting period
		startingInfo := k.distrKeeper.GetDelegatorStartingInfo(ctx, del.GetValidatorAddr(), del.GetDelegatorAddr())
		startingPeriod := startingInfo.PreviousPeriod
		k.decrementReferenceCount(ctx, del.GetValidatorAddr(), startingPeriod)

		// reinitialize the delegation
		k.initializeDelegation(ctx, valAddr, delAddr)

		return false
	})

	// distribute total remainder to community pool
	feePool := k.distrKeeper.GetFeePool(ctx)
	feePool.CommunityPool = feePool.CommunityPool.Add(remainderTotal...)
	k.distrKeeper.SetFeePool(ctx, feePool)

	// add total reward coins to user account
	if !rewardsTotal.IsZero() {
		withdrawAddr := k.distrKeeper.GetDelegatorWithdrawAddr(ctx, delAddr)
		err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, distrtypes.ModuleName, withdrawAddr, rewardsTotal)
		if err != nil {
			return nil, err
		}
	} else {
		rewardsTotal = sdk.Coins{}
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			distrtypes.EventTypeWithdrawRewards,
			sdk.NewAttribute(sdk.AttributeKeyAmount, rewardsTotal.String()),
			sdk.NewAttribute(distrtypes.AttributeKeyDelegator, delAddr.String()),
		),
	)
	return rewardsTotal, nil
}

func (k Keeper) batchResetDelegation(ctx sdk.Context, msg *types.MsgBatchResetDelegation) error {
	delAddr, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return err
	}

	if len(msg.Validators) != len(msg.Amounts) {
		return types.ErrValidatorsAndAmountsMismatch
	}

	bondDenom := k.stakingKeeper.BondDenom(ctx)

	for i, valStr := range msg.Validators {
		valAddr, valErr := sdk.ValAddressFromBech32(valStr)
		if valErr != nil {
			return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", valErr)
		}

		validator, found := k.stakingKeeper.GetValidator(ctx, valAddr)
		if !found {
			return stakingtypes.ErrNoValidatorFound
		}

		targetAmount := msg.Amounts[i]
		currAmount := sdk.ZeroInt()
		delegation, found := k.stakingKeeper.GetDelegation(ctx, delAddr, valAddr)
		if found {
			currAmount = validator.TokensFromShares(delegation.Shares).RoundInt()
		}

		if currAmount.Equal(targetAmount) {
			continue
		}

		if currAmount.LT(targetAmount) {
			amount := targetAmount.Sub(currAmount)
			// NOTE: source funds are always unbonded
			newShares, err := k.stakingKeeper.Delegate(ctx, delAddr, amount, stakingtypes.Unbonded, validator, true)
			if err != nil {
				return err
			}

			if amount.IsInt64() {
				defer func() {
					telemetry.IncrCounter(1, types.ModuleName, "delegate")
					telemetry.SetGaugeWithLabels(
						[]string{"tx", "msg", sdk.MsgTypeURL(msg)},
						float32(amount.Int64()),
						[]metrics.Label{telemetry.NewLabel("denom", bondDenom)},
					)
				}()
			}

			ctx.EventManager().EmitEvents(sdk.Events{
				sdk.NewEvent(
					stakingtypes.EventTypeDelegate,
					sdk.NewAttribute(stakingtypes.AttributeKeyValidator, valStr),
					sdk.NewAttribute(stakingtypes.AttributeKeyDelegator, msg.DelegatorAddress),
					sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
					sdk.NewAttribute(stakingtypes.AttributeKeyNewShares, newShares.String()),
				),
			})
		} else {
			amount := currAmount.Sub(targetAmount)
			shares, err := k.stakingKeeper.ValidateUnbondAmount(
				ctx, delAddr, valAddr, amount,
			)
			if err != nil {
				return err
			}

			completionTime, err := k.stakingKeeper.Undelegate(ctx, delAddr, valAddr, shares)
			if err != nil {
				return err
			}

			if amount.IsInt64() {
				defer func() {
					telemetry.IncrCounter(1, types.ModuleName, "undelegate")
					telemetry.SetGaugeWithLabels(
						[]string{"tx", "msg", msg.Type()},
						float32(amount.Int64()),
						[]metrics.Label{telemetry.NewLabel("denom", bondDenom)},
					)
				}()
			}

			ctx.EventManager().EmitEvents(sdk.Events{
				sdk.NewEvent(
					stakingtypes.EventTypeUnbond,
					sdk.NewAttribute(stakingtypes.AttributeKeyValidator, valStr),
					sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
					sdk.NewAttribute(stakingtypes.AttributeKeyCompletionTime, completionTime.Format(time.RFC3339)),
				),
			})
		}

	}
	return nil
}
