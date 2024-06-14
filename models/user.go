package models

import (
	"time"

	"gorm.io/gorm"
)

// User resource consisting of all the attributes defining a user
type User struct {
	ID        int       `json:"id" gorm:"primary_key"`
	Name      string    `json:"name"`
	CountryID int       `json:"countryID"`
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
	var user User
	if err := u.DB.First(&user, userID); err.Error != nil {
		return nil, err.Error
	}

	return &user, nil
}

// GetByEmail method takes an email, fetches the user information
// from the database and returns User object along with an error if any
func (u *userStore) GetByEmail(email string) (*User, error) {
	var user User
	if err := u.DB.Where("email = ?", email).First(&user); err.Error != nil {
		return nil, err.Error
	}

	return &user, nil
}

// Create method takes a User object
// creates the user information in the database
// and returns the user ID along with an error if any
func (u *userStore) Create(user *User) (int, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	result := u.DB.Create(user)

	if result.Error != nil {
		return 0, result.Error
	}

	return user.ID, nil
}

// Update method takes a User object
// updates the existing user information in the database
// and returns an error if any encountered
func (u *userStore) Update(user *User) error {
	existingUser, err := u.GetByID(user.ID)
	if err != nil {
		return err
	}

	user.CreatedAt = existingUser.CreatedAt
	user.UpdatedAt = time.Now()
	if result := u.DB.Save(user); result.Error != nil {
		return result.Error
	}

	return nil
}

// Delete method takes a user ID
// deletes the user information and
// return an error if encountered
func (u *userStore) Delete(userID int) error {
	if result := u.DB.Delete(&User{}, userID); result.Error != nil {
		return result.Error
	}

	return nil
}
