package models

type CartItem struct {
	ID 	  	  uint    `gorm:"primaryKey" json:"id"`
	Cart      Cart `json:"cart"`
	CartID    uint    `json:"cart_id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `gorm:"not null" json:"quantity"`
}
