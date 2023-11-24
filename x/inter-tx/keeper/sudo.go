package keeper

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Team-Kujira/core/x/inter-tx/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
)

func (k Keeper) HasContractInfo(ctx sdk.Context, contractAddress sdk.AccAddress) bool {
	return k.wasmKeeper.HasContractInfo(ctx, contractAddress)
}

func (k Keeper) SudoIcaCallback(
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

	x := types.MessageCallback{}
	x.IcaCallback.ConnId = callbackData.ConnectionId
	x.IcaCallback.AccId = callbackData.AccountId
	x.IcaCallback.TxId = callbackData.TxId
	x.IcaCallback.Result = result

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
