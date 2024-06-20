package controllers

import (
	"e-commerce/db"
	"e-commerce/models"
	"encoding/json"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

func GetUsers(ctx *gin.Context, client *redis.Client) {
	var users []models.User
	cachedUsers, err := client.Get("users:all").Result()
	if err != nil && err != redis.Nil {
		ctx.JSON(500, gin.H{"error": "Error retrieving data"})
		return
	}

	if cachedUsers != "" {
		if err := json.Unmarshal([]byte(cachedUsers), &users); err != nil {
			ctx.JSON(500, gin.H{"error": "Error unmarshaling data"})
			return
		}
		ctx.JSON(200, users)
		return
	}

	if err := db.DB.Find(&users).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "Error querying database"})
		return
	}

	usersJSON, err := json.Marshal(users)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Error marshaling data"})
		return
	}

	client.Set("users:all", usersJSON, time.Hour)

	for _, user := range users {
		userJSON, err := json.Marshal(user)
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Error marshaling data"})
			return
		}
		client.Set("user:"+strconv.Itoa(int(user.ID)), userJSON, time.Hour)
	}

	ctx.JSON(200, users)

}
