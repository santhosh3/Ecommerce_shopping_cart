package user

import (
	"errors"
	"fmt"
	"sync"

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

func (s *Store) UpdateUserById(id uint64, userUpdates models.User) (*models.User, error) {
	var user models.User

	// Find the existing user by ID
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	// Update only the provided fields
	if err := s.db.Model(&user).Updates(userUpdates).Error; err != nil {
		return nil, err
	}

	// Return the updated user
	return &user, nil
}


func (s *Store) DeleteUserById(id uint64) (string, error) {
	err := s.db.Delete(&models.User{ID : id}).Error // Pass a pointer to `user` and handle the error
	if err != nil {
		return "", err
	}
	return "User deleted successfully", nil
}


func (s *Store) CreateUser(user models.User) (*models.User, error) {
	err := s.db.Create(&user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create the user")
	}
	return &user, nil
}

func (s *Store) CreateAddress(address types.Address) (*models.User, error) {
	var wg sync.WaitGroup
	var errChan = make(chan error,2) // Buffered channel to capture errors

	billing := &address.BillingAddress
	shipping := &address.ShippingAddress

	// Launch goroutine to create ShippingAddress
	wg.Add(1)
	go func ()  {
		defer wg.Done()
		if err := s.db.Create(&shipping).Error; err != nil {
			errChan <- err
		}	
	}()

	// Launch goroutine to create BillingAddress
	wg.Add(1)
	go func ()  {
		defer wg.Done()	
		if err := s.db.Create(&billing).Error; err != nil {
			errChan <- err
		}
	}()

	// Wait for both goroutines to finish
	wg.Wait()
	close(errChan)

	// Check if any errors occured
	if len(errChan) > 0 {
		return nil, <-errChan
	}
	
	var user models.User
	err := s.db.Preload("ShippingAddress").
		Preload("BillingAddress").
		First(&user, shipping.UserID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	// Use GORM's Where method to filter by email and First method to retrieve the first matching record
	err := s.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If no record is found, return a user not found error
			return nil, fmt.Errorf("user not found try to register")
		}
		// Return any other errors encountered during the query
		return nil, err
	}
	// Return the found user
	return &user, nil
}

func (s *Store) GetUserById(id int16) (*models.User, error) {
	var user models.User

	err := s.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	return &user, nil
}
