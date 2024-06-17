package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"e-commerce/pkg/handlers"
)

func main() {
	env := godotenv.Load()
	if env != nil {
		panic("cannot find environment variables")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	} else {
		port = ":" + port
	}

	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.GET("/", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{
				"message": "Hello, World!",
			})

		})
		api.GET("/products", handlers.GetProducts)
	}

	r.Run(port)
}
