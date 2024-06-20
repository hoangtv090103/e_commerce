package main

import (
	"os"

	"e-commerce/config"
	"e-commerce/db"
	"e-commerce/routes"

	"github.com/joho/godotenv"
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

	db.InitDB()
	client := config.ConnectRedis()
	r := routes.SetupRouter(client)
	
	r.Run(":" + port)
}
