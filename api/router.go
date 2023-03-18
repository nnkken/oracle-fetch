package api

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gin-gonic/gin"

	"github.com/nnkken/oracle-fetch/api/avgprice"
	"github.com/nnkken/oracle-fetch/api/price"
	"github.com/nnkken/oracle-fetch/api/utils"
)

func NewRouter(connPool *pgxpool.Pool) *gin.Engine {
	router := gin.Default()
	router.Use(utils.WithPool(connPool))
	router.GET("/price", price.HandlePriceRequest)
	router.GET("/avg_price", avgprice.HandleAvgPriceRequest)
	return router
}
