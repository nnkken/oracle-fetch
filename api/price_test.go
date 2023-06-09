package api

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/nnkken/oracle-fetch/utils"
	"github.com/stretchr/testify/require"
)

func TestParsePriceRequest(t *testing.T) {
	utils.MockTimeNow(t, time.Unix(1234567890, 0))

	// missing everything
	c, _ := newTestContext("/price")
	_, err := ParsePriceRequest(c)
	require.ErrorContains(t, err, "fail to parse price request")

	// missing token
	c, _ = newTestContext("/price?unit=ETH&timestamp=2000-01-01T00:00:00Z")
	_, err = ParsePriceRequest(c)
	require.ErrorContains(t, err, "fail to parse price request")

	// missing timestamp, fallback to now
	c, _ = newTestContext("/price?unit=ETH&token=BTC")
	req, err := ParsePriceRequest(c)
	require.NoError(t, err)
	require.Equal(t, PriceRequest{
		Token:     "BTC",
		Unit:      "ETH",
		Timestamp: time.Unix(1234567890, 0).UTC(),
	}, req)

	// missing unit, fallback to USD
	c, _ = newTestContext("/price?token=BTC&timestamp=2000-01-01T00:00:00Z")
	req, err = ParsePriceRequest(c)
	require.NoError(t, err)
	require.Equal(t, PriceRequest{
		Token:     "BTC",
		Unit:      "USD",
		Timestamp: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}, req)

	// timezone, `%2B` is `+` in URL
	c, _ = newTestContext("/price?unit=ETH&token=BTC&timestamp=2000-01-01T00:00:00%2B08:00")
	req, err = ParsePriceRequest(c)
	require.NoError(t, err)
	require.Equal(t, PriceRequest{
		Token:     "BTC",
		Unit:      "ETH",
		Timestamp: time.Date(1999, 12, 31, 16, 0, 0, 0, time.UTC),
	}, req)
}

func TestHandlePriceRequest(t *testing.T) {
	c, writer := newTestContext("/price?unit=USD&token=BTC&timestamp=2000-01-01T00:00:00Z")
	HandlePriceRequest(c)
	require.Equal(t, 404, writer.Code)

	c, writer = newTestContext("/price?unit=USD&token=BTC&timestamp=2000-01-01T00:00:01Z")
	HandlePriceRequest(c)
	require.Equal(t, 200, writer.Code)
	var res PriceResponse
	err := json.Unmarshal(writer.Written, &res)
	require.NoError(t, err)
	require.Equal(t, "BTC", res.Token)
	require.Equal(t, "USD", res.Unit)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), res.PriceTimestamp)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC), res.FetchTimestamp)
	priceFloat, err := strconv.ParseFloat(res.Price, 64)
	require.NoError(t, err)
	require.Equal(t, float64(12000e8), priceFloat)
}

