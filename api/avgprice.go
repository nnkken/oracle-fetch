package api

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

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

func HandleAvgPriceRequest(c *gin.Context) {
	var q AvgPriceRequest
	if err := c.ShouldBindQuery(&q); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}
	q.From = q.From.UTC()
	q.To = q.To.UTC()
	if q.From.After(q.To) {
		c.AbortWithStatusJSON(400, gin.H{"error": "'from' must be before 'to'"})
		return
	}
	res, err := QueryAvgPrice(q, GetConn(c))
	if ok := HandleDBError(c, err); !ok {
		return
	}
	c.JSON(200, res)
}

func QueryAvgPrice(q AvgPriceRequest, conn *pgxpool.Conn) (AvgPriceResponse, error) {
	// TODO: when no data is available, null happens for the fields, which makes Scan return error and cause 500
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
		// TODO: log
		return AvgPriceResponse{}, err
	}
	if res.PriceCount == 0 {
		return AvgPriceResponse{}, pgx.ErrNoRows
	}
	res.FirstFetchTimestamp = res.FirstFetchTimestamp.UTC()
	res.LastFetchTimestamp = res.LastFetchTimestamp.UTC()
	return res, nil
}
