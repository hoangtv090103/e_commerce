package models

type OrderItem struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	OrderID    uint      `json:"order_id"`
	ProductID  uint      `json:"product_id"`
	Quantity   int       `gorm:"not null" json:"quantity"`
	Price      float64   `gorm:"not null" json:"price"`
	TotalPrice float64   `gorm:"not null" json:"total_price"`
	Product    Product   `json:"product"`
}