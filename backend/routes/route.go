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
	
	userRouter := r.Group("/api/v1/users")
	{
		userRouter.GET("/", func(ctx *gin.Context) {
			controllers.GetUsers(db, ctx, client)
		})
		userRouter.GET("/:id", func(ctx *gin.Context) {
			controllers.GetUser(db, ctx, client)
		})
	}
	
	orderRouter := r.Group("/api/v1/orders")
	{
		orderRouter.GET("/", func(ctx *gin.Context) {
			controllers.GetOrders(db, ctx, client)
		})
		
		orderRouter.GET("/:id", func(ctx *gin.Context) {
			controllers.GetOrder(db, ctx, client)
		})
		
		orderRouter.POST("/", func(ctx *gin.Context) {
			controllers.AddOrder(db, ctx, client)
		})
		
		orderRouter.POST("/:id", func(ctx *gin.Context) {
			controllers.UpdateOrder(db, ctx, client)
		})
		
		orderRouter.DELETE("/:id", func(ctx *gin.Context) {
			controllers.DeleteOrder(db, ctx, client)
		})
	}
		
	return r
}
