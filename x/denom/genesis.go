package denom

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Team-Kujira/core/x/denom/keeper"
	"github.com/Team-Kujira/core/x/denom/types"
)

// InitGenesis initializes the denom module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.CreateModuleAccount(ctx)

	if genState.Params.CreationFee == nil {
		genState.Params.CreationFee = sdk.NewCoins()
	}
	k.SetParams(ctx, genState.Params)

	for _, genDenom := range genState.GetFactoryDenoms() {
		creator, nonce, err := types.DeconstructDenom(genDenom.GetDenom())
		if err != nil {
			panic(err)
		}
		_, err = k.CreateDenom(ctx, creator, nonce)
		if err != nil {
			panic(err)
		}
		err = k.SetAuthorityMetadata(ctx, genDenom.GetDenom(), genDenom.GetAuthorityMetadata())
		if err != nil {
			panic(err)
		}
	}

	for _, addr := range genState.NoFeeAccounts {
		k.SetNoFeeAccount(ctx, addr)
	}
}

// ExportGenesis returns the denom module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genDenoms := []types.GenesisDenom{}
	iterator := k.GetAllDenomsIterator(ctx)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		denom := string(iterator.Value())

		authorityMetadata, err := k.GetAuthorityMetadata(ctx, denom)
		if err != nil {
			panic(err)
		}

		genDenoms = append(genDenoms, types.GenesisDenom{
			Denom:             denom,
			AuthorityMetadata: authorityMetadata,
		})
	}

	return &types.GenesisState{
		FactoryDenoms: genDenoms,
		Params:        k.GetParams(ctx),
		NoFeeAccounts: k.GetNoFeeAccounts(ctx),
	}
}
