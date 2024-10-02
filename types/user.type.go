package types

import (
	"time"

	"github.com/santhosh3/ECOM/models"
)

type UserStore interface {
	GetUserByEmail(email string) (*models.User, error)
	CreateUser(userPayload models.User) (*RegisterUserResponse, error)
	GetUserById(id int16) (*models.User, error)
	CreateAddress(payload Address) (*models.User, error)
	DeleteUserById(id uint64) (string, error)
	InsertOTP(user models.User, otp string) error
	UpdateUserById(id uint64, userPayload models.User) (*models.User, error)
	RemoveOTP(user models.User) error
	LogOutUser(id int16) error
	LoggingUser(id uint64) error
	CheckOTPByEmail(email, otp string) (bool, error)
	UpdatePasswordByEmail(email, password string) (string, error)
}

type CheckOTPPayload struct {
	Email string `json:"email" validate:"required"`
	OTP   string `json:"otp" validate:"required"`
}

type UpdatePasswordCreds struct {
	Email           string `json:"email" validate:"required,min=8"`
	Password        string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required"`
}

type ForgetUserPassword struct {
	Email string `json:"email" validate:"required"`
}

type RefreshTokenPayload struct {
	Token string `json:"token"`
}

type RateLimitStruct struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

type Address struct {
	ShippingAddress models.ShippingAddress `json:"shippingAddress"`
	BillingAddress  models.BillingAddress  `json:"billingAddress"`
}

type RegisterUserResponse struct {
	ID           uint64 `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	ProfileImage string `json:"profile_image"`
	PhoneNumber  string `json:"phone_number"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
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
	UpdatedAt       time.Time `json:"updatedAt"`
}

type LoginUser struct {
	Email    string `json:"email" validate:"required,email"`
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

