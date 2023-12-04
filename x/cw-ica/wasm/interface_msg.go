package wasm

import (
	"cosmossdk.io/errors"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"

	cwicakeeper "github.com/Team-Kujira/core/x/cw-ica/keeper"
	"github.com/Team-Kujira/core/x/cw-ica/types"
	cosmostypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// ProtobufAny is a hack-struct to serialize protobuf Any message into JSON object
// See https://github.com/neutron-org/neutron/blob/main/wasmbinding/bindings/msg.go
type ProtobufAny struct {
	TypeURL string `json:"type_url"`
	Value   []byte `json:"value"`
}

type ICAMsg struct {
	/// Contracts can register a new interchain account.
	Register *Register `json:"register,omitempty"`
	/// Contracts can submit transactions to the ICA
	/// associated with the given information.
	Submit *Submit `json:"submit,omitempty"`
}

// / Register creates a new interchain account.
// / If the account was created in the past, this will
// / re-establish a dropped connection, or do nothing if
// / the connection is still active.
// / The account is registered using (port, channel, sender, id)
// / as the unique identifier.
type Register struct {
	ConnectionId string `json:"connection_id"`
	AccountId    string `json:"account_id"`
	Version      string `json:"version"`
	Callback     []byte `json:"callback"`
}

// / Submit submits transactions to the ICA
// / associated with the given address.
type Submit struct {
	ConnectionId string        `json:"connection_id"`
	AccountId    string        `json:"account_id"`
	Msgs         []ProtobufAny `json:"msgs"`
	Memo         string        `json:"memo"`
	Timeout      uint64        `json:"timeout"`
	Callback     []byte        `json:"callback"`
}

func register(ctx sdk.Context, contractAddr sdk.AccAddress, register *Register, cwicak cwicakeeper.Keeper, ik icacontrollerkeeper.Keeper) ([]sdk.Event, [][]byte, error) {
	_, err := PerformRegisterICA(cwicak, ik, ctx, contractAddr, register)
	if err != nil {
		return nil, nil, errors.Wrap(err, "perform register ICA")
	}
	// Construct an sdk.Event from the MsgRegisterInterchainAccountResponse.
	// Somewhat hacky way to get the data back to the contract.
	// attrs := []sdk.Attribute{
	// 	sdk.NewAttribute()
	return nil, nil, nil
}

// PerformRegisterICA is used with register to validate the register message and register the ICA.
func PerformRegisterICA(cwicak cwicakeeper.Keeper, f icacontrollerkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msg *Register) (*icacontrollertypes.MsgRegisterInterchainAccountResponse, error) {
	if msg == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "register ICA null message"}
	}

	msgServer := icacontrollerkeeper.NewMsgServerImpl(&f)

	// format "{owner}-{id}"
	owner := contractAddr.String() + "-" + msg.AccountId
	msgRegister := icacontrollertypes.NewMsgRegisterInterchainAccount(msg.ConnectionId, owner, msg.Version)

	if err := msgRegister.ValidateBasic(); err != nil {
		return nil, errors.Wrap(err, "failed validating MsgRegisterInterchainAccount")
	}

	res, err := msgServer.RegisterInterchainAccount(
		sdk.WrapSDKContext(ctx),
		msgRegister,
	)

	portID, err := icatypes.NewControllerPortID(owner)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, errors.Wrap(err, "registering ICA")
	}

	f.SetMiddlewareEnabled(ctx, portID, msg.ConnectionId)

	cwicak.SetCallbackData(ctx, types.CallbackData{
		PortId:       portID,
		ChannelId:    "",
		Sequence:     0,
		Contract:     contractAddr.String(),
		ConnectionId: msg.ConnectionId,
		AccountId:    msg.AccountId,
		Callback:     msg.Callback,
	})

	return res, nil
}

func submit(ctx sdk.Context, contractAddr sdk.AccAddress, submitTx *Submit, cwicak cwicakeeper.Keeper, ik icacontrollerkeeper.Keeper) ([]sdk.Event, [][]byte, error) {
	_, err := PerformSubmitTxs(ik, cwicak, ctx, contractAddr, submitTx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "perform submit txs")
	}
	return nil, nil, nil
}

// PerformSubmitTxs is used with submitTxs to validate the submitTxs message and submit the txs.
func PerformSubmitTxs(f icacontrollerkeeper.Keeper, cwicak cwicakeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, submitTx *Submit) (*icacontrollertypes.MsgSendTxResponse, error) {
	if submitTx == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "submit txs null message"}
	}
	msgs := []*cosmostypes.Any{}
	for _, msg := range submitTx.Msgs {
		msgs = append(msgs, &cosmostypes.Any{
			TypeUrl: msg.TypeURL,
			Value:   msg.Value,
		})
	}
	data, err := types.SerializeCosmosTx(cwicak.Codec, msgs)
	if err != nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "failed to serialize txs"}
	}

	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
		Memo: submitTx.Memo,
	}

	msgServer := icacontrollerkeeper.NewMsgServerImpl(&f)

	owner := contractAddr.String() + "-" + submitTx.AccountId
	res, err := msgServer.SendTx(sdk.WrapSDKContext(ctx), icacontrollertypes.NewMsgSendTx(owner, submitTx.ConnectionId, submitTx.Timeout, packetData))
	if err != nil {
		return nil, errors.Wrap(err, "submitting txs")
	}

	portID, err := icatypes.NewControllerPortID(owner)
	if err != nil {
		return nil, err
	}

	activeChannelID, found := f.GetOpenActiveChannel(ctx, submitTx.ConnectionId, portID)
	if !found {
		return nil, sdkerrors.Wrapf(icatypes.ErrActiveChannelNotFound, "failed to retrieve active channel on connection %s for port %s", submitTx.ConnectionId, portID)
	}

	cwicak.SetCallbackData(ctx, types.CallbackData{
		PortId:       portID,
		ChannelId:    activeChannelID,
		Sequence:     res.Sequence,
		Contract:     contractAddr.String(),
		ConnectionId: submitTx.ConnectionId,
		AccountId:    submitTx.AccountId,
		Callback:     submitTx.Callback,
	})
	return res, nil
}

func HandleMsg(ctx sdk.Context, cwicak cwicakeeper.Keeper, icak icacontrollerkeeper.Keeper, contractAddr sdk.AccAddress, msg *ICAMsg) ([]sdk.Event, [][]byte, error) {
	if msg.Register != nil {
		return register(ctx, contractAddr, msg.Register, cwicak, icak)
	}
	if msg.Submit != nil {
		return submit(ctx, contractAddr, msg.Submit, cwicak, icak)
	}
	return nil, nil, wasmvmtypes.InvalidRequest{Err: "unknown ICA Message variant"}
}
