package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/Team-Kujira/core/x/distrib/types"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		bankKeeper bankkeeper.Keeper
		distrKeeper distrkeeper.Keeper
		stakingKeeper *stakingkeeper.Keeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey storetypes.StoreKey,
	bankKeeper bankkeeper.Keeper,
	distrKeeper distrkeeper.Keeper,
	stakingKeeper *stakingkeeper.Keeper,
) Keeper {
	return Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		bankKeeper: bankKeeper,
		distrKeeper: distrKeeper,
		stakingKeeper: stakingKeeper,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// withdraw all delegation rewards for a delegator
func (k Keeper) WithdrawAllDelegationRewards(ctx sdk.Context, delAddr sdk.AccAddress) (sdk.Coins, error) {
	return k.withdrawAllDelegationRewards(ctx, delAddr)
}
