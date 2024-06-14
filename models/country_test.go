package models

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Test_countryStore_GetAll runs unit tests on the method GetAll
func Test_countryStore_GetAll(t *testing.T) {
	fDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error '%v' when opening a mock database connection", err)
	}
	defer fDB.Close()

	tests := []struct {
		name    string
		mock    func()
		want    []Country
		wantErr error
	}{
		{
			name: "Success case",
			mock: func() {
				versionRows := sqlmock.NewRows([]string{"version"}).AddRow("1")
				mock.ExpectQuery("SELECT VERSION").WillReturnRows(versionRows)
				rows := sqlmock.NewRows([]string{"id", "common_name", "official_name", "country_code", "capital", "region", "sub_region"}).
					AddRow(1, "United States", "United States of America", "US", "DC", "America", "North America").
					AddRow(2, "India", "Republic of India", "IN", "Delhi", "Asia", "South Asia")
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			want: []Country{
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
			},
			wantErr: nil,
		},
		{
			name: "Failure case",
			mock: func() {
				versionRows := sqlmock.NewRows([]string{"version"}).AddRow("1")
				mock.ExpectQuery("SELECT VERSION").WillReturnRows(versionRows)
				mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)
			},
			wantErr: sql.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			dialector := mysql.New(mysql.Config{
				Conn:       fDB,
				DriverName: "mysql",
			})
			gormDB, err := gorm.Open(dialector, &gorm.Config{})
			if err != nil {
				t.Fatalf("Error initializing gormDB: %v", err)
			}

			cS := NewCountryStore(gormDB)

			got, err := cS.GetAll()
			if err != tt.wantErr {
				t.Errorf("countryStore.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("countryStore.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_countryStore_GetByID runs unit tests on the method GetByID
func Test_countryStore_GetByID(t *testing.T) {
	fDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error '%v' when opening a mock database connection", err)
	}
	defer fDB.Close()

	tests := []struct {
		name    string
		id      int
		mock    func()
		want    *Country
		wantErr error
	}{
		{
			name: "Success case",
			id:   1,
			mock: func() {
				versionRows := sqlmock.NewRows([]string{"version"}).AddRow("1")
				mock.ExpectQuery("SELECT VERSION").WillReturnRows(versionRows)
				rows := sqlmock.NewRows([]string{"id", "common_name", "official_name", "country_code", "capital", "region", "sub_region"}).
					AddRow(1, "United States", "United States of America", "US", "DC", "America", "North America")
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			want: &Country{
				ID:           1,
				CommonName:   "United States",
				OfficialName: "United States of America",
				CountryCode:  "US",
				Capital:      "DC",
				Region:       "America",
				SubRegion:    "North America",
			},
			wantErr: nil,
		},
		{
			name: "Failure case",
			id:   1,
			mock: func() {
				versionRows := sqlmock.NewRows([]string{"version"}).AddRow("1")
				mock.ExpectQuery("SELECT VERSION").WillReturnRows(versionRows)
				mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)
			},
			wantErr: sql.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			dialector := mysql.New(mysql.Config{
				Conn:       fDB,
				DriverName: "mysql",
			})
			gormDB, err := gorm.Open(dialector, &gorm.Config{})
			if err != nil {
				t.Fatalf("Error initializing gormDB: %v", err)
			}

			cS := NewCountryStore(gormDB)

			got, err := cS.GetByID(tt.id)
			if err != tt.wantErr {
				t.Errorf("countryStore.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("countryStore.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_countryStore_GetByCode runs unit tests on the method GetByCode
func Test_countryStore_GetByCode(t *testing.T) {
	fDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error '%v' when opening a mock database connection", err)
	}
	defer fDB.Close()

	tests := []struct {
		name    string
		code    string
		mock    func()
		want    *Country
		wantErr error
	}{
		{
			name: "Success case",
			code: "US",
			mock: func() {
				versionRows := sqlmock.NewRows([]string{"version"}).AddRow("1")
				mock.ExpectQuery("SELECT VERSION").WillReturnRows(versionRows)
				rows := sqlmock.NewRows([]string{"id", "common_name", "official_name", "country_code", "capital", "region", "sub_region"}).
					AddRow(1, "United States", "United States of America", "US", "DC", "America", "North America")
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			want: &Country{
				ID:           1,
				CommonName:   "United States",
				OfficialName: "United States of America",
				CountryCode:  "US",
				Capital:      "DC",
				Region:       "America",
				SubRegion:    "North America",
			},
			wantErr: nil,
		},
		{
			name: "Failure case",
			code: "US",
			mock: func() {
				versionRows := sqlmock.NewRows([]string{"version"}).AddRow("1")
				mock.ExpectQuery("SELECT VERSION").WillReturnRows(versionRows)
				mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)
			},
			wantErr: sql.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			dialector := mysql.New(mysql.Config{
				Conn:       fDB,
				DriverName: "mysql",
			})
			gormDB, err := gorm.Open(dialector, &gorm.Config{})
			if err != nil {
				t.Fatalf("Error initializing gormDB: %v", err)
			}

			cS := NewCountryStore(gormDB)

			got, err := cS.GetByCode(tt.code)
			if err != tt.wantErr {
				t.Errorf("countryStore.GetByCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("countryStore.GetByCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_countryStore_GetByName runs unit tests on the method GetByName
func Test_countryStore_GetByName(t *testing.T) {
	fDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error '%v' when opening a mock database connection", err)
	}
	defer fDB.Close()

	tests := []struct {
		name       string
		commonName string
		mock       func()
		want       *Country
		wantErr    error
	}{
		{
			name:       "Success case",
			commonName: "United States",
			mock: func() {
				versionRows := sqlmock.NewRows([]string{"version"}).AddRow("1")
				mock.ExpectQuery("SELECT VERSION").WillReturnRows(versionRows)
				rows := sqlmock.NewRows([]string{"id", "common_name", "official_name", "country_code", "capital", "region", "sub_region"}).
					AddRow(1, "United States", "United States of America", "US", "DC", "America", "North America")
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			want: &Country{
				ID:           1,
				CommonName:   "United States",
				OfficialName: "United States of America",
				CountryCode:  "US",
				Capital:      "DC",
				Region:       "America",
				SubRegion:    "North America",
			},
			wantErr: nil,
		},
		{
			name:       "Failure case",
			commonName: "United States",
			mock: func() {
				versionRows := sqlmock.NewRows([]string{"version"}).AddRow("1")
				mock.ExpectQuery("SELECT VERSION").WillReturnRows(versionRows)
				mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)
			},
			wantErr: sql.ErrNoRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			dialector := mysql.New(mysql.Config{
				Conn:       fDB,
				DriverName: "mysql",
			})
			gormDB, err := gorm.Open(dialector, &gorm.Config{})
			if err != nil {
				t.Fatalf("Error initializing gormDB: %v", err)
			}

			cS := NewCountryStore(gormDB)

			got, err := cS.GetByName(tt.commonName)
			if err != tt.wantErr {
				t.Errorf("countryStore.GetByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("countryStore.GetByName() = %v, want %v", got, tt.want)
			}
		})
	}
}
