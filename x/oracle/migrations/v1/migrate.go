package v1

import (
	"cosmossdk.io/math"
	"cosmossdk.io/store"
	oracletypes "github.com/Team-Kujira/core/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

func MigrateParams(
	ctx sdk.Context,
	store store.KVStore,
	subspace paramtypes.Subspace,
	cdc codec.BinaryCodec,
) error {
	var (
		voteThreshold     math.LegacyDec
		rewardBand        math.LegacyDec
		whitelist         oracletypes.DenomList
		slashFraction     math.LegacyDec
		slashWindow       uint64
		minValidPerWindow math.LegacyDec
	)

	subspace.Get(ctx, []byte("VoteThreshold"), &voteThreshold)
	subspace.Get(ctx, []byte("RewardBand"), &rewardBand)
	subspace.Get(ctx, []byte("Whitelist"), &whitelist)
	subspace.Get(ctx, []byte("SlashFraction"), &slashFraction)
	subspace.Get(ctx, []byte("SlashWindow"), &slashWindow)
	subspace.Get(ctx, []byte("MinValidPerWindow"), &minValidPerWindow)

	symbols := []oracletypes.Symbol{}
	for id, denom := range whitelist {
		symbols = append(symbols, oracletypes.Symbol{
			Id:     uint32(id + 1),
			Symbol: denom.Name,
		})
	}

	oracleParams := oracletypes.Params{
		VoteThreshold:     voteThreshold,
		MaxDeviation:      rewardBand,
		RequiredSymbols:   symbols,
		LastSymbolId:      uint32(len(symbols)),
		SlashFraction:     slashFraction,
		SlashWindow:       slashWindow,
		MinValidPerWindow: minValidPerWindow,
	}

	bz := cdc.MustMarshal(&oracleParams)
	store.Set(oracletypes.ParamsKey, bz)

	return nil
}
