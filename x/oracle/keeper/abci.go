package keeper

import (
	"time"

	"github.com/Team-Kujira/core/x/oracle/types"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker is called at the end of every block
func (k Keeper) EndBlocker(_ sdk.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)

	return nil
}
