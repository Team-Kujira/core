package keeper

import (
	"context"

	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	cosmostypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"

	"github.com/Team-Kujira/core/x/inter-tx/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

// NewMsgServerImpl creates and returns a new types.MsgServer, fulfilling the intertx Msg service interface
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// RegisterAccount implements the Msg/RegisterAccount interface
func (k msgServer) RegisterAccount(goCtx context.Context, msg *types.MsgRegisterAccount) (*types.MsgRegisterAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	owner := msg.Sender + "-" + msg.AccountId
	msgRegister := icacontrollertypes.NewMsgRegisterInterchainAccount(msg.ConnectionId, owner, msg.Version)

	if err := msgRegister.ValidateBasic(); err != nil {
		return nil, errors.Wrap(err, "failed validating MsgRegisterInterchainAccount")
	}

	icaMsgServer := icacontrollerkeeper.NewMsgServerImpl(&k.Keeper.icaControllerKeeper)

	_, err := icaMsgServer.RegisterInterchainAccount(
		sdk.WrapSDKContext(ctx),
		msgRegister,
	)

	if err != nil {
		return nil, errors.Wrap(err, "registering ICA")
	}
	return &types.MsgRegisterAccountResponse{}, nil
}

func (k msgServer) SubmitTx(goCtx context.Context, msg *types.MsgSubmitTx) (*types.MsgSubmitTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	data, err := SerializeCosmosTx(k.Keeper.Codec, msg.Msgs)
	if err != nil {
		return nil, err
	}

	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
		Memo: msg.Memo,
	}

	msgServer := icacontrollerkeeper.NewMsgServerImpl(&k.Keeper.icaControllerKeeper)

	owner := msg.Sender + "-" + msg.AccountId
	res, err := msgServer.SendTx(sdk.WrapSDKContext(ctx), icacontrollertypes.NewMsgSendTx(owner, msg.ConnectionId, msg.Timeout, packetData))

	if err != nil {
		return nil, errors.Wrap(err, "submitting txs")
	}

	return &types.MsgSubmitTxResponse{Sequence: res.Sequence}, nil
}

// From neutron inter-tx
func SerializeCosmosTx(cdc codec.BinaryCodec, msgs []*cosmostypes.Any) (bz []byte, err error) {
	// only ProtoCodec is supported
	if _, ok := cdc.(*codec.ProtoCodec); !ok {
		return nil, errors.Wrap(err, "only ProtoCodec is supported on host chain")
	}

	cosmosTx := &icatypes.CosmosTx{
		Messages: msgs,
	}

	bz, err = cdc.Marshal(cosmosTx)
	if err != nil {
		return nil, err
	}

	return bz, nil
}
