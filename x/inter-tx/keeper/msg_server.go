package keeper

import (
	"context"
	"time"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
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
func (k msgServer) RegisterAccount(goCtx context.Context, msg *types.MsgSubmitTx) (*types.MsgRegisterAccountResponse, error) {
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
		return nil, err
	}
	//TODO: Maybe we want to return the icacontrollertypes.MsgRegisterInterchainAccountResponse instead our own,
	return &types.MsgRegisterAccountResponse{}, nil
}

func (k msgServer) SubmitTx(goCtx context.Context, msg *types.MsgSubmitTx) (*types.MsgSubmitTxResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	portID, err := icatypes.NewControllerPortID(msg.Owner)
	if err != nil {
		return nil, err
	}

	data, err := icatypes.SerializeCosmosTx(k.cdc, []proto.Message{msg.GetTxMsg()})
	if err != nil {
		return nil, err
	}

	packetData := icatypes.InterchainAccountPacketData{
		Type: icatypes.EXECUTE_TX,
		Data: data,
	}

	// timeoutTimestamp set to max value with the unsigned bit shifted to sastisfy hermes timestamp conversion
	// it is the responsibility of the auth module developer to ensure an appropriate timeout timestamp
	timeoutTimestamp := ctx.BlockTime().Add(time.Minute).UnixNano()
	_, err = k.icaControllerKeeper.SendTx(ctx, nil, msg.ConnectionId, portID, packetData, uint64(timeoutTimestamp)) //nolint:staticcheck //
	if err != nil {
		return nil, err
	}

	return &types.MsgSubmitTxResponse{}, nil
}
