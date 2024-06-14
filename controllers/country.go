package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nehul-rangappa/gigawrks-user-service/models"
	"gorm.io/gorm"
)

// MetaCountry resource consisting of all the meta data attributes defining a country
type MetaCountry struct {
	Name struct {
		Common   string `json:"common"`
		Official string `json:"official"`
	}
	Cca2      string   `json:"cca2"`
	Capital   []string `json:"capital"`
	Region    string   `json:"region"`
	SubRegion string   `json:"subregion"`
}

type countryController struct {
	countryStore models.Countries
}

func NewCountryController(c models.Countries) *countryController {
	return &countryController{
		countryStore: c,
	}
}

// GetMetaCountries method takes a gin context
// interacts with the client API to fetch meta data of all countries
// saves data using model and writes back to the API response
func (c *countryController) GetMetaCountries(ctx *gin.Context) {
	metaCountries := make([]MetaCountry, 0)

	// Getting the HOST of rest countries from environment variable
	restCountiesHost := os.Getenv("REST_COUNTRIES_HOST")
	response, err := http.Get(restCountiesHost + "/v3.1/all")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer response.Body.Close()

	countryData, err := io.ReadAll(response.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// We are only using limited country information available from the external data
	if err1 := json.Unmarshal(countryData, &metaCountries); err1 != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err1.Error()})
		return
	}

	countries := make([]models.Country, 0, len(metaCountries))

	for i, mc := range metaCountries {
		countries = append(countries, models.Country{
			CommonName:   mc.Name.Common,
			OfficialName: mc.Name.Official,
			CountryCode:  mc.Cca2,
			Region:       mc.Region,
			SubRegion:    mc.SubRegion,
		})

		if len(mc.Capital) > 0 {
			countries[i].Capital = mc.Capital[0]
		}
	}

	// Go routine to handle creating country records in our database
	// Note: It can be ideal to wait for processing these records in Database.
	// But it depends on the use case, where the assumption here is not to worry about the DB records.
	go c.countryStore.Create(countries)

	ctx.JSON(http.StatusOK, metaCountries)
}

// GetCountries method takes a gin context
// checks for filters, validates them and
// interacts with the appropriate model fetch of all countries
// writes back to the API response
func (c *countryController) GetCountries(ctx *gin.Context) {
	countries := make([]models.Country, 0)

	countryCode := ctx.Query("code")
	if countryCode != "" {
		result, err := c.countryStore.GetByCode(countryCode)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		countries = append(countries, *result)

		ctx.JSON(http.StatusOK, countries)
		return
	}

	name := ctx.Query("name")
	if name != "" {
		result, err := c.countryStore.GetByName(name)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		} else if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		countries = append(countries, *result)

		ctx.JSON(http.StatusOK, countries)
		return
	}

	id := ctx.Query("id")

	if id == "" {
		countries, err := c.countryStore.GetAll()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, countries)
		return
	}

	cID, err := strconv.Atoi(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidPathParam.Error()})
		return
	}

	result, err := c.countryStore.GetByID(cID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	countries = append(countries, *result)

	ctx.JSON(http.StatusOK, countries)
}
