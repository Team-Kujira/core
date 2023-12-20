package bindings

import (
	"cosmossdk.io/errors"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	alliancekeeper "github.com/terra-money/alliance/x/alliance/keeper"
	alliancetypes "github.com/terra-money/alliance/x/alliance/types"
)

// Messages
type AllianceMsg struct {
	// Contracts can delegate from the contract.
	Delegate *Delegate `json:"delegate,omitempty"`
	// Contracts can redelegate from the contract.
	Redelegate *Redelegate `json:"redelegate,omitempty"`
	// Contracts can undelegate from the contract.
	Undelegate *Undelegate `json:"undelegate,omitempty"`
	// Contracts can claim delegation rewards from the contract.
	ClaimDelegationRewards *ClaimDelegationRewards `json:"claim_delegation_rewards,omitempty"`
}

type Delegate struct {
	ValidatorAddress string
	Amount           wasmvmtypes.Coin
}

type Redelegate struct {
	ValidatorSrcAddress string
	ValidatorDstAddress string
	Amount              wasmvmtypes.Coin
}

type Undelegate struct {
	ValidatorAddress string
	Amount           wasmvmtypes.Coin
}

type ClaimDelegationRewards struct {
	ValidatorAddress string
	Denom            string
}

// Queries
type AllianceQuery struct {
	Params      *Params      `json:"params,omitempty"`
	Alliance    *Alliance    `json:"alliance,omitempty"`
	IBCAlliance *IBCAlliance `json:"ibc_alliance,omitempty"`
}

type Params struct{}

type Alliance struct {
	Denom string
}

type IBCAlliance struct {
	Hash string
}

func HandleAllianceQuery(ctx sdk.Context, keeper alliancekeeper.Keeper, q *AllianceQuery) (any, error) {
	switch {
	case q.Params != nil:
		return nil, nil

	case q.Alliance != nil:
		return nil, nil

	case q.IBCAlliance != nil:
		return nil, nil
	default:
		return nil, errors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized query request")
	}
}

func HandleAllianceMsg(ctx sdk.Context, keeper alliancekeeper.Keeper, contractAddr sdk.AccAddress, msg *AllianceMsg) ([]sdk.Event, [][]byte, error) {
	if msg.Delegate != nil {
		return delegate(ctx, contractAddr, msg.Delegate, keeper)
	}
	if msg.Redelegate != nil {
		return redelegate(ctx, contractAddr, msg.Redelegate, keeper)
	}
	if msg.Undelegate != nil {
		return undelegate(ctx, contractAddr, msg.Undelegate, keeper)
	}
	if msg.ClaimDelegationRewards != nil {
		return claimDelegationRewards(ctx, contractAddr, msg.ClaimDelegationRewards, keeper)
	}
	return nil, nil, wasmvmtypes.InvalidRequest{Err: "unknown ICA Message variant"}
}

func delegate(ctx sdk.Context, contractAddr sdk.AccAddress, delegate *Delegate, allianceKeeper alliancekeeper.Keeper) ([]sdk.Event, [][]byte, error) {
	_, err := PerformDelegate(allianceKeeper, ctx, contractAddr, delegate)
	if err != nil {
		return nil, nil, errors.Wrap(err, "perform delegate")
	}
	return nil, nil, nil
}

func PerformDelegate(allianceKeeper alliancekeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msg *Delegate) (*alliancetypes.MsgDelegateResponse, error) {
	if msg == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "delegate null message"}
	}

	msgServer := alliancekeeper.NewMsgServerImpl(allianceKeeper)

	amount, err := wasmkeeper.ConvertWasmCoinToSdkCoin(msg.Amount)
	if err != nil {
		return nil, err
	}
	sdkMsg := alliancetypes.NewMsgDelegate(contractAddr.String(), msg.ValidatorAddress, amount)

	if err := sdkMsg.ValidateBasic(); err != nil {
		return nil, errors.Wrap(err, "failed validating MsgDelegate")
	}

	res, err := msgServer.Delegate(sdk.WrapSDKContext(ctx), sdkMsg)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func redelegate(ctx sdk.Context, contractAddr sdk.AccAddress, redelegate *Redelegate, allianceKeeper alliancekeeper.Keeper) ([]sdk.Event, [][]byte, error) {
	_, err := PerformRedelegate(allianceKeeper, ctx, contractAddr, redelegate)
	if err != nil {
		return nil, nil, errors.Wrap(err, "perform redelegate")
	}
	return nil, nil, nil
}

