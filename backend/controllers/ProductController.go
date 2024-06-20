package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"e-commerce/db"
	"e-commerce/models"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

func GetProducts(c *gin.Context, r *redis.Client) {
	cachedProducts, err := r.Get("products:all").Result()
	if err != nil && err != redis.Nil {
		fmt.Println("Error retrieving from Redis:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving data"})
		return
	}

	if cachedProducts != "" {
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
	if err := db.DB.Preload("Categories").Find(&products).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "No products found"})
		} else {
			fmt.Println("Error querying database:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying database"})
		}
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
		r.Set("products:"+strconv.Itoa(int(product.ID)), string(productJSON), time.Hour)
	}

	c.JSON(http.StatusOK, products)
}

func GetProduct(c *gin.Context, r *redis.Client) {
	productId := c.Param("id")
	cachedProduct, err := r.Get("products:" + productId).Result()
	if err != nil && err != redis.Nil {
		fmt.Println("Error retrieving from Redis:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving data"})
		return
	}

	if cachedProduct != "" {
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
	if err := db.DB.Preload("Categories").First(&product, productId).Error; err != nil {
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
	r.Set("products:"+productId, string(productJSON), time.Hour)

	c.JSON(http.StatusOK, product)
}

func AddProduct(c *gin.Context, r *redis.Client) {
	product := models.Product{}
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if product.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	tx := db.DB.Begin()
	if err := db.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating product"})
		tx.Rollback()
		return
	}

	productJSON, err := json.Marshal(product)
	if err != nil {
		log.Println("Error marshaling product:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshaling data"})
		return
	}
	r.Set("products:"+strconv.Itoa(int(product.ID)), string(productJSON), time.Hour)

	cachedProducts, err := r.Get("products:all").Result()
	if err != nil && err != redis.Nil {
		log.Println("Error retrieving from Redis:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving data"})
		return
	}

	var products []models.Product
	if cachedProducts != "" {
		if err := json.Unmarshal([]byte(cachedProducts), &products); err != nil {
			log.Println("Error unmarshaling cached products:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error unmarshaling data"})
			return
		}
	}
	products = append(products, product)
	productsJSON, err := json.Marshal(products)
	if err != nil {
		log.Println("Error marshaling products:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshaling data"})
		return
	}
	r.Set("products:all", productsJSON, time.Hour)

	c.JSON(http.StatusOK, product)
}

func UpdateProduct(c *gin.Context, r *redis.Client) {
	product := models.Product{}
	id := c.Param("id")
	db.DB.First(&product, id)
	c.BindJSON(&product)
	db.DB.Save(&product)

	productJSON, err := json.Marshal(product)
	if err != nil {
		fmt.Println("Error marshaling product:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshaling data"})
		return
	}
	r.Set("products:"+id, string(productJSON), time.Hour)

	c.JSON(http.StatusOK, product)
}

func DeleteProduct(c *gin.Context, r *redis.Client) {
	id := c.Param("id")
	r.Del("products:" + id)
	var product models.Product
	if err := db.DB.First(&product, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	db.DB.Delete(&models.Product{}, id)

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})
}
