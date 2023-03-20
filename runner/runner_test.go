package runner

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/stretchr/testify/require"

	"github.com/nnkken/oracle-fetch/db"
	"github.com/nnkken/oracle-fetch/types"
)

var testEntries = []types.DBEntry{
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

func TestInsert(t *testing.T) {
	conn := db.SetupTestConn(t)
	defer conn.Close(context.Background())

	insert(testEntries[0], conn, zap.S())
	insert(testEntries[1], conn, zap.S())

	rows, err := conn.Query(context.Background(), "SELECT token, unit, price, price_timestamp, fetch_timestamp FROM prices ORDER BY id")
	require.NoError(t, err)
	hasNextRow := rows.Next()
	require.True(t, hasNextRow)
	var rowEntry types.DBEntry
	var priceStr string
	err = rows.Scan(&rowEntry.Token, &rowEntry.Unit, &priceStr, &rowEntry.PriceTimestamp, &rowEntry.FetchTimestamp)
	require.NoError(t, err)
	priceFloat, _, err := new(big.Float).Parse(priceStr, 10)
	require.NoError(t, err)
	rowEntry.Price, _ = priceFloat.Int(nil)
	require.Equal(t, testEntries[0], rowEntry)
	hasNextRow = rows.Next()
	require.True(t, hasNextRow)
	err = rows.Scan(&rowEntry.Token, &rowEntry.Unit, &priceStr, &rowEntry.PriceTimestamp, &rowEntry.FetchTimestamp)
	require.NoError(t, err)
	priceFloat, _, err = new(big.Float).Parse(priceStr, 10)
	require.NoError(t, err)
	rowEntry.Price, _ = priceFloat.Int(nil)
	require.Equal(t, testEntries[1], rowEntry)
	hasNextRow = rows.Next()
	require.False(t, hasNextRow)
}

func TestFetch(t *testing.T) {
	mockDataSource := types.NewMockDataSource(t)
	mockCall := mockDataSource.EXPECT().Fetch().Return(testEntries, nil)

	ch := make(chan types.DBEntry, 100)
	fetch(mockDataSource, ch, zap.S())
	require.Len(t, ch, 2)
	entry := <-ch
	require.Equal(t, testEntries[0], entry)
	entry = <-ch
	require.Equal(t, testEntries[1], entry)
	mockCall.Unset()

	mockDataSource.EXPECT().Fetch().Return(nil, errors.New("err"))
	require.NotPanics(t, func() {
		fetch(mockDataSource, ch, zap.S())
	})
	require.Len(t, ch, 0)
	mockCall.Unset()
}
