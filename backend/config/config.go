package config

import (
	"log"
	"os"

	"gorm.io/gorm"

	"e-commerce/models"
	"gorm.io/driver/postgres"
)

var DBConn *gorm.DB

func ConnectDatabase() {
	dbURL := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Panic("Failed to connect to database: ", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Product{}, &models.Category{}, &models.ProductCategory{}, &models.Order{}, &models.OrderItem{}, &models.Cart{}, &models.CartItem{}, &models.Payment{})

	if err != nil {
		log.Panic("Failed to migrate database: ", err)
	}
	DBConn = db
}

// GetDB is a function that returns the current database connection.
// It does not take any parameters.
//
// Returns:
// *gorm.DB: The current database connection.
func GetDB() *gorm.DB {
	return DBConn
}
