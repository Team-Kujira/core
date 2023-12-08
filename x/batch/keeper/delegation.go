package keeper

import (
	"context"
	"time"

	"cosmossdk.io/math"
	"github.com/Team-Kujira/core/x/batch/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/hashicorp/go-metrics"
)

// increment the reference count for a historical rewards value
// this func was copied from
// https://github.com/cosmos/cosmos-sdk/blob/main/x/distribution/keeper/validator.go
func (k Keeper) incrementReferenceCount(ctx context.Context, valAddr sdk.ValAddress, period uint64) error {
	historical, err := k.distrKeeper.GetValidatorHistoricalRewards(ctx, valAddr, period)
	if err != nil {
		return err
	}
	if historical.ReferenceCount > 2 {
		panic("reference count should never exceed 2")
	}
	historical.ReferenceCount++
	return k.distrKeeper.SetValidatorHistoricalRewards(ctx, valAddr, period, historical)
}

// decrement the reference count for a historical rewards value, and delete if zero references remain
// this func was copied from
// https://github.com/cosmos/cosmos-sdk/blob/main/x/distribution/keeper/validator.go
func (k Keeper) decrementReferenceCount(ctx context.Context, valAddr sdk.ValAddress, period uint64) error {
	historical, err := k.distrKeeper.GetValidatorHistoricalRewards(ctx, valAddr, period)
	if err != nil {
		return err
	}

	if historical.ReferenceCount == 0 {
		panic("cannot set negative reference count")
	}
	historical.ReferenceCount--
	if historical.ReferenceCount == 0 {
		return k.distrKeeper.DeleteValidatorHistoricalReward(ctx, valAddr, period)
	}
	return k.distrKeeper.SetValidatorHistoricalRewards(ctx, valAddr, period, historical)
}

// initialize starting info for a new delegation
// this func was copied from
// https://github.com/cosmos/cosmos-sdk/blob/main/x/distribution/keeper/delegation.go
func (k Keeper) initializeDelegation(ctx context.Context, val sdk.ValAddress, del sdk.AccAddress) error {
	// period has already been incremented - we want to store the period ended by this delegation action
	valCurrentRewards, err := k.distrKeeper.GetValidatorCurrentRewards(ctx, val)
	if err != nil {
		return err
	}
	previousPeriod := valCurrentRewards.Period - 1

	// increment reference count for the period we're going to track
	err = k.incrementReferenceCount(ctx, val, previousPeriod)
	if err != nil {
		return err
	}

	validator, err := k.stakingKeeper.Validator(ctx, val)
	if err != nil {
		return err
	}

	delegation, err := k.stakingKeeper.Delegation(ctx, del, val)
	if err != nil {
		return err
	}

	// calculate delegation stake in tokens
	// we don't store directly, so multiply delegation shares * (tokens per share)
	// note: necessary to truncate so we don't allow withdrawing more rewards than owed
	stake := validator.TokensFromSharesTruncated(delegation.GetShares())
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return k.distrKeeper.SetDelegatorStartingInfo(ctx, val, del, distrtypes.NewDelegatorStartingInfo(previousPeriod, stake, uint64(sdkCtx.BlockHeight())))
}

