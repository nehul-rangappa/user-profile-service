package models

import (
	"gorm.io/gorm"
)

// Country resource consisting of all the attributes defining a country
type Country struct {
	ID           int    `json:"id" gorm:"primary_key"`
	CommonName   string `json:"commonName"`
	OfficialName string `json:"officialName"`
	CountryCode  string `json:"countryCode" gorm:"unique"`
	Capital      string `json:"capital"`
	Region       string `json:"region"`
	SubRegion    string `json:"subregion"`
}

type countryStore struct {
	DB *gorm.DB
}

func NewCountryStore(db *gorm.DB) Countries {
	return &countryStore{
		DB: db,
	}
}

// GetAll method fetches the country information
// from the database and returns slice of Country object along with an error if any
func (c *countryStore) GetAll() ([]Country, error) {
	countries := make([]Country, 0)

	if err := c.DB.Find(&countries); err.Error != nil {
		return nil, err.Error
	}

	return countries, nil
}

// GetByID method takes a id, fetches the country information
// from the database and returns Country object along with an error if any
func (c *countryStore) GetByID(id int) (*Country, error) {
	var country Country
	if err := c.DB.First(&country, id); err.Error != nil {
		return nil, err.Error
	}

	return &country, nil
}

// GetByCode method takes a country code, fetches the country information
// from the database and returns Country object along with an error if any
func (c *countryStore) GetByCode(countryCode string) (*Country, error) {
	var country Country
	if err := c.DB.Where("country_code", countryCode).First(&country); err.Error != nil {
		return nil, err.Error
	}

	return &country, nil
}

// GetByID method takes a id, fetches the country information
// from the database and returns Country object along with an error if any
func (c *countryStore) GetByName(name string) (*Country, error) {
	var country Country
	if err := c.DB.Where("common_name", name).First(&country); err.Error != nil {
		return nil, err.Error
	}

	return &country, nil
}

// Create method takes a slice of Country object
// creates the user information in the database
// and returns the user ID along with an error if any
func (c *countryStore) Create(countries []Country) error {
	for _, country := range countries {
		result := c.DB.FirstOrCreate(&country, Country{CountryCode: country.CountryCode})
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}
