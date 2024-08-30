package user

import (
	"errors"
	"fmt"

	"github.com/santhosh3/ECOM/models"
	"github.com/santhosh3/ECOM/types"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(user *types.User) error {
	return nil
}

func (s *Store) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	// Use GORM's Where method to filter by email and First method to retrieve the first matching record
	err := s.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If no record is found, return a user not found error
			return nil, fmt.Errorf("user not found")
		}
		// Return any other errors encountered during the query
		return nil, err
	}
	// Return the found user
	return &user, nil
}