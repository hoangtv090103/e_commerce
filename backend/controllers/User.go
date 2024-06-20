package controllers

import (
	"e-commerce/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

func GetUsers(db *gorm.DB, ctx *gin.Context, client *redis.Client) {
	var users []models.User
	db.Model(&models.User{}).Find(&users)
	if len(users) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No users found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": users})
}

func GetUser(db *gorm.DB, ctx *gin.Context, client *redis.Client) {
	userID := ctx.MustGet("user_id").(uint)
	var user models.User
	db.Model(&models.User{}).First(&user, userID)
	ctx.JSON(http.StatusOK, gin.H{"data": user})
}
