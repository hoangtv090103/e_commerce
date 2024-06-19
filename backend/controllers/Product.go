package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"e-commerce/models"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

// HSet vs Set
// HSet is used to store a hash map in Redis
// Set is used to store a string in Redis

func GetProducts(db *gorm.DB,c *gin.Context, r *redis.Client) {
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

	var products []models.Product
	err = db.Model(&models.Product{}).Preload("Categories").Find(&products).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "No products found"})
			return
		}
		fmt.Println("Error querying database:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying database"})
		return
	}

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

func GetProduct(db *gorm.DB, c *gin.Context, r *redis.Client) {
	productId := c.Param("id")

	cachedProduct, err := r.HGet("products:"+productId, productId).Result()
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

	var product models.Product
	if err := db.Model(&models.Product{}).Preload("Categories").First(&product, productId).Error; err != nil {
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
	r.HSet("products:"+productId, string(productJSON), time.Hour)

	c.JSON(http.StatusOK, product)
}

func AddProduct(db *gorm.DB, c *gin.Context, r *redis.Client) {
	products := []models.Product{}

	c.BindJSON(&products)
	if len(products) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No products provided"})
		return
	}
	
	db.Model(&models.Product{}).Create(&products)

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

func UpdateProduct(db *gorm.DB, c *gin.Context, r *redis.Client) {
	product := models.Product{}
	id := c.Param("id")
	db.Model(&models.Product{}).First(&product, id)
	
	c.BindJSON(&product)
	
	db.Save(&product)

	productJSON, err := json.Marshal(product)
	if err != nil {
		fmt.Println("Error marshaling product:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshaling data"})
		return
	}
	r.HSet("products:"+id, string(productJSON), time.Hour)

	c.JSON(http.StatusOK, product)
}

func DeleteProduct(db *gorm.DB, c *gin.Context, r *redis.Client) {
	product := models.Product{}
	id := c.Param("id")
	db.Model(&models.Product{}).First(&product, id)

	r.HDel("products", id)

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})
}
