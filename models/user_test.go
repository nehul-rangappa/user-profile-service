package models

import (
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Test_userStore_GetByID runs unit tests on the method GetByID
func Test_userStore_GetByID(t *testing.T) {
	fDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error '%v' when opening a mock database connection", err)
	}
	defer fDB.Close()

	tests := []struct {
		name    string
		userID  int
		mock    func()
		want    *User
		wantErr error
	}{
		{
			name:   "Success case",
			userID: 1,
			mock: func() {
				versionRows := sqlmock.NewRows([]string{"version"}).AddRow("1")
				mock.ExpectQuery("SELECT VERSION").WillReturnRows(versionRows)
				rows := sqlmock.NewRows([]string{"id", "name", "country_id", "email", "password"}).
					AddRow(1, "Test User", 1, "test@gmail.com", "xasf2415g46")
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			want: &User{
				ID:        1,
				Name:      "Test User",
				CountryID: 1,
				Email:     "test@gmail.com",
				Password:  "xasf2415g46",
			},
			wantErr: nil,
		},
		{
			name:   "Failure case",
			userID: 1,
			mock: func() {
				versionRows := sqlmock.NewRows([]string{"version"}).AddRow("1")
				mock.ExpectQuery("SELECT VERSION").WillReturnRows(versionRows)
				mock.ExpectQuery("SELECT").WillReturnError(sqlmock.ErrCancelled)
			},
			wantErr: sqlmock.ErrCancelled,
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

			uS := NewUserStore(gormDB)

			got, err := uS.GetByID(tt.userID)
			if err != tt.wantErr {
				t.Errorf("userStore.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userStore.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_userStore_GetByEmail runs unit tests on the method GetByEmail
func Test_userStore_GetByEmail(t *testing.T) {
	fDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error '%v' when opening a mock database connection", err)
	}
	defer fDB.Close()

	tests := []struct {
		name    string
		email   string
		mockExp func()
		want    *User
		wantErr error
	}{
		{
			name:  "Success case",
			email: "test@gmail.com",
			mockExp: func() {
				versionRows := sqlmock.NewRows([]string{"version"}).AddRow("1")
				mock.ExpectQuery("SELECT VERSION").WillReturnRows(versionRows)
				rows := sqlmock.NewRows([]string{"id", "name", "country_id", "email", "password"}).
					AddRow(1, "Test User", 1, "test@gmail.com", "xasf2415g46")
				mock.ExpectQuery("SELECT").WillReturnRows(rows)
			},
			want: &User{
				ID:        1,
				Name:      "Test User",
				CountryID: 1,
				Email:     "test@gmail.com",
				Password:  "xasf2415g46",
			},
			wantErr: nil,
		},
		{
			name:  "Failure case",
			email: "test@gmail.com",
			mockExp: func() {
				versionRows := sqlmock.NewRows([]string{"version"}).AddRow("1")
				mock.ExpectQuery("SELECT VERSION").WillReturnRows(versionRows)
				mock.ExpectQuery("SELECT").WillReturnError(sqlmock.ErrCancelled)
			},
			wantErr: sqlmock.ErrCancelled,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExp()

			dialector := mysql.New(mysql.Config{
				Conn:       fDB,
				DriverName: "mysql",
			})
			gormDB, err := gorm.Open(dialector, &gorm.Config{})
			if err != nil {
				t.Fatalf("Error initializing gormDB: %v", err)
			}

			uS := NewUserStore(gormDB)

			got, err := uS.GetByEmail(tt.email)
			if err != tt.wantErr {
				t.Errorf("userStore.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userStore.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_userStore_Create runs unit tests on the method Create
func Test_userStore_Create(t *testing.T) {
	fDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error '%v' when opening a mock database connection", err)
	}
	defer fDB.Close()

	tests := []struct {
		name    string
		user    *User
		mockExp func()
		want    int
		wantErr error
	}{
		// Unable to test these cases due to time related attributes which is expected real time
		// {
		// 	name: "Success case",
		// 	user: &User{
		// 		ID:        1,
		// 		Name:      "Test User",
		// 		CountryID: 1,
		// 		Email:     "test@gmail.com",
		// 		Password:  "xasf2415g46",
		// 		CreatedAt: time.Time{},
		// 		UpdatedAt: time.Time{},
		// 	},
		// 	mockExp: func() {
		// 		versionRows := sqlmock.NewRows([]string{"version"}).AddRow("1")
		// 		mock.ExpectQuery("SELECT VERSION").WillReturnRows(versionRows)
		// 		mock.ExpectBegin()
		// 		mock.ExpectExec("INSERT").WithArgs("Test User", 1, "test@gmail.com", "xasf2415g46", time.Time{}, time.Time{}, 1).WillReturnResult(sqlmock.NewResult(1, 1))
		// 	},
		// 	want:    1,
		// 	wantErr: nil,
		// },
		// {
		// 	name: "Failure case",
		// 	user: &User{
		// 		ID:        1,
		// 		Name:      "Test User",
		// 		CountryID: 1,
		// 		Email:     "test@gmail.com",
		// 		Password:  "xasf2415g46",
		// 		CreatedAt: time.Time{},
		// 		UpdatedAt: time.Time{},
		// 	},
		// 	mockExp: func() {
		// 		versionRows := sqlmock.NewRows([]string{"version"}).AddRow("1")
		// 		mock.ExpectQuery("SELECT VERSION").WillReturnRows(versionRows)
		// 		mock.ExpectBegin()
		// 		mock.ExpectExec("INSERT").WithArgs("Test User", 1, "test@gmail.com", "xasf2415g46", time.Time{}, time.Time{}, 1).WillReturnError(sqlmock.ErrCancelled)
		// 	},
		// 	wantErr: sqlmock.ErrCancelled,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExp()

			dialector := mysql.New(mysql.Config{
				Conn:       fDB,
				DriverName: "mysql",
			})
			gormDB, err := gorm.Open(dialector, &gorm.Config{})
			if err != nil {
				t.Fatalf("Error initializing gormDB: %v", err)
			}

			uS := NewUserStore(gormDB)

			got, err := uS.Create(tt.user)
			if err != tt.wantErr {
				t.Errorf("userStore.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("userStore.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test_userStore_Update runs unit tests on the method Update
func Test_userStore_Update(t *testing.T) {
	fDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error '%v' when opening a mock database connection", err)
	}
	defer fDB.Close()

	tests := []struct {
		name    string
		user    *User
		mockExp func()
		wantErr error
	}{
		{
			name: "Failure GetByID case",
			user: &User{
				ID:        1,
				Name:      "Test User",
				CountryID: 1,
				Email:     "test@gmail.com",
				Password:  "xasf2415g46",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			mockExp: func() {
				versionRows := sqlmock.NewRows([]string{"version"}).AddRow("1")
				mock.ExpectQuery("SELECT VERSION").WillReturnRows(versionRows)
				mock.ExpectQuery("SELECT").WillReturnError(sqlmock.ErrCancelled)

			},
			wantErr: sqlmock.ErrCancelled,
		},
		// Unable to test these cases due to time related attributes which is expected real time
		// {
		// 	name: "Success case",
		// 	user: &User{
		// 		ID:        1,
		// 		Name:      "Test User",
		// 		CountryID: 1,
		// 		Email:     "test@gmail.com",
		// 		Password:  "xasf2415g46",
		// 		CreatedAt: time.Time{},
		// 		UpdatedAt: time.Time{},
		// 	},
		// 	mockExp: func() {
		// 		versionRows := sqlmock.NewRows([]string{"version"}).AddRow("1")
		// 		mock.ExpectQuery("SELECT VERSION").WillReturnRows(versionRows)
		// 		rows := sqlmock.NewRows([]string{"id", "name", "country_id", "email", "password"}).
		// 			AddRow(1, "Test User", 1, "test@gmail.com", "xasf2415g46")
		// 		mock.ExpectQuery("SELECT").WillReturnRows(rows)
		// 		mock.ExpectBegin()
		// 		mock.ExpectExec("UPDATE").WithArgs("Test User", 1, "test@gmail.com", "xasf2415g46", time.Time{}, time.Time{}, 1).WillReturnResult(sqlmock.NewResult(0, 1))
		// 	},
		// 	wantErr: nil,
		// },
		// {
		// 	name: "Failure case",
		// 	user: &User{
		// 		ID:        1,
		// 		Name:      "Test User",
		// 		CountryID: 1,
		// 		Email:     "test@gmail.com",
		// 		Password:  "xasf2415g46",
		// 		CreatedAt: time.Time{},
		// 		UpdatedAt: time.Time{},
		// 	},
		// 	mockExp: func() {
		// 		versionRows := sqlmock.NewRows([]string{"version"}).AddRow("1")
		// 		mock.ExpectQuery("SELECT VERSION").WillReturnRows(versionRows)
		// 		rows := sqlmock.NewRows([]string{"id", "name", "country_id", "email", "password"}).
		// 			AddRow(1, "Test User", 1, "test@gmail.com", "xasf2415g46")
		// 		mock.ExpectQuery("SELECT").WillReturnRows(rows)
		// 		mock.ExpectBegin()
		// 		mock.ExpectExec("UPDATE").WithArgs("Test User", 1, "test@gmail.com", "xasf2415g46", time.Time{}, time.Time{}, 1).WillReturnError(sqlmock.ErrCancelled)
		// 	},
		// 	wantErr: sqlmock.ErrCancelled,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockExp()

			dialector := mysql.New(mysql.Config{
				Conn:       fDB,
				DriverName: "mysql",
			})
			gormDB, err := gorm.Open(dialector, &gorm.Config{})
			if err != nil {
				t.Fatalf("Error initializing gormDB: %v", err)
			}

			uS := NewUserStore(gormDB)

			err = uS.Update(tt.user)
			if err != tt.wantErr {
				t.Errorf("userStore.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

// Test_userStore_Delete runs unit tests on the method Delete
func Test_userStore_Delete(t *testing.T) {
	fDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Unexpected error '%v' when opening a mock database connection", err)
	}
	defer fDB.Close()

	tests := []struct {
		name    string
		userID  int
		mock    func()
		wantErr error
	}{
		{
			name:   "Success case",
			userID: 1,
			mock: func() {
				versionRows := sqlmock.NewRows([]string{"version"}).AddRow("1")
				mock.ExpectQuery("SELECT VERSION").WillReturnRows(versionRows)
				mock.ExpectBegin()
				mock.ExpectExec("DELETE").WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			wantErr: nil,
		},
		{
			name:   "Failure case",
			userID: 1,
			mock: func() {
				versionRows := sqlmock.NewRows([]string{"version"}).AddRow("1")
				mock.ExpectQuery("SELECT VERSION").WillReturnRows(versionRows)
				mock.ExpectBegin()
				mock.ExpectExec("DELETE").WithArgs(1).WillReturnError(sqlmock.ErrCancelled)
				mock.ExpectRollback()
			},
			wantErr: sqlmock.ErrCancelled,
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

			uS := NewUserStore(gormDB)

			if err := uS.Delete(tt.userID); err != tt.wantErr {
				t.Errorf("userStore.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
