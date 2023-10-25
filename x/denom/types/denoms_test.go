package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Team-Kujira/core/x/denom/types"
)

func TestDecomposeDenoms(t *testing.T) {
	for _, tc := range []struct {
		desc  string
		denom string
		valid bool
	}{
		{
			desc:  "empty is invalid",
			denom: "",
			valid: false,
		},
		{
			desc:  "normal",
			denom: "factory/kujira1t7egva48prqmzl59x5ngv4zx0dtrwewcug02wd/bitcoin",
			valid: true,
		},
		{
			desc:  "multiple slashes in nonce",
			denom: "factory/kujira1t7egva48prqmzl59x5ngv4zx0dtrwewcug02wd/bitcoin/1",
			valid: true,
		},
		{
			desc:  "no nonce",
			denom: "factory/kujira1t7egva48prqmzl59x5ngv4zx0dtrwewcug02wd/",
			valid: true,
		},
		{
			desc:  "incorrect prefix",
			denom: "ibc/kujira1t7egva48prqmzl59x5ngv4zx0dtrwewcug02wd/bitcoin",
			valid: false,
		},
		{
			desc:  "nonce of only slashes",
			denom: "factory/kujira1t7egva48prqmzl59x5ngv4zx0dtrwewcug02wd/////",
			valid: true,
		},
		{
			desc:  "too long name",
			denom: "factory/kujira1t7egva48prqmzl59x5ngv4zx0dtrwewcug02wd/adsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsf",
			valid: false,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			_, _, err := types.DeconstructDenom(tc.denom)
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
