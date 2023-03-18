package price

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nnkken/oracle-fetch/api/utils"
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

func HandlePriceRequest(c *gin.Context) {
	var q PriceRequest
	if err := c.ShouldBindQuery(&q); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}
	if q.Timestamp.IsZero() {
		q.Timestamp = time.Now()
	}
	q.Timestamp = q.Timestamp.UTC()
	res, err := QueryPrice(q, utils.GetConn(c))
	if ok := utils.HandleDBError(c, err); !ok {
		return
	}
	c.JSON(200, res)
}

func QueryPrice(q PriceRequest, conn *pgxpool.Conn) (PriceResponse, error) {
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
