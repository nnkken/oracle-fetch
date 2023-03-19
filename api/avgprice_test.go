package api

import (
	"testing"
	"time"

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
	// TODO
}
