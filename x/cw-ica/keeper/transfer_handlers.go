package keeper

import (
	"fmt"

	"github.com/Team-Kujira/core/x/cw-ica/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
)

// HandleAcknowledgement passes the acknowledgement data to the appropriate contract via a Sudo call.
func (k *Keeper) HandleTransferAcknowledgement(ctx sdk.Context, packet channeltypes.Packet, acknowledgement []byte, _ sdk.AccAddress) {
	k.Logger(ctx).Debug("Handling transfer acknowledgement")
	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return
	}

	var ack channeltypes.Acknowledgement
	if err := channeltypes.SubModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		k.Logger(ctx).Error("HandleTransferAcknowledgement: cannot unmarshal IBC transfer packet acknowledgement", "error", err)
		return
	}

	cacheCtx, writeFn, newGasMeter := k.createCachedContext(ctx)
	defer k.outOfGasRecovery(ctx, newGasMeter)

	// Actually we have only one kind of error returned from acknowledgement
	// maybe later we'll retrieve actual errors from events
	errorText := ack.GetError()
	var err error
	if errorText != "" {
		err = k.SudoIbcTransferCallback(cacheCtx, packet, data, types.IcaCallbackResult{
			Error: &types.IcaCallbackError{
				Error: errorText,
			},
		})
	} else {
		err = k.SudoIbcTransferCallback(cacheCtx, packet, data, types.IcaCallbackResult{
			Success: &types.IcaCallbackSuccess{
				Data: ack.GetResult(),
			},
		})
	}

	if err != nil {
		k.Logger(ctx).Debug(
			"HandleTransferAcknowledgement: failed to Sudo contract on transfer packet acknowledgement",
			"source_port", packet.SourcePort,
			"source_channel", packet.SourceChannel,
			"sequence", packet.Sequence,
			"error", err)
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeICATxCallbackFailure,
				sdk.NewAttribute(types.AttributePacketSourcePort, packet.SourcePort),
				sdk.NewAttribute(types.AttributePacketSourceChannel, packet.SourceChannel),
				sdk.NewAttribute(types.AttributePacketSequence, fmt.Sprintf("%d", packet.Sequence)),
			),
		})
	} else {
		ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())
		writeFn()
	}

	ctx.GasMeter().ConsumeGas(newGasMeter.GasConsumed(), "consume from cached context")
}

func (k *Keeper) HandleTransferTimeout(ctx sdk.Context, packet channeltypes.Packet, _ sdk.AccAddress) {
	k.Logger(ctx).Debug("Transfer HandleTimeout")

	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return
	}

	cacheCtx, writeFn, newGasMeter := k.createCachedContext(ctx)
	defer k.outOfGasRecovery(ctx, newGasMeter)

	err := k.SudoIbcTransferCallback(ctx, packet, data, types.IcaCallbackResult{
		Timeout: &types.IcaCallbackTimeout{},
	})
	if err != nil {
		k.Logger(ctx).Debug(
			"HandleTransferTimeout: failed to Sudo contract on transfer packet timeout",
			"source_port", packet.SourcePort,
			"source_channel", packet.SourceChannel,
			"sequence", packet.Sequence,
			"error", err)
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeICATimeoutCallbackFailure,
				sdk.NewAttribute(types.AttributePacketSourcePort, packet.SourcePort),
				sdk.NewAttribute(types.AttributePacketSourceChannel, packet.SourceChannel),
				sdk.NewAttribute(types.AttributePacketSequence, fmt.Sprintf("%d", packet.Sequence)),
			),
		})
	} else {
		ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())
		writeFn()
	}

	ctx.GasMeter().ConsumeGas(newGasMeter.GasConsumed(), "consume from cached context")
}

func (k *Keeper) HandleTransferReceipt(ctx sdk.Context, packet channeltypes.Packet, _ sdk.AccAddress) {
	k.Logger(ctx).Debug("Transfer Receipt")

	var data transfertypes.FungibleTokenPacketData
	if err := transfertypes.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return
	}

	cacheCtx, writeFn, newGasMeter := k.createCachedContext(ctx)
	defer k.outOfGasRecovery(ctx, newGasMeter)

	err := k.SudoIbcTransferReceipt(ctx, packet, data)
	if err != nil {
		k.Logger(ctx).Debug(
			"HandleTransferReceipt: failed to Sudo contract on transfer receipt",
			"destination_port", packet.DestinationPort,
			"destination_channel", packet.DestinationChannel,
			"sequence", packet.Sequence,
			"error", err)
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeICATimeoutCallbackFailure,
				sdk.NewAttribute(types.AttributePacketSourcePort, packet.SourcePort),
				sdk.NewAttribute(types.AttributePacketSourceChannel, packet.SourceChannel),
				sdk.NewAttribute(types.AttributePacketSequence, fmt.Sprintf("%d", packet.Sequence)),
			),
		})
	} else {
		ctx.EventManager().EmitEvents(cacheCtx.EventManager().Events())
		writeFn()
	}

	ctx.GasMeter().ConsumeGas(newGasMeter.GasConsumed(), "consume from cached context")
}
