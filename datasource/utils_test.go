package datasource

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestComputeDecimalShift(t *testing.T) {
	require.Equal(t, 8, ComputeDecimalShift(0))
	require.Equal(t, 0, ComputeDecimalShift(8))
	require.Equal(t, 1, ComputeDecimalShift(7))
	require.Equal(t, -1, ComputeDecimalShift(9))
	require.Equal(t, -10, ComputeDecimalShift(18))
}

func TestNormalizePrice(t *testing.T) {
	price := big.NewInt(100000000)
	decimals := uint8(8)
	require.Equal(t, int64(100000000), NormalizePrice(price, decimals).Int64())

	price = big.NewInt(10000000000)
	decimals = 10
	require.Equal(t, int64(100000000), NormalizePrice(price, decimals).Int64())

	price = big.NewInt(123456789)
	decimals = 8
	require.Equal(t, int64(123456789), NormalizePrice(price, decimals).Int64())

	price = big.NewInt(123456789)
	decimals = 9
	require.Equal(t, int64(12345678), NormalizePrice(price, decimals).Int64())

	price = big.NewInt(123456789)
	decimals = 0
	require.Equal(t, int64(12345678900000000), NormalizePrice(price, decimals).Int64())
}
