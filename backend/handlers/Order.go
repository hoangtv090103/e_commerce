package handlers

import (
	"e-commerce/db"
	"e-commerce/models"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

func GetOrders(db *gorm.DB, c *gin.Context, r *redis.Client) {
	cachedOrders, err := r.Get("orders:all").Result()
	var orders []models.Order
	if err == redis.Nil {
		log.Println("Cache miss, querying database")
	} else if err != nil {
		log.Println("Error retrieving from Redis: ", err)
 	} else {
  		log.Println("Returning cached Orders")
	    if err := json.Unmarshal([]byte(cachedOrders), &orders); err != nil {
	    	log.Println("Error unmarshaling cached Orders: ", err)
	    }

	    c.JSON(http.StatusOK,orders)
	    return
  	}

	db.Model(&models.Order{}).Find(&orders)

	ordersJson, err := json.Marshal(orders)

	if err != nil {
	  	log.Println("Error  marshaling orders: ", err)
	   	c.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshaling data"})
	   	return
	}

	r.Set("orders:all", ordersJson, time.Hour)

	for _, order := range orders {
	  	orderJSON, err := json.Marshal(order)
	   if err != nil {
		   log.Println("Error marshaling product: ", err)
		   c.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshaling data"})
	   }
	   r.HSet("orders:" + string(rune(order.ID)), string(orderJSON), time.Hour)
	}
}

func GetOrder(c *gin.Context, r *redis.Client) {
	id := c.Param("id")
	cachedOrder, err := r.HGet("orders:" + id, id).Result()
	var order models.Order
	if err == redis.Nil {
		log.Println("Cache miss, querying database")
	} else if err != nil {
		log.Println("Error retrieving from Redis: ", err)
	} else {
		log.Println("Returning cached Order")
		if err := json.Unmarshal([]byte(cachedOrder), &order); err != nil {
			log.Println("Error unmarshaling cached Order: ", err)
		}

		c.JSON(http.StatusOK, order)
		return
	}

	DB := db.Connect()
	defer db.Close(DB)
	DB.First(&order, id)

	orderJson, err := json.Marshal(order)

	if err != nil {
		log.Println("Error marshaling order: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshaling data"})
		return
	}

	r.HSet("orders:" + id, string(orderJson), time.Hour)
	c.JSON(http.StatusOK, order)
}
