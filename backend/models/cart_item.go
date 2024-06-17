package models

type CartItem struct {
	CartID    uint    `gorm:"primaryKey" json:"cart_id"`
	ProductID uint    `gorm:"primaryKey" json:"product_id"`
	Quantity  int     `gorm:"not null" json:"quantity"`
	Product   Product `json:"product"`
}
