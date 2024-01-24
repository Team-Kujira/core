package keeper

import (
	"fmt"

	gogotypes "github.com/cosmos/gogoproto/types"

	"cosmossdk.io/errors"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/Team-Kujira/core/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Keeper of the oracle store
type Keeper struct {
	cdc        codec.BinaryCodec
	storeKey   storetypes.StoreKey
	paramSpace paramstypes.Subspace

	accountKeeper  types.AccountKeeper
	bankKeeper     types.BankKeeper
	distrKeeper    types.DistributionKeeper
	SlashingKeeper types.SlashingKeeper
	StakingKeeper  types.StakingKeeper

	distrName   string
	rewardDenom string
	authority   string
}

// NewKeeper constructs a new keeper for oracle
func NewKeeper(cdc codec.BinaryCodec, storeKey storetypes.StoreKey,
	paramspace paramstypes.Subspace, accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper, distrKeeper types.DistributionKeeper,
	slashingkeeper types.SlashingKeeper, stakingKeeper types.StakingKeeper, distrName string, authority string,
) Keeper {
	// ensure oracle module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// set KeyTable if it has not already been set
	if !paramspace.HasKeyTable() {
		paramspace = paramspace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		cdc:            cdc,
		storeKey:       storeKey,
		paramSpace:     paramspace,
		accountKeeper:  accountKeeper,
		bankKeeper:     bankKeeper,
		distrKeeper:    distrKeeper,
		SlashingKeeper: slashingkeeper,
		StakingKeeper:  stakingKeeper,
		distrName:      distrName,
		rewardDenom:    "ukuji",
		authority:      authority,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

//-----------------------------------
// ExchangeRate logic

// GetExchangeRate gets the consensus exchange rate of the denom asset from the store.
func (k Keeper) GetExchangeRate(ctx sdk.Context, denom string) (math.LegacyDec, error) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.GetExchangeRateKey(denom))
	if b == nil {
		return math.LegacyZeroDec(), errors.Wrap(types.ErrUnknownDenom, denom)
	}

	dp := sdk.DecProto{}
	k.cdc.MustUnmarshal(b, &dp)
	return dp.Dec, nil
}

// SetExchangeRate sets the consensus exchange rate of the denom asset to the store.
func (k Keeper) SetExchangeRate(ctx sdk.Context, denom string, exchangeRate math.LegacyDec) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&sdk.DecProto{Dec: exchangeRate})
	store.Set(types.GetExchangeRateKey(denom), bz)
}

// SetExchangeRateWithEvent sets the consensus exchange rate of the denom asset to the store with ABCI event
func (k Keeper) SetExchangeRateWithEvent(ctx sdk.Context, denom string, exchangeRate math.LegacyDec) {
	k.SetExchangeRate(ctx, denom, exchangeRate)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventTypeExchangeRateUpdate,
			sdk.NewAttribute(types.AttributeKeyDenom, denom),
			sdk.NewAttribute(types.AttributeKeyExchangeRate, exchangeRate.String()),
		),
	)
}

// DeleteExchangeRate deletes the consensus exchange rate of the denom asset from the store.
func (k Keeper) DeleteExchangeRate(ctx sdk.Context, denom string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetExchangeRateKey(denom))
}

// IterateExchangeRates iterates over luna rates in the store
func (k Keeper) IterateExchangeRates(ctx sdk.Context, handler func(denom string, exchangeRate math.LegacyDec) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := storetypes.KVStorePrefixIterator(store, types.ExchangeRateKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		denom := string(iter.Key()[len(types.ExchangeRateKey):])
		dp := sdk.DecProto{}
		k.cdc.MustUnmarshal(iter.Value(), &dp)
		if handler(denom, dp.Dec) {
			break
		}
	}
}

//-----------------------------------
// Miss counter logic

// GetMissCounter retrieves the # of vote periods missed in this oracle slash window
func (k Keeper) GetMissCounter(ctx sdk.Context, operator sdk.ValAddress) uint64 {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetMissCounterKey(operator))
	if bz == nil {
		// By default the counter is zero
		return 0
	}

	var missCounter gogotypes.UInt64Value
	k.cdc.MustUnmarshal(bz, &missCounter)
	return missCounter.Value
}

// SetMissCounter updates the # of vote periods missed in this oracle slash window
func (k Keeper) SetMissCounter(ctx sdk.Context, operator sdk.ValAddress, missCounter uint64) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&gogotypes.UInt64Value{Value: missCounter})
	store.Set(types.GetMissCounterKey(operator), bz)
}

// DeleteMissCounter removes miss counter for the validator
func (k Keeper) DeleteMissCounter(ctx sdk.Context, operator sdk.ValAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetMissCounterKey(operator))
}

// IterateMissCounters iterates over the miss counters and performs a callback function.
func (k Keeper) IterateMissCounters(ctx sdk.Context,
	handler func(operator sdk.ValAddress, missCounter uint64) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)
	iter := storetypes.KVStorePrefixIterator(store, types.MissCounterKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		operator := sdk.ValAddress(iter.Key()[2:])

		var missCounter gogotypes.UInt64Value
		k.cdc.MustUnmarshal(iter.Value(), &missCounter)

		if handler(operator, missCounter.Value) {
			break
		}
	}
}

//-----------------------------------
// AggregateExchangeRateVote logic

func (k Keeper) GetSubspace() paramstypes.Subspace {
	return k.paramSpace
}

func (k Keeper) SetOraclePrices(ctx sdk.Context, prices map[string]math.LegacyDec) error {
	for b, q := range prices {
		k.SetExchangeRateWithEvent(ctx, b, q)
	}
	return nil
}

type (
	CurrencyPair struct {
		Base  string
		Quote string
	}
)

func (cp CurrencyPair) String() string {
	return cp.Base + cp.Quote
}
