package simulation

// DONTCOVER

import (
	"github.com/Team-Kujira/core/x/oracle/keeper"
	"github.com/Team-Kujira/core/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/codec"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// Simulation operation weights constants
//
//nolint:gosec //these aren't hard coded credentials
const (
	salt = "fc5bb0bc63e54b2918d9334bf3259f5dc575e8d7a4df4e836dd80f1ad62aa89b"
)

var (
	whitelist                          = []string{types.TestDenomA, types.TestDenomB, types.TestDenomC}
	voteHashMap                        = make(map[string]string)
	DefaultWeightMsgSend               = 100
	DefaultWeightMsgSetWithdrawAddress = 50
)

// WeightedOperations returns all the operations from the module with their respective weights
func WeightedOperations(
	appParams simtypes.AppParams,
	_ codec.JSONCodec,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simulation.WeightedOperations {
	return simulation.WeightedOperations{}
}
