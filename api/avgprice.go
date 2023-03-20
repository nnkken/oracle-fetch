package api

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/gin-gonic/gin"
)

type AvgPriceRequest struct {
	Token string `form:"token" binding:"required"`
	Unit  string `form:"unit,default=USD"`

	// TODO: allow 0 = default now?
	From time.Time `form:"from" binding:"required"`
	// TODO: allow 0 = default 1970-01-01?
	To time.Time `form:"to" binding:"required"`
}

type AvgPriceResponse struct {
	Token               string    `json:"token"`
	Unit                string    `json:"unit"`
	AvgPrice            string    `json:"avg_price"`
	PriceCount          uint      `json:"price_count"`
	FirstFetchTimestamp time.Time `json:"first_price_timestamp"`
	LastFetchTimestamp  time.Time `json:"last_price_timestamp"`
}

func ParseAvgPriceRequest(c *gin.Context) (AvgPriceRequest, error) {
	var req AvgPriceRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		return AvgPriceRequest{}, fmt.Errorf("fail to parse avg_price request: %w", err)
	}
	req.From = req.From.UTC()
	req.To = req.To.UTC()
	if req.From.After(req.To) {
		return AvgPriceRequest{}, fmt.Errorf("'from' must be before 'to'")
	}
	return req, nil
}

func HandleAvgPriceRequest(c *gin.Context) {
	req, err := ParseAvgPriceRequest(c)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}
	res, err := QueryAvgPrice(req, GetConn(c))
	if ok := HandleDBError(c, err); !ok {
		return
	}
	c.JSON(200, res)
}

func QueryAvgPrice(q AvgPriceRequest, conn *pgx.Conn) (AvgPriceResponse, error) {
	row := conn.QueryRow(
		context.Background(),
		`
			SELECT
				coalesce(AVG(price), 0),
				COUNT(*),
				coalesce(MIN(fetch_timestamp),to_timestamp(0)),
				coalesce(MAX(fetch_timestamp), to_timestamp(0))
			FROM prices
			WHERE token = $1
				AND unit = $2
				AND fetch_timestamp >= $3
				AND fetch_timestamp <= $4
		`, q.Token, q.Unit, q.From, q.To,
	)

	res := AvgPriceResponse{
		Token: q.Token,
		Unit:  q.Unit,
	}
	err := row.Scan(&res.AvgPrice, &res.PriceCount, &res.FirstFetchTimestamp, &res.LastFetchTimestamp)
	if err != nil {
		return AvgPriceResponse{}, err
	}
	if res.PriceCount == 0 {
		return AvgPriceResponse{}, pgx.ErrNoRows
	}
	res.FirstFetchTimestamp = res.FirstFetchTimestamp.UTC()
	res.LastFetchTimestamp = res.LastFetchTimestamp.UTC()
	return res, nil
}
