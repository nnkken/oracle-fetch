package api

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gin-gonic/gin"
)

func NewRouter(connPool *pgxpool.Pool) *gin.Engine {
	router := gin.Default()
	router.Use(WithPool(connPool))
	router.GET("/price", HandlePriceRequest)
	router.GET("/avg_price", HandleAvgPriceRequest)
	return router
}
