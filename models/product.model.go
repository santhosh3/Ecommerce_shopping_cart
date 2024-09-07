package models

import (
	"gorm.io/gorm"
	"gorm.io/datatypes"
)


type Product struct {
	ID             uint64          `json:"id,omitempty" gorm:"primaryKey"`
	Title          string          `json:"title,omitempty" validate:"required" gorm:"unique"`
	Description    string          `json:"description,omitempty" validate:"required"`
	Price          float64         `json:"price,omitempty" validate:"required"`
	IsFreeShipping bool            `json:"is_free_shipping,omitempty" gorm:"default:false"`
	ProductImage   string          `json:"product_image,omitempty" validate:"required"`
	CurrencyId     string          `json:"currency_id,omitempty"`
	Installments   int             `json:"installments,omitempty"`
	AvailableSize  datatypes.JSON  `json:"available_size,omitempty"` // Use JSON type for GORM
	gorm.Model
}

var ProductModel = []interface{}{
	Product{},
}