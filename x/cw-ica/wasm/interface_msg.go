package wasm

import (
	"cosmossdk.io/errors"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	cwicakeeper "github.com/Team-Kujira/core/x/cw-ica/keeper"
	"github.com/Team-Kujira/core/x/cw-ica/types"
	cosmostypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
)

// ProtobufAny is a hack-struct to serialize protobuf Any message into JSON object
// See https://github.com/neutron-org/neutron/blob/main/wasmbinding/bindings/msg.go
type ProtobufAny struct {
	TypeURL string `json:"type_url"`
	Value   []byte `json:"value"`
}

type CwIcaMsg struct {
	/// Contracts can register a new interchain account.
	Register *Register `json:"register,omitempty"`
	/// Contracts can submit transactions to the ICA
	/// associated with the given information.
	Submit *Submit `json:"submit,omitempty"`
	/// Transfer submits ibc-transfer msg with
	/// optional callback field
	Transfer *Transfer `json:"transfer,omitempty"`
}

// / Register creates a new interchain account.
// / If the account was created in the past, this will
// / re-establish a dropped connection, or do nothing if
// / the connection is still active.
// / The account is registered using (port, channel, sender, id)
// / as the unique identifier.
type Register struct {
	ConnectionID string `json:"connection_id"`
	AccountID    string `json:"account_id"`
	Version      string `json:"version"`
	Callback     []byte `json:"callback"`
}

// / Submit submits transactions to the ICA
// / associated with the given address.
type Submit struct {
	ConnectionID string        `json:"connection_id"`
	AccountID    string        `json:"account_id"`
	Msgs         []ProtobufAny `json:"msgs"`
	Memo         string        `json:"memo"`
	Timeout      uint64        `json:"timeout"`
	Callback     []byte        `json:"callback"`
}

type Transfer struct {
	ChannelID string                 `json:"channel_id"`
	ToAddress string                 `json:"to_address"`
	Amount    wasmvmtypes.Coin       `json:"amount"`
	Timeout   wasmvmtypes.IBCTimeout `json:"timeout"`
	Callback  []byte                 `json:"callback"`
}

func register(ctx sdk.Context, contractAddr sdk.AccAddress, register *Register, cwicak cwicakeeper.Keeper, ik icacontrollerkeeper.Keeper) ([]sdk.Event, [][]byte, error) {
	_, err := PerformRegisterICA(cwicak, ik, ctx, contractAddr, register)
	if err != nil {
		return nil, nil, errors.Wrap(err, "perform register ICA")
	}
	return nil, nil, nil
}

