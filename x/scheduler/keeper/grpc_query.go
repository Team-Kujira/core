package keeper

import (
	"github.com/Team-Kujira/core/x/scheduler/types"
)

var _ types.QueryServer = Keeper{}
