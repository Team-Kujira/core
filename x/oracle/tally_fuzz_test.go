package oracle_test

import (
	"sort"
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"kujira/x/oracle"
	"kujira/x/oracle/types"
)

func TestFuzz_Tally(t *testing.T) {
	validators := map[string]int64{}

	f := fuzz.New().NilChance(0).Funcs(
		func(e *sdk.Dec, c fuzz.Continue) {
			*e = sdk.NewDec(c.Int63())
		},
		func(e *map[string]int64, c fuzz.Continue) {
			numValidators := c.Intn(100) + 5

			for i := 0; i < numValidators; i++ {
				(*e)[sdk.ValAddress(secp256k1.GenPrivKey().PubKey().Address()).String()] = c.Int63n(100)
			}
		},
		func(e *map[string]types.Claim, c fuzz.Continue) {
			for validator, power := range validators {
				addr, err := sdk.ValAddressFromBech32(validator)
				require.NoError(t, err)
				(*e)[validator] = types.NewClaim(power, 0, 0, addr)
			}
		},
		func(e *types.ExchangeRateBallot, c fuzz.Continue) {

			ballot := types.ExchangeRateBallot{}
			for addr, power := range validators {
				addr, _ := sdk.ValAddressFromBech32(addr)

				var rate sdk.Dec
				c.Fuzz(&rate)

				ballot = append(ballot, types.NewVoteForTally(rate, c.RandString(), addr, power))
			}

			sort.Sort(ballot)

			*e = ballot
		},
	)

	// set random denoms and validators
	f.Fuzz(&validators)

	input, _ := setup(t)

	claimMap := map[string]types.Claim{}
	f.Fuzz(&claimMap)

	ballot := types.ExchangeRateBallot{}
	f.Fuzz(&ballot)

	var rewardBand sdk.Dec
	f.Fuzz(&rewardBand)

	require.NotPanics(t, func() {
		oracle.Tally(input.Ctx, ballot, rewardBand, claimMap)
	})
}
