package api

import (
	"strconv"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestParseAvgPriceRequest(t *testing.T) {
	// missing everything
	c := newContextWithUrl("/avg_price")
	_, err := ParseAvgPriceRequest(c)
	require.ErrorContains(t, err, "fail to parse avg_price request")

	// missing token
	c = newContextWithUrl("/avg_price?unit=ETH&from=2000-01-01T00:00:00Z&to=2000-01-02T00:00:00Z")
	_, err = ParseAvgPriceRequest(c)
	require.ErrorContains(t, err, "fail to parse avg_price request")

	// missing from
	c = newContextWithUrl("/avg_price?unit=ETH&token=BTC&to=2000-01-02T00:00:00Z")
	_, err = ParseAvgPriceRequest(c)
	require.ErrorContains(t, err, "fail to parse avg_price request")

	// missing to
	c = newContextWithUrl("/avg_price?unit=ETH&token=BTC&from=2000-01-01T00:00:00Z")
	_, err = ParseAvgPriceRequest(c)
	require.ErrorContains(t, err, "fail to parse avg_price request")

	// from after to
	c = newContextWithUrl("/avg_price?unit=ETH&token=BTC&from=2000-01-01T00:00:01Z&to=2000-01-01T00:00:00Z")
	_, err = ParseAvgPriceRequest(c)
	require.ErrorContains(t, err, "'from' must be before 'to'")

	// missing unit, fallback to USD
	c = newContextWithUrl("/avg_price?token=BTC&from=2000-01-01T00:00:00Z&to=2000-01-02T00:00:00Z")
	req, err := ParseAvgPriceRequest(c)
	require.NoError(t, err)
	require.Equal(t, AvgPriceRequest{
		Token: "BTC",
		Unit:  "USD",
		From:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		To:    time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
	}, req)

	// custom unit
	c = newContextWithUrl("/avg_price?unit=ETH&token=BTC&from=2000-01-01T00:00:00Z&to=2000-01-02T00:00:00Z")
	req, err = ParseAvgPriceRequest(c)
	require.NoError(t, err)
	require.Equal(t, AvgPriceRequest{
		Token: "BTC",
		Unit:  "ETH",
		From:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		To:    time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
	}, req)

	// timezone, `%2B` is `+` in URL
	c = newContextWithUrl("/avg_price?unit=ETH&token=BTC&from=2000-01-01T00:00:00%2B08:00&to=2000-01-02T00:00:00%2B08:00")
	req, err = ParseAvgPriceRequest(c)
	require.NoError(t, err)
	require.Equal(t, AvgPriceRequest{
		Token: "BTC",
		Unit:  "ETH",
		From:  time.Date(1999, 12, 31, 16, 0, 0, 0, time.UTC),
		To:    time.Date(2000, 1, 1, 16, 0, 0, 0, time.UTC),
	}, req)
}

func TestHandleAvgPriceRequest(t *testing.T) {
	// TODO
}

func TestQueryAvgPrice(t *testing.T) {
	conn := setupTestData(t)
	res, err := QueryAvgPrice(AvgPriceRequest{
		Token: "BTC",
		Unit:  "USD",
		From:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		To:    time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
	}, conn)
	require.NoError(t, err)
	require.Equal(t, "BTC", res.Token)
	require.Equal(t, "USD", res.Unit)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC), res.FirstFetchTimestamp)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 10, 0, time.UTC), res.LastFetchTimestamp)
	require.Equal(t, uint(6), res.PriceCount)
	priceFloat, err := strconv.ParseFloat(res.AvgPrice, 64)
	require.NoError(t, err)
	require.Equal(t, float64(12400e8), priceFloat)

	res, err = QueryAvgPrice(AvgPriceRequest{
		Token: "BTC",
		Unit:  "USD",
		From:  time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC),
		To:    time.Date(2000, 1, 1, 0, 0, 5, 0, time.UTC),
	}, conn)
	require.NoError(t, err)
	require.Equal(t, "BTC", res.Token)
	require.Equal(t, "USD", res.Unit)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC), res.FirstFetchTimestamp)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 5, 0, time.UTC), res.LastFetchTimestamp)
	require.Equal(t, uint(3), res.PriceCount)
	priceFloat, err = strconv.ParseFloat(res.AvgPrice, 64)
	require.NoError(t, err)
	require.Equal(t, float64((12000e8+12000e8+12500e8)/3), priceFloat)

	res, err = QueryAvgPrice(AvgPriceRequest{
		Token: "ETH",
		Unit:  "USD",
		From:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		To:    time.Date(2000, 1, 1, 0, 0, 10, 0, time.UTC),
	}, conn)
	require.NoError(t, err)
	require.Equal(t, "ETH", res.Token)
	require.Equal(t, "USD", res.Unit)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 3, 0, time.UTC), res.FirstFetchTimestamp)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 8, 0, time.UTC), res.LastFetchTimestamp)
	require.Equal(t, uint(4), res.PriceCount)
	priceFloat, err = strconv.ParseFloat(res.AvgPrice, 64)
	require.NoError(t, err)
	require.Equal(t, float64(2050e8), priceFloat)

	_, err = QueryAvgPrice(AvgPriceRequest{
		Token: "BTC",
		Unit:  "HKD",
		From:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		To:    time.Date(2000, 1, 1, 0, 0, 10, 0, time.UTC),
	}, conn)
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = QueryAvgPrice(AvgPriceRequest{
		Token: "USDT",
		Unit:  "USD",
		From:  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		To:    time.Date(2000, 1, 1, 0, 0, 10, 0, time.UTC),
	}, conn)
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = QueryAvgPrice(AvgPriceRequest{
		Token: "BTC",
		Unit:  "USD",
		From:  time.Date(2000, 1, 1, 0, 0, 11, 0, time.UTC),
		To:    time.Date(2000, 1, 1, 0, 0, 12, 0, time.UTC),
	}, conn)
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = QueryAvgPrice(AvgPriceRequest{
		Token: "BTC",
		Unit:  "USD",
		From:  time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC),
		To:    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}, conn)
	require.ErrorIs(t, err, pgx.ErrNoRows)

	res, err = QueryAvgPrice(AvgPriceRequest{
		Token: "BTC",
		Unit:  "USD",
		From:  time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC),
		To:    time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC),
	}, conn)
	require.NoError(t, err)
	require.Equal(t, "BTC", res.Token)
	require.Equal(t, "USD", res.Unit)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC), res.FirstFetchTimestamp)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC), res.LastFetchTimestamp)
	require.Equal(t, uint(1), res.PriceCount)
	priceFloat, err = strconv.ParseFloat(res.AvgPrice, 64)
	require.NoError(t, err)
	require.Equal(t, float64(12000e8), priceFloat)

	res, err = QueryAvgPrice(AvgPriceRequest{
		Token: "BTC",
		Unit:  "USD",
		From:  time.Date(2000, 1, 1, 0, 0, 10, 0, time.UTC),
		To:    time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
	}, conn)
	require.NoError(t, err)
	require.Equal(t, "BTC", res.Token)
	require.Equal(t, "USD", res.Unit)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 10, 0, time.UTC), res.FirstFetchTimestamp)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 10, 0, time.UTC), res.LastFetchTimestamp)
	require.Equal(t, uint(1), res.PriceCount)
	priceFloat, err = strconv.ParseFloat(res.AvgPrice, 64)
	require.NoError(t, err)
	require.Equal(t, float64(12700e8), priceFloat)

	res, err = QueryAvgPrice(AvgPriceRequest{
		Token: "BTC",
		Unit:  "USD",
		From:  time.Date(2000, 1, 1, 8, 0, 10, 0, time.FixedZone("HKT", 8*60*60)),
		To:    time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
	}, conn)
	require.NoError(t, err)
	require.Equal(t, "BTC", res.Token)
	require.Equal(t, "USD", res.Unit)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 10, 0, time.UTC), res.FirstFetchTimestamp)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 10, 0, time.UTC), res.LastFetchTimestamp)
	require.Equal(t, uint(1), res.PriceCount)
	priceFloat, err = strconv.ParseFloat(res.AvgPrice, 64)
	require.NoError(t, err)
	require.Equal(t, float64(12700e8), priceFloat)

	res, err = QueryAvgPrice(AvgPriceRequest{
		Token: "BTC",
		Unit:  "USD",
		From:  time.Date(2000, 1, 1, 0, 0, 10, 0, time.UTC),
		To:    time.Date(2001, 1, 1, 8, 0, 0, 0, time.FixedZone("HKT", 8*60*60)),
	}, conn)
	require.NoError(t, err)
	require.Equal(t, "BTC", res.Token)
	require.Equal(t, "USD", res.Unit)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 10, 0, time.UTC), res.FirstFetchTimestamp)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 10, 0, time.UTC), res.LastFetchTimestamp)
	require.Equal(t, uint(1), res.PriceCount)
	priceFloat, err = strconv.ParseFloat(res.AvgPrice, 64)
	require.NoError(t, err)
	require.Equal(t, float64(12700e8), priceFloat)

	res, err = QueryAvgPrice(AvgPriceRequest{
		Token: "BTC",
		Unit:  "USD",
		From:  time.Date(2000, 1, 1, 8, 0, 10, 0, time.FixedZone("HKT", 8*60*60)),
		To:    time.Date(2001, 1, 1, 8, 0, 0, 0, time.FixedZone("HKT", 8*60*60)),
	}, conn)
	require.NoError(t, err)
	require.Equal(t, "BTC", res.Token)
	require.Equal(t, "USD", res.Unit)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 10, 0, time.UTC), res.FirstFetchTimestamp)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 10, 0, time.UTC), res.LastFetchTimestamp)
	require.Equal(t, uint(1), res.PriceCount)
	priceFloat, err = strconv.ParseFloat(res.AvgPrice, 64)
	require.NoError(t, err)
	require.Equal(t, float64(12700e8), priceFloat)
}
