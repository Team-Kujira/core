package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"

	"github.com/Team-Kujira/core/x/cw-ica/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

// NewMsgServerImpl creates and returns a new types.MsgServer, fulfilling the cwica Msg service interface
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

	data, err := types.SerializeCosmosTx(k.Keeper.Codec, msg.Msgs)
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
