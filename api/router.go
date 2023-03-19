package api

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gin-gonic/gin"

	"github.com/nnkken/oracle-fetch/logger"
)

func NewRouter(connPool *pgxpool.Pool) *gin.Engine {
	router := gin.Default()
	log := logger.GetLogger("api")
	router.Use(WithLogger(log))
	router.Use(WithPool(connPool))
	router.GET("/price", HandlePriceRequest)
	router.GET("/avg_price", HandleAvgPriceRequest)
	return router
}
