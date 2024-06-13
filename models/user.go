package models

import (
	"time"

	"gorm.io/gorm"
)

// User resource consisting of all the attributes defining a user
type User struct {
	ID        int       `json:"id" gorm:"primary_key"`
	Name      string    `json:"name"`
	Country   string    `json:"country"`
	Email     string    `json:"email"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type userStore struct {
	DB *gorm.DB
}

func NewUserStore(db *gorm.DB) Users {
	return &userStore{
		DB: db,
	}
}

// GetByID method takes a userID, fetches the user information
// from the database and returns User object along with an error if any
func (u *userStore) GetByID(userID int) (*User, error) {
	return nil, nil
}

// GetByEmail method takes an email, fetches the user information
// from the database and returns User object along with an error if any
func (u *userStore) GetByEmail(email string) (*User, error) {
	return nil, nil
}

// Create method takes a User object
// creates the user information in the database
// and returns the user ID along with an error if any
func (u *userStore) Create(user *User) (int, error) {
	return 0, nil
}

// Update method takes a User object
// updates the existing user information in the database
// and returns an error if any encountered
func (u *userStore) Update(user *User) error {
	return nil
}

// Delete method takes a user ID
// deletes the user information and
// return an error if encountered
func (u *userStore) Delete(userID int) error {
	return nil
}
