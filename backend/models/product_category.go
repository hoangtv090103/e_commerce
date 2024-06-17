package models

type ProductCategory struct {
	ProductID  uint `gorm:"primaryKey" json:"product_id"`
	CategoryID uint `gorm:"primaryKey" json:"category_id"`
}