func PerformRedelegate(allianceKeeper alliancekeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msg *Redelegate) (*alliancetypes.MsgRedelegateResponse, error) {
	if msg == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "redelegate null message"}
	}

	msgServer := alliancekeeper.NewMsgServerImpl(allianceKeeper)

	amount, err := wasmkeeper.ConvertWasmCoinToSdkCoin(msg.Amount)
	if err != nil {
		return nil, err
	}
	sdkMsg := alliancetypes.NewMsgRedelegate(
		contractAddr.String(),
		msg.ValidatorSrcAddress,
		msg.ValidatorDstAddress,
		amount,
	)

	if err := sdkMsg.ValidateBasic(); err != nil {
		return nil, errors.Wrap(err, "failed validating MsgRedelegate")
	}

	res, err := msgServer.Redelegate(sdk.WrapSDKContext(ctx), sdkMsg)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func undelegate(ctx sdk.Context, contractAddr sdk.AccAddress, undelegate *Undelegate, allianceKeeper alliancekeeper.Keeper) ([]sdk.Event, [][]byte, error) {
	_, err := PerformUndelegate(allianceKeeper, ctx, contractAddr, undelegate)
	if err != nil {
		return nil, nil, errors.Wrap(err, "perform undelegate")
	}
	return nil, nil, nil
}

func PerformUndelegate(allianceKeeper alliancekeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msg *Undelegate) (*alliancetypes.MsgUndelegateResponse, error) {
	if msg == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "undelegate null message"}
	}

	msgServer := alliancekeeper.NewMsgServerImpl(allianceKeeper)

	amount, err := wasmkeeper.ConvertWasmCoinToSdkCoin(msg.Amount)
	if err != nil {
		return nil, err
	}
	sdkMsg := alliancetypes.NewMsgUndelegate(
		contractAddr.String(),
		msg.ValidatorAddress,
		amount,
	)

	if err := sdkMsg.ValidateBasic(); err != nil {
		return nil, errors.Wrap(err, "failed validating MsgUndelegate")
	}

	res, err := msgServer.Undelegate(sdk.WrapSDKContext(ctx), sdkMsg)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func claimDelegationRewards(ctx sdk.Context, contractAddr sdk.AccAddress, claim *ClaimDelegationRewards, allianceKeeper alliancekeeper.Keeper) ([]sdk.Event, [][]byte, error) {
	_, err := PerformClaimDelegationRewards(allianceKeeper, ctx, contractAddr, claim)
	if err != nil {
		return nil, nil, errors.Wrap(err, "perform delegation")
	}
	return nil, nil, nil
}

func PerformClaimDelegationRewards(allianceKeeper alliancekeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msg *ClaimDelegationRewards) (*alliancetypes.MsgClaimDelegationRewardsResponse, error) {
	if msg == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "claimDelegationRewards null message"}
	}

	msgServer := alliancekeeper.NewMsgServerImpl(allianceKeeper)

	sdkMsg := alliancetypes.NewMsgClaimDelegationRewards(
		contractAddr.String(),
		msg.ValidatorAddress,
		msg.Denom,
	)

	if err := sdkMsg.ValidateBasic(); err != nil {
		return nil, errors.Wrap(err, "failed validating MsgClaimDelegationRewards")
	}

	res, err := msgServer.ClaimDelegationRewards(sdk.WrapSDKContext(ctx), sdkMsg)
	if err != nil {
		return nil, err
	}

	return res, nil
}
