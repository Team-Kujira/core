package abci

import (
	"math/big"

	"cosmossdk.io/math"
	"github.com/Team-Kujira/core/x/oracle/keeper"
	"github.com/Team-Kujira/core/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CompressDecimal(dec math.LegacyDec, desiredValidDigit int64) []byte {
	validDigits := int64(0)
	newAmount := dec
	for !newAmount.IsZero() {
		newAmount = newAmount.QuoInt64(10)
		validDigits++
	}
	cuttingDigits := validDigits - desiredValidDigit
	if cuttingDigits < 0 {
		cuttingDigits = 0
	}
	cutAmount := dec.Quo(math.LegacyNewDec(10).Power(uint64(cuttingDigits)))
	return append(cutAmount.BigInt().Bytes(), byte(cuttingDigits))
}

func DecompressDecimal(bz []byte) math.LegacyDec {
	if len(bz) == 0 {
		return math.LegacyZeroDec()
	}
	cuttingDigits := int(bz[len(bz)-1])
	amountInt := new(big.Int).SetBytes(bz[:len(bz)-1])
	amountDec := math.LegacyNewDecFromBigIntWithPrec(amountInt, math.LegacyPrecision)
	amountDec = amountDec.Mul(math.LegacyNewDec(10).Power(uint64(cuttingDigits)))
	return amountDec
}

func ComposeVoteExtension2(k keeper.Keeper, ctx sdk.Context, height int64, exchangeRates sdk.DecCoins) types.VoteExtension2 {
	params := k.GetParams(ctx)
	denomIDs := make(map[string]uint32)
	for _, denom := range params.RequiredDenoms {
		denomIDs[denom.Denom] = denom.Id
	}
	prices := make(map[uint32][]byte)
	for _, rate := range exchangeRates {
		id, ok := denomIDs[rate.Denom]
		if !ok {
			continue
		}
		prices[id] = CompressDecimal(rate.Amount, 8)
	}
	return types.VoteExtension2{
		Height: height,
		Prices: prices,
	}
}

func ExchangeRatesFromVoteExtension2(k keeper.Keeper, ctx sdk.Context, voteExt types.VoteExtension2) sdk.DecCoins {
	params := k.GetParams(ctx)
	idToDenom := make(map[uint32]string)
	for _, denom := range params.RequiredDenoms {
		idToDenom[denom.Id] = denom.Denom
	}

	exchangeRates := sdk.DecCoins{}
	for id, priceBytes := range voteExt.Prices {
		denom, ok := idToDenom[id]
		if !ok {
			continue
		}

		exchangeRates = exchangeRates.Add(sdk.NewDecCoinFromDec(denom, DecompressDecimal(priceBytes)))
	}
	return exchangeRates
}
