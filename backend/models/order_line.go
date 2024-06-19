package models

type OrderLine struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Order      Order     `json:"order"`
	OrderID    uint      `json:"order_id"`
	ProductID  uint      `json:"product_id"`
	Quantity   int       `gorm:"not null" json:"quantity"`
	Price      float64   `gorm:"not null" json:"price"`
	TotalPrice float64   `gorm:"not null" json:"total_price"`
}