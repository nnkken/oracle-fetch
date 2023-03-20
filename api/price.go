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
	Token          string    `json:"token"`
	Unit           string    `json:"unit"`
	Price          string    `json:"price"`
	PriceTimestamp time.Time `json:"price_timestamp"`
	FetchTimestamp time.Time `json:"fetch_timestamp"`
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

func HandlePriceRequest(c *gin.Context) {
	req, err := ParsePriceRequest(c)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
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
