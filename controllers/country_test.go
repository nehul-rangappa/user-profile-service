package controllers

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/nehul-rangappa/gigawrks-user-service/models"
)

// Test_countryController_GetCountries runs unit tests on the method GetCountries
func Test_countryController_GetCountries(t *testing.T) {
	ctrl := gomock.NewController(t)
	countryModel := models.NewMockCountries(ctrl)

	tests := []struct {
		name       string
		id         string
		commonName string
		code       string
		expMock    func()
		wantCode   int
	}{
		{
			name: "Success case for Get All",
			expMock: func() {
				countryModel.EXPECT().GetAll().Return([]models.Country{
					{
						ID:           1,
						CommonName:   "United States",
						OfficialName: "United States of America",
						CountryCode:  "US",
						Capital:      "DC",
						Region:       "America",
						SubRegion:    "North America",
					},
					{
						ID:           2,
						CommonName:   "India",
						OfficialName: "Republic of India",
						CountryCode:  "IN",
						Capital:      "Delhi",
						Region:       "Asia",
						SubRegion:    "South Asia",
					},
				}, nil)
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Failure case due for Get All",
			expMock: func() {
				countryModel.EXPECT().GetAll().Return(nil, sql.ErrNoRows)
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			name:       "Success case for Get By Name",
			commonName: "United States",
			expMock: func() {
				countryModel.EXPECT().GetByName("United States").Return(&models.Country{
					ID:           1,
					CommonName:   "United States",
					OfficialName: "United States of America",
					CountryCode:  "US",
					Capital:      "DC",
					Region:       "America",
					SubRegion:    "North America",
				}, nil)
			},
			wantCode: http.StatusOK,
		},
		{
			name:       "Failure case due for Get By Name",
			commonName: "United States",
			expMock: func() {
				countryModel.EXPECT().GetByName("United States").Return(nil, sql.ErrNoRows)
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			name: "Success case for Get By Code",
			code: "US",
			expMock: func() {
				countryModel.EXPECT().GetByCode("US").Return(&models.Country{
					ID:           1,
					CommonName:   "United States",
					OfficialName: "United States of America",
					CountryCode:  "US",
					Capital:      "DC",
					Region:       "America",
					SubRegion:    "North America",
				}, nil)
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Failure case due for Get By Code",
			code: "US",
			expMock: func() {
				countryModel.EXPECT().GetByCode("US").Return(nil, sql.ErrNoRows)
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			name: "Success case for Get By ID",
			id:   "1",
			expMock: func() {
				countryModel.EXPECT().GetByID(1).Return(&models.Country{
					ID:           1,
					CommonName:   "United States",
					OfficialName: "United States of America",
					CountryCode:  "US",
					Capital:      "DC",
					Region:       "America",
					SubRegion:    "North America",
				}, nil)
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Failure case due for Get By ID",
			id:   "1",
			expMock: func() {
				countryModel.EXPECT().GetByID(1).Return(nil, sql.ErrNoRows)
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "Failure case due to wrong ID param",
			id:       "a",
			expMock:  func() {},
			wantCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.expMock()
			w := httptest.NewRecorder()
			gin.SetMode(gin.TestMode)

			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = &http.Request{
				Header: make(http.Header),
				URL:    &url.URL{},
			}
			ctx.Request.Method = "GET"

			u := url.Values{}

			u.Add("id", tt.id)
			u.Add("name", tt.commonName)
			u.Add("code", tt.code)

			ctx.Request.URL.RawQuery = u.Encode()

			c := NewCountryController(countryModel)

			c.GetCountries(ctx)

			if !reflect.DeepEqual(tt.wantCode, w.Code) {
				t.Errorf("userController.GetCountries() = %v, want %v", w.Code, tt.wantCode)
			}
		})
	}
}
