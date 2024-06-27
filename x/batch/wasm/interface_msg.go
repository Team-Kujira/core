package wasm

import (
	"cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/v2/types"
	batchkeeper "github.com/Team-Kujira/core/x/batch/keeper"
	batchtypes "github.com/Team-Kujira/core/x/batch/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

type BatchMsg struct {
	// Contracts can withdraw all rewards for the delegations from the contract.
	WithdrawAllDelegatorRewards *WithdrawAllDelegatorRewards `json:"withdrawAllDelegatorRewards,omitempty"`
}

type WithdrawAllDelegatorRewards struct{}

// withdrawAllDelegatorRewards withdraw all delegation rewards for the delegations from the contract address
func withdrawAllDelegatorRewards(
	ctx sdk.Context,
	contractAddr sdk.AccAddress,
	withdrawAllDelegatorRewards *WithdrawAllDelegatorRewards,
	bk batchkeeper.Keeper,
) (*batchtypes.MsgWithdrawAllDelegatorRewardsResponse, error) {
	res, err := PerformWithdrawAllDelegatorRewards(bk, ctx, contractAddr, withdrawAllDelegatorRewards)
	if err != nil {
		return nil, errors.Wrap(err, "perform withdrawAllDelegatorRewards")
	}
	return res, nil
}

// PerformWithdrawAllDelegatorRewards is used to perform delegation rewards from the contract delegations
func PerformWithdrawAllDelegatorRewards(
	bk batchkeeper.Keeper,
	ctx sdk.Context,
	contractAddr sdk.AccAddress,
	_ *WithdrawAllDelegatorRewards,
) (*batchtypes.MsgWithdrawAllDelegatorRewardsResponse, error) {
	msgServer := batchkeeper.NewMsgServerImpl(bk)
	msgWithdrawAllDelegatorRewards := batchtypes.NewMsgWithdrawAllDelegatorRewards(contractAddr)

	if err := msgWithdrawAllDelegatorRewards.ValidateBasic(); err != nil {
		return nil, errors.Wrap(err, "failed validating MsgWithdrawAllDelegatorRewards")
	}

	res, err := msgServer.WithdrawAllDelegatorRewards(
		ctx,
		msgWithdrawAllDelegatorRewards,
	)
	if err != nil {
		return nil, errors.Wrap(err, "batch claim")
	}
	return res, nil
}

// QueryCustom implements custom msg interface
func HandleMsg(
	dk batchkeeper.Keeper,
	contractAddr sdk.AccAddress,
	ctx sdk.Context,
	q *BatchMsg,
) ([]sdk.Event, [][]byte, [][]*codectypes.Any, error) {
	var res proto.Message
	var err error

	if q.WithdrawAllDelegatorRewards != nil {
		res, err = withdrawAllDelegatorRewards(ctx, contractAddr, q.WithdrawAllDelegatorRewards, dk)
	}

	if err != nil {
		return nil, nil, nil, err
	}

	if res == nil {
		return nil, nil, nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown Custom variant"}
	}

	x, err := codectypes.NewAnyWithValue(res)
	if err != nil {
		return nil, nil, nil, err
	}
	msgResponses := [][]*codectypes.Any{{x}}

	return nil, nil, msgResponses, err
}
