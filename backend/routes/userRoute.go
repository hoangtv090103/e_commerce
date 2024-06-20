package routes

import (
	"e-commerce/controllers"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func UserRouter(superGroup *gin.RouterGroup, client *redis.Client) {
	userGroup := superGroup.Group("/users")
	{
		userGroup.GET("/", func (c *gin.Context) {
			controllers.GetUsers(c, client)
		})		
	}
}