func TestQueryPrice(t *testing.T) {
	res, err := QueryPrice(PriceRequest{
		Token:     "BTC",
		Unit:      "USD",
		Timestamp: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}, testConn)
	require.ErrorIs(t, err, pgx.ErrNoRows)

	res, err = QueryPrice(PriceRequest{
		Token:     "BTC",
		Unit:      "USD",
		Timestamp: time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC),
	}, testConn)
	require.NoError(t, err)
	require.Equal(t, "BTC", res.Token)
	require.Equal(t, "USD", res.Unit)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), res.PriceTimestamp)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 1, 0, time.UTC), res.FetchTimestamp)
	priceFloat, err := strconv.ParseFloat(res.Price, 64)
	require.NoError(t, err)
	require.Equal(t, float64(12000e8), priceFloat)

	res, err = QueryPrice(PriceRequest{
		Token:     "BTC",
		Unit:      "USD",
		Timestamp: time.Date(2000, 1, 1, 0, 0, 2, 0, time.UTC),
	}, testConn)
	require.NoError(t, err)
	require.Equal(t, "BTC", res.Token)
	require.Equal(t, "USD", res.Unit)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), res.PriceTimestamp)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 2, 0, time.UTC), res.FetchTimestamp)
	priceFloat, err = strconv.ParseFloat(res.Price, 64)
	require.NoError(t, err)
	require.Equal(t, float64(12000e8), priceFloat)

	res, err = QueryPrice(PriceRequest{
		Token:     "BTC",
		Unit:      "USD",
		Timestamp: time.Date(2000, 1, 1, 0, 0, 4, 0, time.UTC),
	}, testConn)
	require.NoError(t, err)
	require.Equal(t, "BTC", res.Token)
	require.Equal(t, "USD", res.Unit)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), res.PriceTimestamp)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 2, 0, time.UTC), res.FetchTimestamp)
	priceFloat, err = strconv.ParseFloat(res.Price, 64)
	require.NoError(t, err)
	require.Equal(t, float64(12000e8), priceFloat)

	res, err = QueryPrice(PriceRequest{
		Token:     "BTC",
		Unit:      "USD",
		Timestamp: time.Date(2000, 1, 1, 0, 0, 10, 0, time.UTC),
	}, testConn)
	require.NoError(t, err)
	require.Equal(t, "BTC", res.Token)
	require.Equal(t, "USD", res.Unit)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 8, 0, time.UTC), res.PriceTimestamp)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 10, 0, time.UTC), res.FetchTimestamp)
	priceFloat, err = strconv.ParseFloat(res.Price, 64)
	require.NoError(t, err)
	require.Equal(t, float64(12700e8), priceFloat)

	res, err = QueryPrice(PriceRequest{
		Token:     "BTC",
		Unit:      "USD",
		Timestamp: time.Date(2000, 1, 1, 0, 0, 11, 0, time.UTC),
	}, testConn)
	require.NoError(t, err)
	require.Equal(t, "BTC", res.Token)
	require.Equal(t, "USD", res.Unit)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 8, 0, time.UTC), res.PriceTimestamp)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 10, 0, time.UTC), res.FetchTimestamp)
	priceFloat, err = strconv.ParseFloat(res.Price, 64)
	require.NoError(t, err)
	require.Equal(t, float64(12700e8), priceFloat)

	res, err = QueryPrice(PriceRequest{
		Token:     "ETH",
		Unit:      "USD",
		Timestamp: time.Date(2000, 1, 1, 0, 0, 4, 0, time.UTC),
	}, testConn)
	require.NoError(t, err)
	require.Equal(t, "ETH", res.Token)
	require.Equal(t, "USD", res.Unit)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), res.PriceTimestamp)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 4, 0, time.UTC), res.FetchTimestamp)
	priceFloat, err = strconv.ParseFloat(res.Price, 64)
	require.NoError(t, err)
	require.Equal(t, float64(2000e8), priceFloat)

	_, err = QueryPrice(PriceRequest{
		Token:     "USDT",
		Unit:      "USD",
		Timestamp: time.Date(2000, 1, 1, 0, 0, 4, 0, time.UTC),
	}, testConn)
	require.ErrorIs(t, err, pgx.ErrNoRows)

	_, err = QueryPrice(PriceRequest{
		Token:     "BTC",
		Unit:      "USDT",
		Timestamp: time.Date(2000, 1, 1, 0, 0, 4, 0, time.UTC),
	}, testConn)
	require.ErrorIs(t, err, pgx.ErrNoRows)

	res, err = QueryPrice(PriceRequest{
		Token:     "BTC",
		Unit:      "USD",
		Timestamp: time.Date(2000, 1, 1, 8, 0, 10, 0, time.FixedZone("HKT", 8*60*60)),
	}, testConn)
	require.NoError(t, err)
	require.Equal(t, "BTC", res.Token)
	require.Equal(t, "USD", res.Unit)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 8, 0, time.UTC), res.PriceTimestamp)
	require.Equal(t, time.Date(2000, 1, 1, 0, 0, 10, 0, time.UTC), res.FetchTimestamp)
	priceFloat, err = strconv.ParseFloat(res.Price, 64)
	require.NoError(t, err)
	require.Equal(t, float64(12700e8), priceFloat)
}
