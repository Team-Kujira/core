package keeper

import (
	"time"

	"github.com/Team-Kujira/core/x/oracle/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker is called at the end of every block
func (k Keeper) EndBlocker(ctx sdk.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	// Set exchange rate snapshots
	params := k.GetParams(ctx)
	for _, epoch := range params.ExchangeRateSnapEpochs {
		for _, denom := range params.RequiredDenoms {
			latestRate := k.LatestHistoricalExchangeRateByEpochDenom(ctx, epoch.Epoch, denom)
			if latestRate.Timestamp == 0 || latestRate.Timestamp+epoch.Duration <= ctx.BlockTime().Unix() {
				rate, err := k.GetExchangeRate(ctx, denom)
				if err == nil {
					latestRate = types.HistoricalExchangeRate{
						Epoch:        epoch.Epoch,
						Timestamp:    ctx.BlockTime().Unix(),
						Denom:        denom,
						ExchangeRate: rate,
					}
					k.SetHistoricalExchangeRate(ctx, latestRate)
				}
			}

			oldestRate := k.OldestHistoricalExchangeRateByEpochDenom(ctx, epoch.Epoch, denom)
			if oldestRate.Timestamp != 0 && ctx.BlockTime().Unix() > oldestRate.Timestamp+epoch.Duration*epoch.MaxCount {
				k.DeleteHistoricalExchangeRate(ctx, oldestRate)
			}
		}
	}
	return nil
}
