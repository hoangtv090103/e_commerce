package models

import "time"

type Order struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	UserID      uint       `json:"user_id"`
	Status      string     `gorm:"not null" json:"status"`
	TotalAmount float64    `gorm:"not null" json:"total_amount"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	OrderItems  []OrderItem `json:"order_items"`
}
