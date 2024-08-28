package types

import "time"

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
}

type User struct {
	ID              int       `json:"id"`
	FirstName       string    `json:"firstname"`
	LastName        string    `json:"lastname"`
	Email           string    `json:"email"`
	ProfileImage    string    `json:"Image"`
	Password        string    `json:"-"`
	ShippingAddress string    `json:"shippingAddress"`
	BillingAddress  string    `json:"billingAddress"`
	CreatedAt       time.Time `json:"createdAt"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"-"`
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
