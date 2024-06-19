package models

import "time"

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	Price       float64   `gorm:"not null" json:"price"`
	Stock       int       `gorm:"not null" json:"stock"`
	Categories  []Category `gorm:"many2many:product_categories" json:"categories"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}