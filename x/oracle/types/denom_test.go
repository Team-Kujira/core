package types_test

import (
	"testing"

	"kujira/x/oracle/types"

	"github.com/stretchr/testify/require"
)

func Test_DenomList(t *testing.T) {
	denoms := types.DenomList{
		types.Denom{
			Name: "denom1",
		},
		types.Denom{
			Name: "denom2",
		},
		types.Denom{
			Name: "denom3",
		},
	}

	require.False(t, denoms[0].Equal(&denoms[1]))
	require.True(t, denoms[0].Equal(&denoms[0]))
	require.Equal(t, "name: denom1\n", denoms[0].String())
	require.Equal(t, "name: denom2\n", denoms[1].String())
	require.Equal(t, "name: denom3\n", denoms[2].String())
	require.Equal(t, "name: denom1\n\nname: denom2\n\nname: denom3", denoms.String())
}
