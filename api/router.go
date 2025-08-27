package api

import (
	"subscription/api/routes"
	"subscription/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(middleware.LoggerMiddleware())

	routes.SubscriptionRouter(router)

	return router
}
