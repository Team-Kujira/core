package keeper

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Team-Kujira/core/x/cw-ica/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
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
		if callbackData.PortId == ibctransfertypes.PortID {
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
		if callbackData.PortId == ibctransfertypes.PortID {
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
