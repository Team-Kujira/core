package v1

import (
	oracletypes "github.com/Team-Kujira/core/x/oracle/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

func MigrateParams(
	ctx sdk.Context,
	store sdk.KVStore,
	subspace paramtypes.Subspace,
	cdc codec.BinaryCodec,
) error {
	var (
		votePeriod        uint64
		voteThreshold     sdk.Dec
		rewardBand        sdk.Dec
		whitelist         oracletypes.DenomList
		slashFraction     sdk.Dec
		slashWindow       uint64
		minValidPerWindow sdk.Dec
	)

	subspace.Get(ctx, []byte("VotePeriod"), &votePeriod)
	subspace.Get(ctx, []byte("VoteThreshold"), &voteThreshold)
	subspace.Get(ctx, []byte("RewardBand"), &rewardBand)
	subspace.Get(ctx, []byte("Whitelist"), &whitelist)
	subspace.Get(ctx, []byte("SlashFraction"), &slashFraction)
	subspace.Get(ctx, []byte("SlashWindow"), &slashWindow)
	subspace.Get(ctx, []byte("MinValidPerWindow"), &minValidPerWindow)

	denoms := []string{}
	for _, denom := range whitelist {
		denoms = append(denoms, denom.Name)
	}

	oracleParams := oracletypes.Params{
		VotePeriod:        votePeriod,
		VoteThreshold:     voteThreshold,
		MaxDeviation:      rewardBand,
		RequiredDenoms:    denoms,
		SlashFraction:     slashFraction,
		SlashWindow:       slashWindow,
		MinValidPerWindow: minValidPerWindow,
	}

	bz := cdc.MustMarshal(&oracleParams)
	store.Set(oracletypes.ParamsKey, bz)

	return nil
}