package runner

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/stretchr/testify/require"

	"github.com/nnkken/oracle-fetch/datasource/types"
	"github.com/nnkken/oracle-fetch/db"
)

var testEntries = []db.DBEntry{
	{
		Token:          "TEST",
		Unit:           "TEST-UNIT",
		Price:          big.NewInt(1234),
		PriceTimestamp: time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC),
		FetchTimestamp: time.Date(2010, 1, 1, 0, 0, 1, 0, time.UTC),
	},
	{
		Token:          "TEST-2",
		Unit:           "TEST-2-UNIT",
		Price:          big.NewInt(5678),
		PriceTimestamp: time.Date(2010, 1, 1, 0, 0, 1, 0, time.UTC),
		FetchTimestamp: time.Date(2010, 1, 1, 0, 0, 2, 0, time.UTC),
	},
}

func TestInsert(t *testing.T) {
	conn := db.SetupTestConn()
	defer conn.Close(context.Background())

	Insert(testEntries[0], conn, zap.S())
	Insert(testEntries[1], conn, zap.S())

	rows, err := conn.Query(context.Background(), `
		SELECT token, unit, price, price_timestamp, fetch_timestamp
		FROM prices
		WHERE fetch_timestamp >= $1
		ORDER BY id
	`, testEntries[0].FetchTimestamp)
	require.NoError(t, err)
	hasNextRow := rows.Next()
	require.True(t, hasNextRow)
	var rowEntry db.DBEntry
	var priceStr string
	err = rows.Scan(&rowEntry.Token, &rowEntry.Unit, &priceStr, &rowEntry.PriceTimestamp, &rowEntry.FetchTimestamp)
	require.NoError(t, err)
	priceFloat, _, err := new(big.Float).Parse(priceStr, 10)
	require.NoError(t, err)
	rowEntry.Price, _ = priceFloat.Int(nil)
	rowEntry.PriceTimestamp = rowEntry.PriceTimestamp.UTC()
	rowEntry.FetchTimestamp = rowEntry.FetchTimestamp.UTC()
	require.Equal(t, testEntries[0], rowEntry)
	hasNextRow = rows.Next()
	require.True(t, hasNextRow)
	err = rows.Scan(&rowEntry.Token, &rowEntry.Unit, &priceStr, &rowEntry.PriceTimestamp, &rowEntry.FetchTimestamp)
	require.NoError(t, err)
	priceFloat, _, err = new(big.Float).Parse(priceStr, 10)
	require.NoError(t, err)
	rowEntry.Price, _ = priceFloat.Int(nil)
	rowEntry.PriceTimestamp = rowEntry.PriceTimestamp.UTC()
	rowEntry.FetchTimestamp = rowEntry.FetchTimestamp.UTC()
	require.Equal(t, testEntries[1], rowEntry)
	hasNextRow = rows.Next()
	require.False(t, hasNextRow)

	_, err = conn.Exec(context.Background(), `
		DELETE FROM prices WHERE fetch_timestamp >= $1
	`, testEntries[0].FetchTimestamp)
	require.NoError(t, err)
}

func TestFetch(t *testing.T) {
	mockDataSource := types.NewMockDataSource(t)
	mockCall := mockDataSource.EXPECT().Fetch().Return(testEntries, nil)

	ch := make(chan db.DBEntry, 100)
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
