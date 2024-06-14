package models

type Users interface {
	GetByID(userID int) (*User, error)
	GetByEmail(email string) (*User, error)
	Create(user *User) (int, error)
	Update(user *User) error
	Delete(userID int) error
}
