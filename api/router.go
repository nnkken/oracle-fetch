package api

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/gin-gonic/gin"

	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	"github.com/nnkken/oracle-fetch/docs"
	"github.com/nnkken/oracle-fetch/utils"
)

func NewRouter(connPool *pgxpool.Pool) *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		log := utils.GetLogger("api")
		v1.Use(WithLogger(log))
		v1.Use(WithPool(connPool))
		v1.GET("/price", HandlePriceRequest)
		v1.GET("/avg_price", HandleAvgPriceRequest)
	}

	docs.SwaggerInfo.Title = "Oracle Fetch API"
	docs.SwaggerInfo.Description = "This is a service for fetching price data periodically from oracles. This API provides access to the fetched price data."
	docs.SwaggerInfo.BasePath = v1.BasePath()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
