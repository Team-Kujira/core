package keeper

import (
	"context"
	"fmt"
	"reflect"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/Team-Kujira/core/x/oracle/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the oracle MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (ms msgServer) AddRequiredSymbols(goCtx context.Context, msg *types.MsgAddRequiredSymbols) (*types.MsgAddRequiredSymbolsResponse, error) {
	if ms.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", ms.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	params := ms.GetParams(ctx)

	existingSymbols := make(map[string]bool)
	for _, symbol := range params.RequiredSymbols {
		existingSymbols[symbol.Symbol] = true
	}

	for _, symbol := range msg.Symbols {
		if existingSymbols[symbol] {
			return nil, fmt.Errorf("symbol '%s' already set as required symbols", symbol)
		}
		params.RequiredSymbols = append(params.RequiredSymbols, types.Symbol{
			Symbol: symbol,
			Id:     params.LastSymbolId + 1,
		})
		params.LastSymbolId++
	}

	err := ms.SetParams(ctx, params)
	if err != nil {
		return nil, types.ErrSetParams
	}

	return &types.MsgAddRequiredSymbolsResponse{}, nil
}

func (ms msgServer) RemoveRequiredSymbols(goCtx context.Context, msg *types.MsgRemoveRequiredSymbols) (*types.MsgRemoveRequiredSymbolsResponse, error) {
	if ms.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", ms.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	params := ms.GetParams(ctx)

	removingSymbols := make(map[string]bool)
	for _, symbol := range msg.Symbols {
		removingSymbols[symbol] = true
	}

	requiredDenoms := []types.Symbol{}
	for _, denom := range params.RequiredSymbols {
		if !removingSymbols[denom.Symbol] {
			requiredDenoms = append(requiredDenoms, denom)
		}
	}

	params.RequiredSymbols = requiredDenoms
	err := ms.SetParams(ctx, params)
	if err != nil {
		return nil, types.ErrSetParams
	}

	return &types.MsgRemoveRequiredSymbolsResponse{}, nil
}

func (ms msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if ms.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", ms.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check id and denom mapping change
	params := ms.GetParams(ctx)
	if !reflect.DeepEqual(params.RequiredSymbols, msg.Params.RequiredSymbols) {
		return nil, types.ErrCanNotUpdateRequiredSymbols
	}

	if err := ms.SetParams(ctx, *msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