// PerformRegisterICA is used with register to validate the register message and register the ICA.
func PerformRegisterICA(cwicak cwicakeeper.Keeper, f icacontrollerkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, msg *Register) (*icacontrollertypes.MsgRegisterInterchainAccountResponse, error) {
	if msg == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "register ICA null message"}
	}

	msgServer := icacontrollerkeeper.NewMsgServerImpl(&f)

	// format "{owner}-{id}"
	owner := contractAddr.String() + "-" + msg.AccountID
	msgRegister := icacontrollertypes.NewMsgRegisterInterchainAccount(msg.ConnectionID, owner, msg.Version)

	if err := msgRegister.ValidateBasic(); err != nil {
		return nil, errors.Wrap(err, "failed validating MsgRegisterInterchainAccount")
	}

	res, err := msgServer.RegisterInterchainAccount(
		sdk.WrapSDKContext(ctx),
		msgRegister,
	)
	if err != nil {
		return nil, err
	}

	portID, err := icatypes.NewControllerPortID(owner)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, errors.Wrap(err, "registering ICA")
	}

	f.SetMiddlewareEnabled(ctx, portID, msg.ConnectionID)

	cwicak.SetCallbackData(ctx, types.CallbackData{
		PortId:       portID,
		ChannelId:    "",
		Sequence:     0,
		Contract:     contractAddr.String(),
		ConnectionId: msg.ConnectionID,
		AccountId:    msg.AccountID,
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

	owner := contractAddr.String() + "-" + submitTx.AccountID
	res, err := msgServer.SendTx(sdk.WrapSDKContext(ctx), icacontrollertypes.NewMsgSendTx(owner, submitTx.ConnectionID, submitTx.Timeout, packetData))
	if err != nil {
		return nil, errors.Wrap(err, "submitting txs")
	}

	portID, err := icatypes.NewControllerPortID(owner)
	if err != nil {
		return nil, err
	}

	activeChannelID, found := f.GetOpenActiveChannel(ctx, submitTx.ConnectionID, portID)
	if !found {
		return nil, errors.Wrapf(icatypes.ErrActiveChannelNotFound, "failed to retrieve active channel on connection %s for port %s", submitTx.ConnectionID, portID)
	}

	cwicak.SetCallbackData(ctx, types.CallbackData{
		PortId:       portID,
		ChannelId:    activeChannelID,
		Sequence:     res.Sequence,
		Contract:     contractAddr.String(),
		ConnectionId: submitTx.ConnectionID,
		AccountId:    submitTx.AccountID,
		Callback:     submitTx.Callback,
	})
	return res, nil
}

func transfer(ctx sdk.Context, contractAddr sdk.AccAddress, transferTx *Transfer, cwicak cwicakeeper.Keeper, tk ibctransferkeeper.Keeper) ([]sdk.Event, [][]byte, error) {
	_, err := PerformTransfer(tk, cwicak, ctx, contractAddr, transferTx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "perform submit txs")
	}
	return nil, nil, nil
}

// PerformTransfer is used to perform ibc transfer through wasmbinding.
func PerformTransfer(f ibctransferkeeper.Keeper, cwicak cwicakeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, transferTx *Transfer) (*ibctransfertypes.MsgTransferResponse, error) {
	if transferTx == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "transfer txs null message"}
	}

	amount, err := wasmkeeper.ConvertWasmCoinToSdkCoin(transferTx.Amount)
	if err != nil {
		return nil, errors.Wrap(err, "amount")
	}
	msg := &ibctransfertypes.MsgTransfer{
		SourcePort:       ibctransfertypes.PortID,
		SourceChannel:    transferTx.ChannelID,
		Token:            amount,
		Sender:           contractAddr.String(),
		Receiver:         transferTx.ToAddress,
		TimeoutHeight:    wasmkeeper.ConvertWasmIBCTimeoutHeightToCosmosHeight(transferTx.Timeout.Block),
		TimeoutTimestamp: transferTx.Timeout.Timestamp,
	}

	res, err := f.Transfer(sdk.WrapSDKContext(ctx), msg)
	if err != nil {
		return nil, errors.Wrap(err, "submitting transfer tx")
	}

	cwicak.SetCallbackData(ctx, types.CallbackData{
		PortId:       msg.SourcePort,
		ChannelId:    msg.SourceChannel,
		Sequence:     res.Sequence,
		Contract:     contractAddr.String(),
		ConnectionId: "",
		AccountId:    "",
		Callback:     transferTx.Callback,
	})
	return res, nil
}

func HandleMsg(ctx sdk.Context, cwicak cwicakeeper.Keeper, icak icacontrollerkeeper.Keeper, transferk ibctransferkeeper.Keeper, contractAddr sdk.AccAddress, msg *CwIcaMsg) ([]sdk.Event, [][]byte, error) {
	if msg.Register != nil {
		return register(ctx, contractAddr, msg.Register, cwicak, icak)
	}
	if msg.Submit != nil {
		return submit(ctx, contractAddr, msg.Submit, cwicak, icak)
	}
	if msg.Transfer != nil {
		return transfer(ctx, contractAddr, msg.Transfer, cwicak, transferk)
	}
	return nil, nil, wasmvmtypes.InvalidRequest{Err: "unknown ICA Message variant"}
}
