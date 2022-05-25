package keeper

import (
	"kujira/x/scheduler/types"
)

var _ types.QueryServer = Keeper{}
