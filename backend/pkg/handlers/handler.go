package handlers

import "gorm.io/gorm"

type handler struct {
	DB *gorm.DB
}

func NewHandler(db *gorm.DB) *handler {
	return &handler{DB: db}
}