func (k Keeper) withdrawAllDelegationRewards(ctx context.Context, delAddr sdk.AccAddress) (sdk.Coins, error) {
	rewardsTotal := sdk.Coins{}
	remainderTotal := sdk.DecCoins{}
	// callback func was referenced from withdrawDelegationRewards func in
	// https://github.com/cosmos/cosmos-sdk/blob/main/x/distribution/keeper/delegation.go
	k.stakingKeeper.IterateDelegations(ctx, delAddr, func(_ int64, del stakingtypes.DelegationI) (stop bool) {
		valAddr, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(del.GetValidatorAddr())
		if err != nil {
			panic(err)
		}

		val, err := k.stakingKeeper.Validator(ctx, valAddr)
		if err != nil {
			panic(err)
		}

		// check existence of delegator starting info
		hasInfo, err := k.distrKeeper.HasDelegatorStartingInfo(ctx, sdk.ValAddress(valAddr), sdk.AccAddress(delAddr))
		if err != nil {
			panic(err)
		}

		if !hasInfo {
			return false
		}

		// end current period and calculate rewards
		endingPeriod, err := k.distrKeeper.IncrementValidatorPeriod(ctx, val)
		if err != nil {
			panic(err)
		}

		rewardsRaw, err := k.distrKeeper.CalculateDelegationRewards(ctx, val, del, endingPeriod)
		if err != nil {
			panic(err)
		}

		outstanding, err := k.distrKeeper.GetValidatorOutstandingRewardsCoins(ctx, sdk.ValAddress(valAddr))
		if err != nil {
			panic(err)
		}

		// defensive edge case may happen on the very final digits
		// of the decCoins due to operation order of the distribution mechanism.
		rewards := rewardsRaw.Intersect(outstanding)
		if !rewards.Equal(rewardsRaw) {
			logger := k.Logger(ctx)
			logger.Info(
				"rounding error withdrawing rewards from validator",
				"delegator", del.GetDelegatorAddr(),
				"validator", val.GetOperator(),
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
		err = k.distrKeeper.SetValidatorOutstandingRewards(ctx, sdk.ValAddress(valAddr), distrtypes.ValidatorOutstandingRewards{Rewards: outstanding.Sub(rewards)})
		if err != nil {
			panic(err)
		}

		// decrement reference count of starting period
		startingInfo, err := k.distrKeeper.GetDelegatorStartingInfo(ctx, sdk.ValAddress(valAddr), sdk.AccAddress(delAddr))
		if err != nil {
			panic(err)
		}

		startingPeriod := startingInfo.PreviousPeriod
		err = k.decrementReferenceCount(ctx, sdk.ValAddress(valAddr), startingPeriod)
		if err != nil {
			panic(err)
		}

		// remove delegator starting info
		err = k.distrKeeper.DeleteDelegatorStartingInfo(ctx, sdk.ValAddress(valAddr), sdk.AccAddress(delAddr))
		if err != nil {
			panic(err)
		}

		// reinitialize the delegation
		err = k.initializeDelegation(ctx, valAddr, delAddr)
		if err != nil {
			panic(err)
		}

		return false
	})

	// distribute total remainder to community pool
	feePool, err := k.distrKeeper.FeePool.Get(ctx)
	if err != nil {
		panic(err)
	}

	feePool.CommunityPool = feePool.CommunityPool.Add(remainderTotal...)
	err = k.distrKeeper.FeePool.Set(ctx, feePool)
	if err != nil {
		panic(err)
	}

	// add total reward coins to user account
	if !rewardsTotal.IsZero() {
		withdrawAddr, err := k.distrKeeper.GetDelegatorWithdrawAddr(ctx, delAddr)
		if err != nil {
			return nil, err
		}

		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, distrtypes.ModuleName, withdrawAddr, rewardsTotal)
		if err != nil {
			return nil, err
		}
	} else {
		rewardsTotal = sdk.Coins{}
	}

	sdk.UnwrapSDKContext(ctx).EventManager().EmitEvent(
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

	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return err
	}

	for i, valStr := range msg.Validators {
		valAddr, valErr := sdk.ValAddressFromBech32(valStr)
		if valErr != nil {
			return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", valErr)
		}

		validator, err := k.stakingKeeper.GetValidator(ctx, valAddr)
		if err != nil {
			return err
		}

		targetAmount := msg.Amounts[i]
		currAmount := math.ZeroInt()
		delegation, err := k.stakingKeeper.GetDelegation(ctx, delAddr, valAddr)
		if err == nil {
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

			completionTime, _, err := k.stakingKeeper.Undelegate(ctx, delAddr, valAddr, shares)
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
