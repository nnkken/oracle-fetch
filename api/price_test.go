package api

import (
	"testing"
	"time"

	"github.com/nnkken/oracle-fetch/utils"
	"github.com/stretchr/testify/require"
)

func TestParsePriceRequest(t *testing.T) {
	utils.MockTimeNow(t, time.Unix(1234567890, 0))

	// missing everything
	c := newContextWithUrl("/price")
	_, err := ParsePriceRequest(c)
	require.ErrorContains(t, err, "fail to parse price request")

	// missing token
	c = newContextWithUrl("/price?unit=ETH&timestamp=2000-01-01T00:00:00Z")
	_, err = ParsePriceRequest(c)
	require.ErrorContains(t, err, "fail to parse price request")

	// missing timestamp, fallback to now
	c = newContextWithUrl("/price?unit=ETH&token=BTC")
	req, err := ParsePriceRequest(c)
	require.NoError(t, err)
	require.Equal(t, PriceRequest{
		Token:     "BTC",
		Unit:      "ETH",
		Timestamp: time.Unix(1234567890, 0).UTC(),
	}, req)

	// missing unit, fallback to USD
	c = newContextWithUrl("/price?token=BTC&timestamp=2000-01-01T00:00:00Z")
	req, err = ParsePriceRequest(c)
	require.NoError(t, err)
	require.Equal(t, PriceRequest{
		Token:     "BTC",
		Unit:      "USD",
		Timestamp: time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
	}, req)

	// timezone, `%2B` is `+` in URL
	c = newContextWithUrl("/price?unit=ETH&token=BTC&timestamp=2000-01-01T00:00:00%2B08:00")
	req, err = ParsePriceRequest(c)
	require.NoError(t, err)
	require.Equal(t, PriceRequest{
		Token:     "BTC",
		Unit:      "ETH",
		Timestamp: time.Date(1999, 12, 31, 16, 0, 0, 0, time.UTC),
	}, req)
}

func TestHandlePriceRequest(t *testing.T) {
	// TODO
}

func TestQueryPrice(t *testing.T) {
	// TODO
}
