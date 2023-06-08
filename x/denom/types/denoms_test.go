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
			denom: "factory/cosmos1ft6e5esdtdegnvcr3djd3ftk4kwpcr6jta8eyh/bitcoin",
			valid: true,
		},
		{
			desc:  "multiple slashes in nonce",
			denom: "factory/cosmos1ft6e5esdtdegnvcr3djd3ftk4kwpcr6jta8eyh/bitcoin/1",
			valid: true,
		},
		{
			desc:  "no nonce",
			denom: "factory/cosmos1ft6e5esdtdegnvcr3djd3ftk4kwpcr6jta8eyh/",
			valid: true,
		},
		{
			desc:  "incorrect prefix",
			denom: "ibc/cosmos1ft6e5esdtdegnvcr3djd3ftk4kwpcr6jta8eyh/bitcoin",
			valid: false,
		},
		{
			desc:  "nonce of only slashes",
			denom: "factory/cosmos1ft6e5esdtdegnvcr3djd3ftk4kwpcr6jta8eyh/////",
			valid: true,
		},
		{
			desc:  "too long name",
			denom: "factory/cosmos1ft6e5esdtdegnvcr3djd3ftk4kwpcr6jta8eyh/adsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsfadsf",
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
