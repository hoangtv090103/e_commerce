package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"e-commerce/models"
	"e-commerce/pkg/db"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

// HSet vs Set
// HSet is used to store a hash map in Redis
// Set is used to store a string in Redis

func GetProducts(c *gin.Context, r *redis.Client) {
	cachedProducts, err := r.Get("products:all").Result()
	if err == redis.Nil {
		fmt.Println("Cache miss, querying database")
	} else if err != nil {
		fmt.Println("Error retrieving from Redis:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving data"})
		return
	} else {
		fmt.Println("Returning cached products")
		var products []models.Product
		if err := json.Unmarshal([]byte(cachedProducts), &products); err != nil {
			fmt.Println("Error unmarshaling cached products:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error unmarshaling data"})
			return
		}
		c.JSON(http.StatusOK, products)
		return
	}

	DB := db.Connect()
	defer db.Close(DB)
	var products []models.Product
	DB.Find(&products)

	productsJSON, err := json.Marshal(products)
	if err != nil {
		fmt.Println("Error marshaling products:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshaling data"})
		return
	}
	r.Set("products:all", productsJSON, time.Hour)

	for _, product := range products {
		productJSON, err := json.Marshal(product)
		if err != nil {
			fmt.Println("Error marshaling product:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshaling data"})
			return
		}
		r.HSet("products:"+string(rune(product.ID)), string(productJSON), time.Hour)
	}

	c.JSON(http.StatusOK, products)
}

func GetProduct(c *gin.Context, r *redis.Client) {
	productId := c.Param("id")

	cachedProduct, err := r.Get("products:" + productId).Result()
	if err == redis.Nil {
		fmt.Println("Cache miss, querying database")
	} else if err != nil {
		fmt.Println("Error retrieving from Redis:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving data"})
		return
	} else {
		var product models.Product
		if err := json.Unmarshal([]byte(cachedProduct), &product); err != nil {
			fmt.Println("Error unmarshaling cached product:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error unmarshaling data"})
			return
		}
		c.JSON(http.StatusOK, product)
		return
	}

	DB := db.Connect()
	defer db.Close(DB)
	var product models.Product
	if err := DB.First(&product, productId).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	productJSON, err := json.Marshal(product)
	if err != nil {
		fmt.Println("Error marshaling product:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshaling data"})
		return
	}
	r.Set("products:"+productId, productJSON, time.Hour)

	c.JSON(http.StatusOK, product)
}

func AddProduct(c *gin.Context, r *redis.Client) {
	DB := db.Connect()
	defer db.Close(DB)
	products := []models.Product{}

	c.BindJSON(&products)
	DB.Create(&products)

	for _, product := range products {
		productJSON, err := json.Marshal(product)
		if err != nil {
			fmt.Println("Error marshaling product:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshaling data"})
			return
		}
		r.HSet("products:"+string(rune(product.ID)), string(productJSON), time.Hour)

		// update redis products:all
		cachedProducts, err := r.Get("products:all").Result()
		if err == redis.Nil {
			fmt.Println("Cache miss, querying database")
		} else if err != nil {
			log.Println("Error: ", err)
		} else {
			var products []models.Product
			if err := json.Unmarshal([]byte(cachedProducts), &products); err != nil {
				log.Println("Error: ", err)
			}
			products = append(products, product)
			productsJSON, err := json.Marshal(products)
			if err != nil {
				log.Println("Error: ", err)
			}
			r.Set("products:all", productsJSON, time.Hour)

		}
	}

	c.JSON(http.StatusOK, products)
}

func UpdateProduct(c *gin.Context, r *redis.Client) {
	DB := db.Connect()
	defer db.Close(DB)
	product := models.Product{}
	id := c.Param("id")
	DB.First(&product, id)
	c.BindJSON(&product)
	DB.Save(&product)

	productJSON, err := json.Marshal(product)
	if err != nil {
		fmt.Println("Error marshaling product:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshaling data"})
		return
	}
	r.HSet("products:"+id, string(productJSON), time.Hour)

	c.JSON(http.StatusOK, product)
}

func DeleteProduct(c *gin.Context, r *redis.Client) {
	DB := db.Connect()
	defer db.Close(DB)
	product := models.Product{}
	id := c.Param("id")
	DB.Delete(&product, id)

	r.HDel("products", id)

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})
}
