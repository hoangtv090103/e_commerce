package controllers

import (
	"e-commerce/models"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)
func AddToCart(ctx *gin.Context, client *redis.Client) {
	var cartItem models.CartItem

	err := ctx.BindJSON(&cartItem)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	cachedCartItem, err := client.Get("cart:" + strconv.Itoa(int(cartItem.ProductID))).Result()

	if err == redis.Nil {
		// product not in the cart
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving data from Redis"})
		return
	}

	if cachedCartItem != "" {
		var storedCartItems models.CartItem
		err := json.Unmarshal([]byte(cachedCartItem), &storedCartItems)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error unmarshaling data"})
			return
		}

		// increment the quantity in the cart incrby
		client.IncrBy("cart:" + strconv.Itoa(int(cartItem.ProductID)), int64(cartItem.Quantity))
	}

	// set the cart item in the cart
	cartItemJSON, err := json.Marshal(cartItem)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshaling data"})
		return
	}
	err = client.Set("cart:" + string(cartItem.ProductID), cartItemJSON, time.Hour).Err()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating Redis cache"})
		return
	}

	ctx.JSON(http.StatusOK, "Cart updated")

	return
}
func RemoveFromCart(ctx *gin.Context, client *redis.Client) {
}

func GetCart(ctx *gin.Context, client *redis.Client) {
}

func Checkout(ctx *gin.Context, client *redis.Client){
}
