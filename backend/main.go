package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"e-commerce/pkg/config"
	"e-commerce/pkg/handlers"
)

func main() {
	env := godotenv.Load()
	if env != nil {
		panic("cannot find environment variables")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()
	client := config.ConnectRedis()
	api := r.Group("/api/v1")
	{
		api.GET("/", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "Hello, World!",
			})

		})
		api.GET("/products", func(ctx *gin.Context) {
			handlers.GetProducts(ctx, client)
		})
		api.GET("/products/:id", func(ctx *gin.Context) {
			handlers.GetProduct(ctx, client)
		})
		api.POST("/products", func(ctx *gin.Context) {
			handlers.AddProduct(ctx, client)
		})
		api.POST("/products/:id", func(ctx *gin.Context) {
			handlers.UpdateProduct(ctx, client)
		})
		api.DELETE("/products/:id", func(ctx *gin.Context) {
			handlers.DeleteProduct(ctx, client)
		})

	}

	r.Run(":" + port)
}
