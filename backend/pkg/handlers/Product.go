package handlers

import (
	"net/http"

	"e-commerce/models"
	"e-commerce/pkg/db"

	"github.com/gin-gonic/gin"
)

func GetProducts(c *gin.Context) {
	DB := db.Connect()
	products := []models.Product{}
	DB.Find(&products)

	c.JSON(http.StatusOK, products)
	defer db.Close(DB)
}
