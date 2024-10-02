package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Product struct {
	ID             uint64         `json:"id,omitempty" gorm:"primaryKey"`
	Title          string         `json:"title,omitempty" validate:"required" gorm:"unique"`
	Description    string         `json:"description,omitempty" validate:"required"`
	Price          float64        `json:"price,omitempty" validate:"required"`
	IsFreeShipping bool           `json:"is_free_shipping,omitempty" gorm:"default:false"`
	ProductImage   string         `json:"product_image,omitempty" validate:"required"`
	CurrencyId     string         `json:"currency_id,omitempty"`
	Installments   int            `json:"installments,omitempty"`
	AvailableSize  datatypes.JSON `json:"available_size,omitempty"` // Use JSON type for GORM
	Quantity       int            `json:"quantity,omitempty" validate:"required"`
	gorm.Model
}

var ProductModel = []interface{}{
	Product{},
}
