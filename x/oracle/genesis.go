package oracle

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Team-Kujira/core/x/oracle/keeper"
	"github.com/Team-Kujira/core/x/oracle/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data *types.GenesisState) {
	for _, ex := range data.ExchangeRates {
		keeper.SetExchangeRate(ctx, ex.Denom, ex.ExchangeRate)
	}

	for _, mc := range data.MissCounters {
		operator, err := sdk.ValAddressFromBech32(mc.ValidatorAddress)
		if err != nil {
			panic(err)
		}

		keeper.SetMissCounter(ctx, operator, mc.MissCounter)
	}

	for _, r := range data.HistoricalExchangeRates {
		keeper.SetHistoricalExchangeRate(ctx, r)
	}

	err := keeper.SetParams(ctx, data.Params)
	if err != nil {
		panic(err)
	}

	// check if the module account exists
	moduleAcc := keeper.GetOracleAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	params := keeper.GetParams(ctx)

	exchangeRates := []types.ExchangeRateTuple{}
	keeper.IterateExchangeRates(ctx, func(denom string, rate math.LegacyDec) (stop bool) {
		exchangeRates = append(exchangeRates, types.ExchangeRateTuple{Denom: denom, ExchangeRate: rate})
		return false
	})

	historicalRates := []types.HistoricalExchangeRate{}
	keeper.IterateHistoricalExchangeRate(ctx, func(rate types.HistoricalExchangeRate) (stop bool) {
		historicalRates = append(historicalRates, rate)
		return false
	})

	missCounters := []types.MissCounter{}
	keeper.IterateMissCounters(ctx, func(operator sdk.ValAddress, missCounter uint64) (stop bool) {
		missCounters = append(missCounters, types.MissCounter{
			ValidatorAddress: operator.String(),
			MissCounter:      missCounter,
		})
		return false
	})

	return types.NewGenesisState(
		params,
		exchangeRates,
		missCounters,
		historicalRates,
	)
}
