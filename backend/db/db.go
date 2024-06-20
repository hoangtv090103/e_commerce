package db

import (
	"log"
	"os"

	"gorm.io/gorm"

	"e-commerce/models"
	"gorm.io/driver/postgres"
)

var DB *gorm.DB

func InitDB() {
 dbURL := os.Getenv("DB_URL")
 var err error
 DB, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{
  PrepareStmt: true,
 })

 if err != nil {
  log.Panic("Failed to connect to database: ", err)
 }

 err = DB.Debug().AutoMigrate(&models.User{}, &models.Product{}, &models.Category{}, &models.Order{}, &models.OrderLine{}, &models.Cart{}, &models.CartItem{}, &models.Payment{})

 if err != nil {
  log.Panic("Failed to migrate database: ", err)
 }

  log.Println("Connected to database")
}