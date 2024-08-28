package types

import "time"

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
}

type RegisterUserPayload struct {
	FirstName    string `json:"firstname" validate:"required"`
	LastName     string `json:"lastname" validate:"required"`
	Email        string `json:"email" validate:"required"`
	ProfileImage string `json:"image" validate:"required,email"`
	Password     string `json:"password" validate:"required,min=3,max=130"`
	Address      struct {
		ShippingAddress ShippingAddressPayload `json:"shipping_address"`
		BillingAddress  BillingAddressPayload  `json:"billing_address"`
	} `json:"address"`
}

type ShippingAddressPayload struct {
	Street  string `json:"street" validate:"required"`
	City    string `json:"city" validate:"required"`
	Pincode string `json:"pincode" validate:"required"`
}

type BillingAddressPayload struct {
	Street  string `json:"street" validate:"required"`
	City    string `json:"city" validate:"required"`
	Pincode string `json:"pincode" validate:"required"`
}

type User struct {
	ID              int       `json:"id"`
	FirstName       string    `json:"firstname"`
	LastName        string    `json:"lastname"`
	Email           string    `json:"email"`
	ProfileImage    string    `json:"profile"`
	Password        string    `json:"-"` // Password will be ignored during JSON marshaling
	ShippingAddress string    `json:"shippingAddress"`
	BillingAddress  string    `json:"billingAddress"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt		 time.Time `json:"updatedAt"`
}

type LoginUser struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ShippingAddress struct {
	ID      int    `json:"id"`
	Street  string `json:"street"`
	City    string `json:"city"`
	Pincode string `json:"pincode"`
}

type BillingAddress struct {
	ID      int    `json:"id"`
	Street  string `json:"street"`
	City    string `json:"city"`
	Pincode string `json:"pincode"`
}
