package models

type Users interface {
	GetByID(userID int) (*User, error)
	GetByEmail(email string) (*User, error)
	Create(user *User) (int, error)
	Update(user *User) error
	Delete(userID int) error
}

type Countries interface {
	GetAll() ([]Country, error)
	GetByID(id int) (*Country, error)
	GetByCode(countryCode string) (*Country, error)
	GetByName(name string) (*Country, error)
	Create(countries []Country) error
}
