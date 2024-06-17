package models

import "time"

type Cart struct {
	UserID    uint      `gorm:"primaryKey" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CartItems []CartItem `json:"cart_items"`
}