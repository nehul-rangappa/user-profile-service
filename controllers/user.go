package controllers

import (
	"errors"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nehul-rangappa/gigawrks-user-service/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userController struct {
	userStore models.Users
}

func NewUserController(us models.Users) *userController {
	return &userController{
		userStore: us,
	}
}

// validate function takes a User object and
// validates all the attributes and
// returns an error for any missing or invalid values
func validate(user *models.User) error {
	if user.Name == "" {
		return errors.New("user name cannot be empty")
	}

	if user.CountryID <= 0 {
		return errors.New("user's country cannot be empty")
	}

	if user.Email == "" || !regexp.MustCompile(`^[a-zA-Z0-9._]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(user.Email) {
		return errors.New("user email is empty or invalid")
	}

	if user.Password == "" || len(user.Password) < 8 {
		return errors.New("password should contain a minimum of 8 characters")
	}

	return nil
}

// createJWTToken function takes the userID
// uses the JWT to generate a token
// with an expiration period of 12 hours and
// returns the token along with any error
func createJWTToken(userID int) (string, error) {
	secretKey := os.Getenv("SECRET_KEY")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":     userID,
			"expiry": time.Now().Add(time.Hour * 12).Unix(), // Keeping an expiration period of 12 hours
		})

	jwtToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

// Signup method takes a gin context, validates the request body
// creates a hash of the password and interacts with the model
// creates a JWT token and writes back to the API response
func (u *userController) Signup(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errPayload.Error()})
		return
	}

	if err := validate(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.Password = string(hash)

	id, err1 := u.userStore.Create(&user)
	if err1 != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err1.Error()})
		return
	}

	jwtToken, err2 := createJWTToken(id)
	if err2 != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "issue while creating a jwt token"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"id": id, "jwtToken": jwtToken})
}

// Login method takes a gin context, validates the request body
// validates the user credentials with existing information using model
// creates a JWT token and writes back to the API response
func (u *userController) Login(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errPayload.Error()})
		return
	}

	if user.Email == "" || user.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing email or password"})
		return
	}

	userData, err := u.userStore.GetByEmail(user.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err1 := bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(user.Password)); err1 != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "credentials do not match. Please try again"})
		return
	}

	jwtToken, err2 := createJWTToken(userData.ID)
	if err2 != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "issue while creating a jwt token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": userData.ID, "jwtToken": jwtToken})
}

// Get method takes a gin context, validates the path parameter
// authorizes the user based on JWT headers, interacts with the model
// to fetch user information and writes back to the API response
func (u *userController) Get(ctx *gin.Context) {
	// Ignoring error as this is already validated in middleware
	id, _ := strconv.Atoi(ctx.Param("id"))

	userData, err := u.userStore.GetByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userData.Password = ""

	ctx.JSON(http.StatusOK, userData)
}

// Update method takes a gin context, validates the path parameter, request body
// authorizes the user based on JWT headers, interacts with the model
// to update the existing user information and writes back to the API response
func (u *userController) Update(ctx *gin.Context) {
	var user models.User

	// Ignoring error as this is already validated in middleware
	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errPayload.Error()})
		return
	}

	user.ID = id

	if err := validate(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.Password = string(hash)

	err1 := u.userStore.Update(&user)
	if err1 != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err1.Error()})
		return
	}

	user.Password = ""

	ctx.JSON(http.StatusOK, user)
}

// Delete method takes a gin context, validates the path parameter
// authorizes the user based on JWT headers, interacts with the model
// to delete the user information and writes back to the API response
func (u *userController) Delete(ctx *gin.Context) {
	// Ignoring error as this is already validated in middleware
	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := u.userStore.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
