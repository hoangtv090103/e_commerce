package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func SetupRouter(client *redis.Client) *gin.Engine {
	r := gin.Default()

	// Middleware
	// r.Use(middleware.LoggingMiddleware())
	// r.Use(middleware.AuthMiddleware())
	apiRouter := r.Group("/api/v1")
	ProductRouter(apiRouter, client)
	AuthRouter(apiRouter, client)
	return r
}