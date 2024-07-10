package abci_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/Team-Kujira/core/x/oracle/abci"
	"github.com/stretchr/testify/require"
)

func TestCompressDecimal(t *testing.T) {
	// compress 0 decimal
	bz := abci.CompressDecimal(math.LegacyZeroDec(), 8)
	require.Len(t, bz, 1)

	// compress negative decimal - only abs value's encoded/decoded
	bz = abci.CompressDecimal(math.LegacyNewDec(-1), 8)
	require.Len(t, bz, 4)

	// compress low number of valid digits decimal
	bz = abci.CompressDecimal(math.LegacyNewDecWithPrec(123, 18), 8)
	require.Len(t, bz, 2)

	// compress high number of valid digits decimal
	bz = abci.CompressDecimal(math.LegacyNewDecWithPrec(123456123456, 18), 8)
	require.Len(t, bz, 4)

	// compress big number
	bz = abci.CompressDecimal(math.LegacyNewDec(123456123456123456), 8)
	require.Len(t, bz, 4)
}

func TestDecompressDecimal(t *testing.T) {
	// empty bytes
	decoded := abci.DecompressDecimal([]byte{})
	require.Equal(t, decoded, math.LegacyZeroDec())

	// decompress 0 decimal bytes
	bz := abci.CompressDecimal(math.LegacyZeroDec(), 8)
	decoded = abci.DecompressDecimal(bz)
	require.Equal(t, decoded, math.LegacyZeroDec())

	// decompress negative decimal bytes - only abs value's encoded/decoded
	bz = abci.CompressDecimal(math.LegacyNewDec(-1), 8)
	decoded = abci.DecompressDecimal(bz)
	require.Equal(t, decoded, math.LegacyNewDec(1))

	// decompress low number of valid digits decimal bytes
	bz = abci.CompressDecimal(math.LegacyNewDecWithPrec(123, 18), 8)
	decoded = abci.DecompressDecimal(bz)
	require.Equal(t, decoded, math.LegacyNewDecWithPrec(123, 18))

	// decompress high number of valid digits decimal bytes
	bz = abci.CompressDecimal(math.LegacyNewDecWithPrec(123456123456, 18), 8)
	decoded = abci.DecompressDecimal(bz)
	require.Equal(t, decoded, math.LegacyNewDecWithPrec(123456120000, 18))

	// decompress big number bytes
	bz = abci.CompressDecimal(math.LegacyNewDec(123456123456123456), 8)
	decoded = abci.DecompressDecimal(bz)
	require.Equal(t, decoded, math.LegacyNewDec(123456120000000000))
}
