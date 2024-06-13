package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/nehul-rangappa/gigawrks-user-service/models"
)

type userController struct {
	userStore models.Users
}

func NewUserController(us models.Users) *userController {
	return &userController{
		userStore: us,
	}
}

// Signup method takes a gin context, validates the request body
// creates a hash of the password and interacts with the model
// creates a JWT token and writes back to the API response
func (u *userController) Signup(ctx *gin.Context) {
	return
}

// Login method takes a gin context, validates the request body
// validates the user credentials with existing information using model
// creates a JWT token and writes back to the API response
func (u *userController) Login(ctx *gin.Context) {
	return
}

// Get method takes a gin context, validates the path parameter
// authorizes the user based on JWT headers, interacts with the model
// to fetch user information and writes back to the API response
func (u *userController) Get(ctx *gin.Context) {
	return
}

// Update method takes a gin context, validates the path parameter, request body
// authorizes the user based on JWT headers, interacts with the model
// to update the existing user information and writes back to the API response
func (u *userController) Update(ctx *gin.Context) {
	return
}

// Delete method takes a gin context, validates the path parameter
// authorizes the user based on JWT headers, interacts with the model
// to delete the user information and writes back to the API response
func (u *userController) Delete(ctx *gin.Context) {
	return
}
