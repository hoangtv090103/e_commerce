package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"

	"e-commerce/controllers"
)

func CartRouter(superGroup *gin.RouterGroup, client *redis.Client) {
	cartRouter := superGroup.Group("/cart")
	{
		cartRouter.POST("/", func(ctx *gin.Context) {
			controllers.AddToCart(ctx, client)
		})
	}
}
