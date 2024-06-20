package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"

	"e-commerce/controllers"
)

func ProductRouter(superRouter *gin.RouterGroup, client *redis.Client){
 productRouter := superRouter.Group("/products")
 {
  productRouter.GET("/", func(ctx *gin.Context) {
   controllers.GetProducts(ctx, client)
  })
  productRouter.GET("/:id", func(ctx *gin.Context) {
   controllers.GetProduct(ctx, client)
  })
  productRouter.POST("/", func(ctx *gin.Context) {
   controllers.AddProduct(ctx, client)
  })
  productRouter.PUT("/:id", func(ctx *gin.Context) {
   controllers.UpdateProduct(ctx, client)
  })
  productRouter.DELETE("/:id", func(ctx *gin.Context) {
   controllers.DeleteProduct(ctx, client)
  })
 }
}