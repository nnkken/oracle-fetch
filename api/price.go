package api

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/jackc/pgx/v5"

	"github.com/nnkken/oracle-fetch/utils"
)

type PriceRequest struct {
	Token     string    `form:"token" binding:"required"`
	Unit      string    `form:"unit,default=USD"`
	Timestamp time.Time `form:"timestamp"`
}

type PriceResponse struct {
	Token          string    `json:"token" example:"BTC"`
	Unit           string    `json:"unit" example:"USD"`
	Price          string    `json:"price" example:"1234500000000.000000 (8 extra decimal places, so it means 12345)"`
	PriceTimestamp time.Time `json:"price_timestamp" example:"2023-03-18T01:23:45Z"`
	FetchTimestamp time.Time `json:"fetch_timestamp" example:"2023-03-18T01:23:45Z"`
}

func ParsePriceRequest(c *gin.Context) (PriceRequest, error) {
	var req PriceRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		return PriceRequest{}, fmt.Errorf("fail to parse price request: %w", err)
	}
	if req.Timestamp.IsZero() {
		req.Timestamp = utils.TimeNow()
	}
	req.Timestamp = req.Timestamp.UTC()
	return req, nil
}

// @Summary Get the price of a token-unit pair at a given timestamp
// @Description Retrieves the most recent price of a token-unit pair in the specified unit before the given timestamp
// @Tags price
// @Produce json
// @Param token query string true "The token part of the pair"
// @Param unit query string false "The unit part of the pair (default: USD)"
// @Param timestamp query string false "The fetch timestamp to retrieve the price for, in RFC3339 format (e.g. 2023-03-18T01:23:45+08:00) (default: current time)"
// @Success 200 {object} PriceResponse "Returns the price of the token"
// @Failure 400 {object} ErrorResponse "Returns an error if the request is invalid"
// @Failure 404 {object} ErrorResponse "Returns 404 if the price info is not found at the given timestamp"
// @Router /price [get]
func HandlePriceRequest(c *gin.Context) {
	req, err := ParsePriceRequest(c)
	if err != nil {
		c.AbortWithStatusJSON(400, Error(err))
		return
	}
	res, err := QueryPrice(req, GetConn(c))
	if ok := HandleDBError(c, err); !ok {
		return
	}
	c.JSON(200, res)
}

func QueryPrice(q PriceRequest, conn *pgx.Conn) (PriceResponse, error) {
	row := conn.QueryRow(
		context.Background(),
		`SELECT price, price_timestamp, fetch_timestamp FROM prices
			WHERE token = $1 AND unit = $2 AND fetch_timestamp <= $3
			ORDER BY fetch_timestamp DESC
			LIMIT 1`,
		q.Token, q.Unit, q.Timestamp,
	)
	res := PriceResponse{
		Token: q.Token,
		Unit:  q.Unit,
	}
	err := row.Scan(&res.Price, &res.PriceTimestamp, &res.FetchTimestamp)
	if err != nil {
		return PriceResponse{}, err
	}
	res.PriceTimestamp = res.PriceTimestamp.UTC()
	res.FetchTimestamp = res.FetchTimestamp.UTC()
	return res, nil
}
