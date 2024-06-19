package models

import "time"

type Cart struct {
	ID 	  	  uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CartItems []CartItem `json:"cart_items"`
}