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

func (ms msgServer) AddRequiredDenoms(goCtx context.Context, msg *types.MsgAddRequiredDenoms) (*types.MsgAddRequiredDenomsResponse, error) {
	if ms.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", ms.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	params := ms.GetParams(ctx)

	existingSymbols := make(map[string]bool)
	for _, denom := range params.RequiredDenoms {
		existingSymbols[denom.Denom] = true
	}

	for _, denom := range msg.Symbols {
		if existingSymbols[denom] {
			return nil, fmt.Errorf("symbol '%s' already set as required denoms", denom)
		}
		params.RequiredDenoms = append(params.RequiredDenoms, types.Denom{
			Denom: denom,
			Id:    params.LastDenomId + 1,
		})
		params.LastDenomId++
	}

	err := ms.SetParams(ctx, params)
	if err != nil {
		return nil, types.ErrSetParams
	}

	return &types.MsgAddRequiredDenomsResponse{}, nil
}

func (ms msgServer) RemoveRequiredDenoms(goCtx context.Context, msg *types.MsgRemoveRequiredDenoms) (*types.MsgRemoveRequiredDenomsResponse, error) {
	if ms.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", ms.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	params := ms.GetParams(ctx)

	removingSymbols := make(map[string]bool)
	for _, symbol := range msg.Symbols {
		removingSymbols[symbol] = true
	}

	requiredDenoms := []types.Denom{}
	for _, denom := range params.RequiredDenoms {
		if !removingSymbols[denom.Denom] {
			requiredDenoms = append(requiredDenoms, denom)
		}
	}

	params.RequiredDenoms = requiredDenoms
	err := ms.SetParams(ctx, params)
	if err != nil {
		return nil, types.ErrSetParams
	}

	return &types.MsgRemoveRequiredDenomsResponse{}, nil
}

func (ms msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if ms.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", ms.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check id and denom mapping change
	params := ms.GetParams(ctx)
	if !reflect.DeepEqual(params.RequiredDenoms, msg.Params.RequiredDenoms) {
		return nil, types.ErrCanNotUpdateRequiredDenoms
	}

	if err := ms.SetParams(ctx, *msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
