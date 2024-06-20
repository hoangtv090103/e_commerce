package models

import "time"

type Order struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	UserID      uint        `gorm:"not null" json:"user_id"`
	Status      string      `gorm:"not null" json:"status"`
	TotalAmount float64     `gorm:"not null" json:"total_amount"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	OrderLines  []OrderLine `json:"order_lines"`
}
