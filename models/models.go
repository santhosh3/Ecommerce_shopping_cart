package models

import (
	"gorm.io/gorm"
)

type User struct {
	ID              uint64 `json:"id" gorm:"primaryKey"`
	FirstName       string `json:"first_name" gorm:"first_name"`
	LastName        string `json:"last_name" gorm:"last_name"`
	Email           string `json:"email" gorm:"email"`
	ProfileImage    string `json:"profile" gorm:"profile"`
	Password        string `json:"password" gorm:"password"`
	PhoneNumber     string `json:"phone_number" gorm:"phone_number"`
	ShippingAddress []ShippingAddress
	BillingAddress  []BillingAddress
	gorm.Model
}

type ShippingAddress struct {
	ID      uint64 `json:"id" gorm:"primaryKey"`
	Street  string `json:"street" gorm:"street"`
	City    string `json:"city" gorm:"city"`
	Pincode string `json:"pincode" gorm:"pincode"`
	UserID  uint64 `gorm:"index"`
	User    User
	gorm.Model
}

type BillingAddress struct {
	ID      uint64 `json:"id" gorm:"primaryKey"`
	Street  string `json:"street" gorm:"street"`
	City    string `json:"city" gorm:"city"`
	Pincode string `json:"pincode" gorm:"pincode"`
	UserID  uint64 `gorm:"index"`
	User    User
	gorm.Model
}

func DBMigrations(db *gorm.DB)  {
	db.AutoMigrate(
		&BillingAddress{},
		&ShippingAddress{},
		&User{},
	)
}