package controllers

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nehul-rangappa/gigawrks-user-service/models"
	"golang.org/x/crypto/bcrypt"
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

	if user.Country == "" {
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":     userID,
			"expiry": time.Now().Add(time.Hour * 12).Unix(), // Keeping an expiration period of 12 hours
		})

	jwtToken, err := token.SignedString([]byte("gigawrks-secret-go-key"))
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

// verifyJWTToken takes a token and user ID
// validates the authenticity of the token followed by
// its validity based on expiration time and
// returns an error in case of any encountered issues
func verifyJWTToken(jwtToken string, id int) error {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		return []byte("gigawrks-secret-go-key"), nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid jwt token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("invalid jwt token")
	}

	if jwtID, ok := claims["id"].(float64); ok {
		if jwtID != float64(id) {
			return errors.New("no authorization to this entity")
		}
	}

	if expiry, ok := claims["expiry"].(float64); ok {
		if expiry < float64(time.Now().Unix()) {
			return errors.New("jwt token is expired")
		}
	}

	return nil
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
	return
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

	ctx.JSON(http.StatusOK, gin.H{"jwtToken": jwtToken})
	return
}

// authorizeUser function takes in a gin context
// validates the path parameter, authorizes the user
// based on JWT token and verifies the ownership
// returns an ID if no error, else it
// writes the response and returns the error
func authorizeUser(ctx *gin.Context) (int, error) {
	userID := ctx.Param("id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errMissingPathParam.Error()})
		return 0, errMissingPathParam
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errInvalidPathParam.Error()})
		return 0, errInvalidPathParam
	}

	authHeaders := ctx.Request.Header["Authorization"]

	if len(authHeaders) == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New("missing Authorization Headers").Error()})
		return 0, errors.New("missing Authorization Headers")
	}

	authToken := strings.Split(authHeaders[0], " ")
	if len(authToken) != 2 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": errors.New("invalid Authorization Headers").Error()})
		return 0, errors.New("invalid Authorization Headers")
	}

	jwtToken := authToken[1]

	if err := verifyJWTToken(jwtToken, id); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return 0, err
	}

	return id, nil
}

// Get method takes a gin context, validates the path parameter
// authorizes the user based on JWT headers, interacts with the model
// to fetch user information and writes back to the API response
func (u *userController) Get(ctx *gin.Context) {
	id, err := authorizeUser(ctx)
	if err != nil {
		return
	}

	userData, err := u.userStore.GetByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userData.Password = ""

	ctx.JSON(http.StatusOK, userData)
	return
}

// Update method takes a gin context, validates the path parameter, request body
// authorizes the user based on JWT headers, interacts with the model
// to update the existing user information and writes back to the API response
func (u *userController) Update(ctx *gin.Context) {
	var user models.User

	id, err := authorizeUser(ctx)
	if err != nil {
		return
	}

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
	return
}

// Delete method takes a gin context, validates the path parameter
// authorizes the user based on JWT headers, interacts with the model
// to delete the user information and writes back to the API response
func (u *userController) Delete(ctx *gin.Context) {
	id, err := authorizeUser(ctx)
	if err != nil {
		return
	}

	if err := u.userStore.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
	return
}
