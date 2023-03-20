package api

import (
	"math/big"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"go.uber.org/zap"

	"github.com/nnkken/oracle-fetch/db"
	"github.com/nnkken/oracle-fetch/runner"
	"github.com/nnkken/oracle-fetch/types"
)

var testData = []types.DBEntry{
	{
		Token:          "BTC",
		Unit:           "USD",
		Price:          big.NewInt(12000e8),
		PriceTimestamp: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		FetchTimestamp: time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC),
	},
	{
		Token:          "BTC",
		Unit:           "USD",
		Price:          big.NewInt(12000e8),
		PriceTimestamp: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		FetchTimestamp: time.Date(2000, 1, 1, 0, 0, 2, 0, time.UTC),
	},
	{
		Token:          "ETH",
		Unit:           "USD",
		Price:          big.NewInt(2000e8),
		PriceTimestamp: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		FetchTimestamp: time.Date(2000, 1, 1, 0, 0, 3, 0, time.UTC),
	},
	{
		Token:          "ETH",
		Unit:           "USD",
		Price:          big.NewInt(2000e8),
		PriceTimestamp: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		FetchTimestamp: time.Date(2000, 1, 1, 0, 0, 4, 0, time.UTC),
	},
	{
		Token:          "BTC",
		Unit:           "USD",
		Price:          big.NewInt(12500e8),
		PriceTimestamp: time.Date(2000, 1, 1, 0, 0, 4, 0, time.UTC),
		FetchTimestamp: time.Date(2000, 1, 1, 0, 0, 5, 0, time.UTC),
	},
	{
		Token:          "BTC",
		Unit:           "USD",
		Price:          big.NewInt(12500e8),
		PriceTimestamp: time.Date(2000, 1, 1, 0, 0, 4, 0, time.UTC),
		FetchTimestamp: time.Date(2000, 1, 1, 0, 0, 6, 0, time.UTC),
	},
	{
		Token:          "ETH",
		Unit:           "USD",
		Price:          big.NewInt(2100e8),
		PriceTimestamp: time.Date(2000, 1, 1, 0, 0, 7, 0, time.UTC),
		FetchTimestamp: time.Date(2000, 1, 1, 0, 0, 7, 0, time.UTC),
	},
	{
		Token:          "ETH",
		Unit:           "USD",
		Price:          big.NewInt(2100e8),
		PriceTimestamp: time.Date(2000, 1, 1, 0, 0, 7, 0, time.UTC),
		FetchTimestamp: time.Date(2000, 1, 1, 0, 0, 8, 0, time.UTC),
	},
	{
		Token:          "BTC",
		Unit:           "USD",
		Price:          big.NewInt(12700e8),
		PriceTimestamp: time.Date(2000, 1, 1, 0, 0, 8, 0, time.UTC),
		FetchTimestamp: time.Date(2000, 1, 1, 0, 0, 9, 0, time.UTC),
	},
	{
		Token:          "BTC",
		Unit:           "USD",
		Price:          big.NewInt(12700e8),
		PriceTimestamp: time.Date(2000, 1, 1, 0, 0, 8, 0, time.UTC),
		FetchTimestamp: time.Date(2000, 1, 1, 0, 0, 10, 0, time.UTC),
	},
}

func setupTestData(t *testing.T) *pgx.Conn {
	conn := db.SetupTestConn(t)
	for _, d := range testData {
		runner.Insert(d, conn, zap.S())
	}
	return conn
}
