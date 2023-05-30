package wasmbinding

import (
	denomkeeper "github.com/Team-Kujira/core/x/denom/keeper"

	oraclekeeper "github.com/Team-Kujira/core/x/oracle/keeper"

	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"

	intertxkeeper "github.com/Team-Kujira/core/x/inter-tx/keeper"
)

type QueryPlugin struct {
	denomKeeper   denomkeeper.Keeper
	bankkeeper    bankkeeper.Keeper
	oraclekeeper  oraclekeeper.Keeper
	intertxkeeper intertxkeeper.Keeper
}

// NewQueryPlugin returns a reference to a new QueryPlugin.
func NewQueryPlugin(bk bankkeeper.Keeper, ok oraclekeeper.Keeper, dk denomkeeper.Keeper, itxk intertxkeeper.Keeper) *QueryPlugin {
	return &QueryPlugin{
		denomKeeper:   dk,
		bankkeeper:    bk,
		oraclekeeper:  ok,
		intertxkeeper: itxk,
	}
}
