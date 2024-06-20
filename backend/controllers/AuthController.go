package controllers

import (
	"e-commerce/db"
	"e-commerce/models"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(ctx *gin.Context, client *redis.Client) error {
	var user models.User

	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return err
	}

	// Check if user already exists. If user exists, return error
	if existingUser := db.DB.Model(&models.User{}).Where("email = ?", user.Email).First(&user); existingUser.Error == nil {
		ctx.JSON(400, gin.H{"error": "User already exists"})
		return nil
	}

	// Start a transaction
	tx := db.DB.Begin()

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Error hashing password"})
		tx.Rollback()
		return err
	}
	user.Password = string(hashedPassword)

	// Create user
	err = tx.Create(&user).Error

	if err != nil {
		ctx.JSON(500, gin.H{"error": "Error creating user"})
		tx.Rollback()
		return err
	}

	// Commit the transaction
	tx.Commit()

	ctx.JSON(200, gin.H{"message": "User created successfully"})

	// Cache user
	userJSON, err := json.Marshal(user)
	if err != nil {
		log.Println("Error marshaling user: ", err)
		return err
	}
	client.Set("users:"+strconv.Itoa(int(user.ID)), userJSON, time.Hour)

	return nil
}

func Login(ctx *gin.Context, client *redis.Client) {
	var loginUser models.User
	var err error

	if err = ctx.BindJSON(&loginUser); err != nil {
		ctx.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	// Get the stored hashed password for the given email
	var user models.User
	existingUser := db.DB.Model(&models.User{}).Where("email = ?", loginUser.Email).First(&user)
	if existingUser.Error != nil {
		ctx.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare the incoming password with the stored hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUser.Password))
	if err != nil {
		ctx.JSON(401, gin.H{"error": "Invalid email or password"})
		return
	}

	// Cache user
	userJSON, err := json.Marshal(user)
	if err != nil {
		log.Println("Error marshaling user: ", err)
		return
	}
	client.Set("users:"+strconv.Itoa(int(user.ID)), userJSON, time.Hour)

	ctx.JSON(200, gin.H{"message": "Login successful"})

}
