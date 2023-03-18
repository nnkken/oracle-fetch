package utils

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gin-gonic/gin"
)

func HandleDBError(c *gin.Context, err error) (ok bool) {
	switch err {
	case nil:
		return true
	case pgx.ErrNoRows:
		c.AbortWithStatus(404)
		return false
	default:
		c.AbortWithStatus(500)
		return false
	}
}

func WithPool(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := pool.Acquire(context.Background())
		if err != nil {
			// TODO: log
			c.AbortWithStatusJSON(500, gin.H{"error": "failed to acquire connection"})
			return
		}
		c.Set("conn", conn)
		defer conn.Release()
		c.Next()
	}
}

func GetConn(c *gin.Context) *pgxpool.Conn {
	return c.MustGet("conn").(*pgxpool.Conn)
}
