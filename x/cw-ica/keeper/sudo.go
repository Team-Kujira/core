package keeper

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Team-Kujira/core/x/cw-ica/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
)

func (k Keeper) HasContractInfo(ctx sdk.Context, contractAddress sdk.AccAddress) bool {
	return k.wasmKeeper.HasContractInfo(ctx, contractAddress)
}

func (k Keeper) SudoIcaRegisterCallback(
	ctx sdk.Context,
	callbackData types.CallbackData,
	result types.IcaCallbackResult,
) ([]byte, error) {
	contractAddr := sdk.MustAccAddressFromBech32(callbackData.Contract)

	if !k.wasmKeeper.HasContractInfo(ctx, contractAddr) {
		if callbackData.PortId == transfertypes.PortID {
			// we want to allow non contract account to send the assets via IBC Transfer module
			// we can determine the originating module by the source port of the request packet
			return nil, nil
		}
		k.Logger(ctx).Debug("SudoCallback: contract not found", "senderAddress", contractAddr)
		return nil, fmt.Errorf("%s is not a contract address and not the Transfer module", contractAddr)
	}

	x := types.MessageRegisterCallback{}
	x.IcaRegisterCallback.ConnID = callbackData.ConnectionId
	x.IcaRegisterCallback.AccID = callbackData.AccountId
	x.IcaRegisterCallback.Callback = callbackData.Callback
	x.IcaRegisterCallback.Result = result

	m, err := json.Marshal(x)
	if err != nil {
		k.Logger(ctx).Error("SudoCallback: failed to marshal MessageResponse message", "error", err, "contractAddress", contractAddr)
		return nil, fmt.Errorf("failed to marshal MessageResponse: %v", err)
	}

	resp, err := k.wasmKeeper.Sudo(ctx, contractAddr, m)
	if err != nil {
		k.Logger(ctx).Debug("SudoResponse: failed to Sudo", "error", err, "contractAddress", contractAddr)
		return nil, fmt.Errorf("failed to Sudo: %v", err)
	}

	return resp, nil
}

func (k Keeper) SudoIcaTxCallback(
	ctx sdk.Context,
	callbackData types.CallbackData,
	result types.IcaCallbackResult,
) ([]byte, error) {
	contractAddr := sdk.MustAccAddressFromBech32(callbackData.Contract)

	if !k.wasmKeeper.HasContractInfo(ctx, contractAddr) {
		if callbackData.PortId == transfertypes.PortID {
			// we want to allow non contract account to send the assets via IBC Transfer module
			// we can determine the originating module by the source port of the request packet
			return nil, nil
		}
		k.Logger(ctx).Debug("SudoCallback: contract not found", "senderAddress", contractAddr)
		return nil, fmt.Errorf("%s is not a contract address and not the Transfer module", contractAddr)
	}

	x := types.MessageTxCallback{}
	x.IcaTxCallback.ConnID = callbackData.ConnectionId
	x.IcaTxCallback.AccID = callbackData.AccountId
	x.IcaTxCallback.Sequence = callbackData.Sequence
	x.IcaTxCallback.Callback = callbackData.Callback
	x.IcaTxCallback.Result = result

	m, err := json.Marshal(x)
	if err != nil {
		k.Logger(ctx).Error("SudoCallback: failed to marshal MessageResponse message", "error", err, "contractAddress", contractAddr)
		return nil, fmt.Errorf("failed to marshal MessageResponse: %v", err)
	}

	resp, err := k.wasmKeeper.Sudo(ctx, contractAddr, m)
	if err != nil {
		k.Logger(ctx).Debug("SudoResponse: failed to Sudo", "error", err, "contractAddress", contractAddr)
		return nil, fmt.Errorf("failed to Sudo: %v", err)
	}

	return resp, nil
}

func (k Keeper) SudoIbcTransferCallback(
	ctx sdk.Context,
	packet channeltypes.Packet,
	data transfertypes.FungibleTokenPacketData,
	result types.IcaCallbackResult,
) error {
	contractAddr, err := sdk.AccAddressFromBech32(data.Sender)
	if err != nil {
		return err
	}

	if !k.wasmKeeper.HasContractInfo(ctx, contractAddr) {
		return nil
	}

	x := types.MessageTransferCallback{}
	x.TransferCallback.Port = packet.SourcePort
	x.TransferCallback.Channel = packet.SourceChannel
	x.TransferCallback.Sequence = packet.Sequence
	x.TransferCallback.Receiver = data.Receiver
	x.TransferCallback.Denom = data.Denom
	x.TransferCallback.Amount = data.Amount
	x.TransferCallback.Memo = data.Memo
	x.TransferCallback.Result = result

	m, err := json.Marshal(x)
	if err != nil {
		k.Logger(ctx).Error("SudoCallback: failed to marshal MessageResponse message", "error", err, "contractAddress", contractAddr)
		return err
	}

	_, err = k.wasmKeeper.Sudo(ctx, contractAddr, m)
	if err != nil {
		k.Logger(ctx).Debug("SudoResponse: failed to Sudo", "error", err, "contractAddress", contractAddr)
	}
	return err
}

func (k Keeper) SudoIbcTransferReceipt(
	ctx sdk.Context,
	packet channeltypes.Packet,
	data transfertypes.FungibleTokenPacketData,
) error {
	contractAddr, err := sdk.AccAddressFromBech32(data.Receiver)
	if err != nil {
		return err
	}

	if !k.wasmKeeper.HasContractInfo(ctx, contractAddr) {
		return nil
	}

	denom := ""
	if transfertypes.ReceiverChainIsSource(packet.GetSourcePort(), packet.GetSourceChannel(), data.Denom) {
		voucherPrefix := transfertypes.GetDenomPrefix(packet.GetSourcePort(), packet.GetSourceChannel())
		denom = data.Denom[len(voucherPrefix):]
	} else {
		sourcePrefix := transfertypes.GetDenomPrefix(packet.GetDestPort(), packet.GetDestChannel())
		prefixedDenom := sourcePrefix + data.Denom
		denom = transfertypes.ParseDenomTrace(prefixedDenom).IBCDenom()
	}

	x := types.MessageTransferReceipt{}
	x.TransferReceipt.Port = packet.DestinationPort
	x.TransferReceipt.Channel = packet.DestinationChannel
	x.TransferReceipt.Sequence = packet.Sequence
	x.TransferReceipt.Sender = data.Sender
	x.TransferReceipt.Denom = denom
	x.TransferReceipt.Amount = data.Amount
	x.TransferReceipt.Memo = data.Memo

	m, err := json.Marshal(x)
	if err != nil {
		k.Logger(ctx).Error("SudoCallback: failed to marshal MessageResponse message", "error", err, "contractAddress", contractAddr)
		return err
	}

	_, err = k.wasmKeeper.Sudo(ctx, contractAddr, m)
	if err != nil {
		k.Logger(ctx).Debug("SudoResponse: failed to Sudo", "error", err, "contractAddress", contractAddr)
	}
	return err
}
