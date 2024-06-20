	package controllers

import (
	"e-commerce/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

func GetOrders(db *gorm.DB, ctx *gin.Context, client *redis.Client) {
	cachedOrders, err := client.Get("orders:all").Result()
	if err == redis.Nil {
		log.Println("Cache miss, querying database")
	} else if err != nil {
		log.Println("Error retrieving from Redis: ", err)
	} else {
		log.Println("Returning cached Orders")
		var orders []models.Order
		if err := json.Unmarshal([]byte(cachedOrders), &orders); err != nil {
			log.Println("Error unmarshaling cached Orders: ", err)
		}
		ctx.JSON(http.StatusOK, orders)
		return
	}

	var orders []models.Order
	db.Model(&models.Order{}).Find(&orders)

	ordersJson, err := json.Marshal(orders)
	if err != nil {
		log.Println("Error marshaling orders: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshaling data"})
		return
	}

	client.Set("orders:all", ordersJson, time.Hour)

	for _, order := range orders {
		orderJSON, err := json.Marshal(order)
		if err != nil {
			log.Println("Error marshaling product: ", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshaling data"})
		}
		client.Set("orders:"+strconv.Itoa(int(order.ID)), orderJSON, time.Hour)
	}

	ctx.JSON(http.StatusOK, gin.H{"data": orders})
}

func GetOrder(db *gorm.DB, ctx *gin.Context, client *redis.Client) {
	id := ctx.Param("id")
	cachedOrder, err := client.Get("orders:" + id).Result()
	if err == redis.Nil {
		log.Println("Cache miss, querying database")
	} else if err != nil {
		log.Println("Error retrieving from Redis: ", err)
	} else {
		log.Println("Returning cached Order")
		var order models.Order
		if err := json.Unmarshal([]byte(cachedOrder), &order); err != nil {
			log.Println("Error unmarshaling cached Order: ", err)
		}
		ctx.JSON(http.StatusOK, order)
		return
	}

	var order models.Order
	db.Model(&models.Order{}).First(&order, id)

	orderJSON, err := json.Marshal(order)
	if err != nil {
		log.Println("Error marshaling order: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshaling data"})
		return
	}

	client.Set("orders:"+id, orderJSON, time.Hour)

	ctx.JSON(http.StatusOK, gin.H{"data": order})
}

func AddOrder(db *gorm.DB, ctx *gin.Context, client *redis.Client) {
	var order models.Order
	if err := ctx.ShouldBindJSON(&order); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Create(&order).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	orderJSON, err := json.Marshal(order)
	if err != nil {
		log.Println("Error marshaling order: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshaling data"})
		return
	}

	client.Set("orders:"+strconv.Itoa(int(order.ID)), orderJSON, time.Hour)

	ctx.JSON(http.StatusOK, gin.H{"data": order})
}

func UpdateOrder(db *gorm.DB, ctx *gin.Context, client *redis.Client) {
	id := ctx.Param("id")
	var order models.Order
	if err := db.Model(&models.Order{}).First(&order, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if err := ctx.ShouldBindJSON(&order); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Save(&order).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	orderJSON, err := json.Marshal(order)
	if err != nil {
		log.Println("Error marshaling order: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshaling data"})
		return
	}

	client.Set("orders:"+id, orderJSON, time.Hour)

	ctx.JSON(http.StatusOK, gin.H{"data": order})
}

func DeleteOrder(db *gorm.DB, ctx *gin.Context, client *redis.Client) {
	id := ctx.Param("id")
	var order models.Order
	if err := db.Model(&models.Order{}).First(&order, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	if err := db.Delete(&order).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	client.Del("orders:" + id)

	ctx.JSON(http.StatusOK, gin.H{"data": "Order deleted"})
}
