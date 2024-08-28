package user

import (
	"database/sql"

	"github.com/santhosh3/ECOM/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(user *types.User) error {
	return nil
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	panic("unimplemented")
}