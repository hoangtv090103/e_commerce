package db

import (
	"log"
	"os"

	"gorm.io/gorm"

	"e-commerce/models"
	"gorm.io/driver/postgres"
)

func Connect() *gorm.DB {
	dbURL := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		PrepareStmt: true,
	})

	if err != nil {
		log.Panic("Failed to connect to database: ", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Product{}, &models.Category{}, &models.ProductCategory{}, &models.Order{}, &models.OrderItem{}, &models.Cart{}, &models.CartItem{}, &models.Payment{})

	if err != nil {
		log.Panic("Failed to migrate database: ", err)
	}

	// return db
	return db
}

func Close(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Panic("Failed to get database connection: ", err)
	}

	err = sqlDB.Close()
	if err != nil {
		log.Panic("Failed to close database connection: ", err)
	}
}
