package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	
	"e-commerce/controllers"
)

func AuthRouter(superRouter *gin.RouterGroup, client *redis.Client) {
	authRouter := superRouter.Group("/auth")
	{
		authRouter.POST("/register", func(ctx *gin.Context) {
			controllers.SignUp(ctx, client)
		})
		authRouter.POST("/login", func(ctx *gin.Context) {
			controllers.Login(ctx, client)
		})
	}
	
}