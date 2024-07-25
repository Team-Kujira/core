package simulation

// DONTCOVER

import (
	"github.com/Team-Kujira/core/x/oracle/keeper"
	"github.com/Team-Kujira/core/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/codec"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

var (
	DefaultWeightMsgSend               = 100
	DefaultWeightMsgSetWithdrawAddress = 50
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	_ simtypes.AppParams,
	_ codec.JSONCodec,
	_ types.AccountKeeper,
	_ types.BankKeeper,
	_ keeper.Keeper,
) simulation.WeightedOperations {
	return simulation.WeightedOperations{}
}
