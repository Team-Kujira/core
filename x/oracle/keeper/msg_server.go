package keeper

import (
	"context"
	"fmt"

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

func (ms msgServer) AddRequiredDenom(goCtx context.Context, msg *types.MsgAddRequiredDenom) (*types.MsgAddRequiredDenomResponse, error) {
	if ms.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", ms.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	params := ms.GetParams(ctx)

	denoms := params.RequiredDenoms
	for _, denom := range denoms {
		if denom == msg.Symbol {
			return nil, fmt.Errorf("symbol '%s' already set as required denoms", msg.Symbol)
		}
	}

	denoms = append(denoms, msg.Symbol)
	params.RequiredDenoms = denoms
	err := ms.SetParams(ctx, params)
	if err != nil {
		return nil, types.ErrSetParams
	}

	return &types.MsgAddRequiredDenomResponse{}, nil
}

func (ms msgServer) RemoveRequiredDenom(goCtx context.Context, msg *types.MsgRemoveRequiredDenom) (*types.MsgRemoveRequiredDenomResponse, error) {
	if ms.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", ms.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	params := ms.GetParams(ctx)

	denoms := params.RequiredDenoms
	index := -1
	for i, denom := range denoms {
		if denom == msg.Symbol {
			index = i
			break
		}
	}

	if index < 0 {
		return nil, fmt.Errorf("symbol '%s' not found in required denoms", msg.Symbol)
	}

	denoms = append(denoms[:index], denoms[index+1:]...)
	params.RequiredDenoms = denoms
	err := ms.SetParams(ctx, params)
	if err != nil {
		return nil, types.ErrSetParams
	}

	return &types.MsgRemoveRequiredDenomResponse{}, nil
}

func (ms msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if ms.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", ms.authority, msg.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := ms.SetParams(ctx, *msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
