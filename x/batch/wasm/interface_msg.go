package wasm

import (
	"cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	batchkeeper "github.com/Team-Kujira/core/x/batch/keeper"
	batchtypes "github.com/Team-Kujira/core/x/batch/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type BatchMsg struct {
	// Contracts can withdraw all rewards for the delegations from the contract.
	WithdrawAllDelegatorRewards *WithdrawAllDelegatorRewards `json:"withdrawAllDelegatorRewards,omitempty"`
}

type WithdrawAllDelegatorRewards struct{}

// withdrawAllDelegatorRewards withdraw all delegation rewards for the delegations from the contract address
func withdrawAllDelegatorRewards(ctx sdk.Context, contractAddr sdk.AccAddress, withdrawAllDelegatorRewards *WithdrawAllDelegatorRewards, bk batchkeeper.Keeper) ([]sdk.Event, [][]byte, error) {
	err := PerformWithdrawAllDelegatorRewards(bk, ctx, contractAddr, withdrawAllDelegatorRewards)
	if err != nil {
		return nil, nil, errors.Wrap(err, "perform withdrawAllDelegatorRewards")
	}
	return nil, nil, nil
}

// PerformWithdrawAllDelegatorRewards is used to perform delegation rewards from the contract delegations
func PerformWithdrawAllDelegatorRewards(bk batchkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, withdrawAllDelegatorRewards *WithdrawAllDelegatorRewards) error {
	msgServer := batchkeeper.NewMsgServerImpl(bk)
	msgWithdrawAllDelegatorRewards := batchtypes.NewMsgWithdrawAllDelegatorRewards(contractAddr)

	if err := msgWithdrawAllDelegatorRewards.ValidateBasic(); err != nil {
		return errors.Wrap(err, "failed validating MsgWithdrawAllDelegatorRewards")
	}

	_, err := msgServer.WithdrawAllDelegatorRewards(
		sdk.WrapSDKContext(ctx),
		msgWithdrawAllDelegatorRewards,
	)
	if err != nil {
		return errors.Wrap(err, "batch claim")
	}
	return nil
}

// QueryCustom implements custom msg interface
func HandleMsg(dk batchkeeper.Keeper, contractAddr sdk.AccAddress, ctx sdk.Context, q *BatchMsg) ([]sdk.Event, [][]byte, error) {
	if q.WithdrawAllDelegatorRewards != nil {
		return withdrawAllDelegatorRewards(ctx, contractAddr, q.WithdrawAllDelegatorRewards, dk)
	}

	return nil, nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown Custom variant"}
}
