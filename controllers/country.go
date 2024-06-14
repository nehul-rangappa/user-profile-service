package controllers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nehul-rangappa/gigawrks-user-service/models"
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

	response, err := http.Get("https://restcountries.com/v3.1/all")
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
	go c.countryStore.Create(countries)

	ctx.JSON(http.StatusOK, metaCountries)
	return
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
		if err != nil {
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
		if err != nil {
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errInvalidPathParam.Error()})
		return
	}

	result, err := c.countryStore.GetByID(cID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	countries = append(countries, *result)

	ctx.JSON(http.StatusOK, countries)
	return
}
