package models

import (
	"gorm.io/gorm"
)

type User struct {
	ID              uint64            `json:"id" gorm:"primaryKey"`
	FirstName       string            `json:"first_name" validate:"required,min=2,max=100"`
	LastName        string            `json:"last_name" validate:"required,min=2,max=100"`
	Email           string            `json:"email" validate:"required,email"`
	ProfileImage    string            `json:"profile_image" validate:"omitempty,url"`
	Password        string            `json:"password" validate:"required,min=8"`
	PhoneNumber     string            `json:"phone_number" validate:"required,max=15"` 
	OTP             string            `json:"otp"`
	Status		 bool
	ShippingAddress []ShippingAddress `gorm:"foreignKey:UserID"`
	BillingAddress  []BillingAddress  `gorm:"foreignKey:UserID"`
	gorm.Model
}

type ShippingAddress struct {
	ID      uint64 `json:"id" gorm:"primaryKey"`
	Street  string `json:"street" validate:"required,min=2,max=100"`
	City    string `json:"city" validate:"required,min=2,max=100"`
	Pincode string `json:"pincode" validate:"required,min=2,max=100"`
	UserID  uint64 `gorm:"index"` 
	gorm.Model
}

type BillingAddress struct {
	ID      uint64 `json:"id" gorm:"primaryKey"`
	Street  string `json:"street" validate:"required,min=2,max=100"`
	City    string `json:"city" validate:"required,min=2,max=100"`
	Pincode string `json:"pincode" validate:"required,min=2,max=100"`
	UserID  uint64 `gorm:"index"` 
	gorm.Model
}

func DBMigrations(db *gorm.DB) {
	db.AutoMigrate(&User{}, &ShippingAddress{}, &BillingAddress{})
}
