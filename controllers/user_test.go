package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/nehul-rangappa/gigawrks-user-service/models"
	"gorm.io/gorm"
)

func Test_userController_Signup(t *testing.T) {
	ctrl := gomock.NewController(t)
	userModel := models.NewMockUsers(ctrl)

	tests := []struct {
		name     string
		expMock  func()
		reqBody  models.User
		wantCode int
	}{
		{
			name:    "Failure case due to missing name",
			expMock: func() {},
			reqBody: models.User{
				CountryID: 1,
				Email:     "test@gmail.com",
				Password:  "xasf2415g46",
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name:    "Failure case due to missing country",
			expMock: func() {},
			reqBody: models.User{
				Name:     "Test User",
				Email:    "test@gmail.com",
				Password: "xasf2415g46",
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name:    "Failure case due to invalid email",
			expMock: func() {},
			reqBody: models.User{
				Name:      "Test User",
				CountryID: 1,
				Email:     "testuser",
				Password:  "xasf2415g46",
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name:    "Failure case due to invalid password",
			expMock: func() {},
			reqBody: models.User{
				Name:      "Test User",
				CountryID: 1,
				Email:     "test@gmail.com",
				Password:  "a",
			},
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
			ctx.Request.Method = "POST"

			jsonbytes, _ := json.Marshal(tt.reqBody)
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))

			uH := NewUserController(userModel)

			uH.Signup(ctx)

			if !reflect.DeepEqual(tt.wantCode, w.Code) {
				t.Errorf("userController.Signup() = %v, want %v", w.Code, tt.wantCode)
			}
		})
	}
}

func Test_userController_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	userModel := models.NewMockUsers(ctrl)

	tests := []struct {
		name     string
		expMock  func()
		reqBody  models.User
		wantCode int
	}{
		{
			name:    "Failure case due to missing data",
			expMock: func() {},
			reqBody: models.User{
				Email: "test@gmail.com",
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "Failure case due to Get By Email model",
			expMock: func() {
				userModel.EXPECT().GetByEmail("test@gmail.com").Return(nil, sql.ErrNoRows)
			},
			reqBody: models.User{
				Email:    "test@gmail.com",
				Password: "xasf2415g46",
			},
			wantCode: http.StatusInternalServerError,
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
			ctx.Request.Method = "POST"

			jsonbytes, _ := json.Marshal(tt.reqBody)
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))

			uH := NewUserController(userModel)

			uH.Login(ctx)

			if !reflect.DeepEqual(tt.wantCode, w.Code) {
				t.Errorf("userController.Login() = %v, want %v", w.Code, tt.wantCode)
			}
		})
	}
}

func Test_userController_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	userModel := models.NewMockUsers(ctrl)

	tests := []struct {
		name      string
		userID    int
		pathParam string
		expMock   func()
		wantCode  int
	}{
		{
			name:      "Success case",
			userID:    1,
			pathParam: "1",
			expMock: func() {
				userModel.EXPECT().GetByID(1).Return(&models.User{
					ID:        1,
					Name:      "Test User",
					CountryID: 1,
					Email:     "test@gmail.com",
					Password:  "xasf2415g46",
				}, nil)
			},
			wantCode: http.StatusOK,
		},
		{
			name:      "Failure case due to model",
			userID:    1,
			pathParam: "1",
			expMock: func() {
				userModel.EXPECT().GetByID(1).Return(nil, sql.ErrNoRows)
			},
			wantCode: http.StatusInternalServerError,
		},
		{
			name:      "Failure case due to no content found",
			userID:    1,
			pathParam: "1",
			expMock: func() {
				userModel.EXPECT().GetByID(1).Return(nil, gorm.ErrRecordNotFound)
			},
			wantCode: http.StatusNotFound,
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

			ctx.Params = []gin.Param{{Key: "id", Value: tt.pathParam}}

			uH := NewUserController(userModel)

			uH.Get(ctx)

			if !reflect.DeepEqual(tt.wantCode, w.Code) {
				t.Errorf("userController.Get() = %v, want %v", w.Code, tt.wantCode)
			}
		})
	}
}

func Test_userController_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	userModel := models.NewMockUsers(ctrl)

	tests := []struct {
		name      string
		userID    int
		pathParam string
		expMock   func()
		reqBody   models.User
		wantCode  int
	}{
		// {
		// 	name:      "Success case",
		// 	userID:    1,
		// 	pathParam: "1",
		// 	expMock: func() {
		// 		userModel.EXPECT().Update(&models.User{
		// 			ID:       1,
		// 			Name:     "Test User",
		// 			CountryID: 1,
		// 			Email:    "test@gmail.com",
		// 			Password: "xasf2415g46",
		// 		}).Return(nil)
		// 	},
		// 	reqBody: models.User{
		// 		ID:       1,
		// 		Name:     "Test User",
		// 		CountryID: 1,
		// 		Email:    "test@gmail.com",
		// 		Password: "xasf2415g46",
		// 	},
		// 	wantCode: http.StatusOK,
		// },
		{
			name:      "Failure case due to request body",
			userID:    1,
			pathParam: "a",
			expMock:   func() {},
			reqBody: models.User{
				ID:        1,
				Name:      "Test User",
				CountryID: 1,
				Email:     "",
				Password:  "",
			},
			wantCode: http.StatusBadRequest,
		},
		// {
		// 	name:      "Failure case due to model",
		// 	userID:    1,
		// 	pathParam: "1",
		// 	expMock: func() {
		// 		userModel.EXPECT().Update(&models.User{
		// 			ID:       1,
		// 			Name:     "Test User",
		// 			CountryID: 1,
		// 			Email:    "test@gmail.com",
		// 			Password: "xasf2415g46",
		// 		}).Return(sql.ErrNoRows)
		// 	},
		// 	reqBody: models.User{
		// 		ID:       1,
		// 		Name:     "Test User",
		// 		CountryID: 1,
		// 		Email:    "test@gmail.com",
		// 		Password: "xasf2415g46",
		// 	},
		// 	wantCode: http.StatusInternalServerError,
		// },
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
			ctx.Request.Method = "PUT"

			ctx.Params = []gin.Param{{Key: "id", Value: tt.pathParam}}

			jsonbytes, _ := json.Marshal(tt.reqBody)
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))

			uH := NewUserController(userModel)

			uH.Update(ctx)

			if !reflect.DeepEqual(tt.wantCode, w.Code) {
				t.Errorf("userController.Update() = %v, want %v", w.Code, tt.wantCode)
			}
		})
	}
}

func Test_userController_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	userModel := models.NewMockUsers(ctrl)

	tests := []struct {
		name      string
		userID    int
		pathParam string
		expMock   func()
		wantCode  int
	}{
		{
			name:      "Success case",
			userID:    1,
			pathParam: "1",
			expMock: func() {
				userModel.EXPECT().Delete(1).Return(nil)
			},
			wantCode: http.StatusNoContent,
		},
		{
			name:      "Failure case due to model",
			userID:    1,
			pathParam: "1",
			expMock: func() {
				userModel.EXPECT().Delete(1).Return(sql.ErrConnDone)
			},
			wantCode: http.StatusInternalServerError,
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
			ctx.Request.Method = "DELETE"

			ctx.Params = []gin.Param{{Key: "id", Value: tt.pathParam}}

			uH := NewUserController(userModel)

			uH.Delete(ctx)

			if !reflect.DeepEqual(tt.wantCode, w.Code) {
				t.Errorf("userController.Delete() = %v, want %v", w.Code, tt.wantCode)
			}
		})
	}
}
