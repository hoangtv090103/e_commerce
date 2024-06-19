package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
	
	"e-commerce/controllers"
)

func SetupRouter(db *gorm.DB, client *redis.Client) *gin.Engine {
	r := gin.Default()

	// Middleware
	// r.Use(middleware.LoggingMiddleware())
	// r.Use(middleware.AuthMiddleware())
	productRouter := r.Group("/api/v1/products")
	{
		productRouter.GET("/", func(ctx *gin.Context) {
			controllers.GetProducts(db, ctx, client)
		})
		productRouter.GET("/:id", func(ctx *gin.Context) {
			controllers.GetProduct(db, ctx, client)
		})
		productRouter.POST("/", func(ctx *gin.Context) {
			controllers.AddProduct(db, ctx, client)
		})
		productRouter.POST("/:id", func(ctx *gin.Context) {
			controllers.UpdateProduct(db, ctx, client)
		})
		productRouter.DELETE("/:id", func(ctx *gin.Context) {
			controllers.DeleteProduct(db, ctx, client)
		})
	}
	return r
}
