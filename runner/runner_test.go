package runner

import (
	"errors"
	"math/big"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/nnkken/oracle-fetch/types"
	"github.com/stretchr/testify/require"
)

func TestInsert(t *testing.T) {

}

func TestFetch(t *testing.T) {
	mockDataSource := types.NewMockDataSource(t)
	expected := []types.DBEntry{
		{
			Token:          "TEST",
			Unit:           "TEST-UNIT",
			Price:          big.NewInt(1234),
			PriceTimestamp: time.Unix(1234567890, 0),
			FetchTimestamp: time.Unix(1300000000, 0),
		},
		{
			Token:          "TEST-2",
			Unit:           "TEST-2-UNIT",
			Price:          big.NewInt(5678),
			PriceTimestamp: time.Unix(1234567891, 0),
			FetchTimestamp: time.Unix(1300000001, 0),
		},
	}
	mockCall := mockDataSource.EXPECT().Fetch().Return(expected, nil)

	ch := make(chan types.DBEntry, 100)
	fetch(mockDataSource, ch, zap.S())
	require.Len(t, ch, 2)
	entry := <-ch
	require.Equal(t, expected[0], entry)
	entry = <-ch
	require.Equal(t, expected[1], entry)
	mockCall.Unset()

	mockDataSource.EXPECT().Fetch().Return(nil, errors.New("err"))
	require.NotPanics(t, func() {
		fetch(mockDataSource, ch, zap.S())
	})
	require.Len(t, ch, 0)
	mockCall.Unset()
}
