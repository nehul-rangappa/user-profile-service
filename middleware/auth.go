package middleware

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nehul-rangappa/gigawrks-user-service/controllers"
)

// verifyJWTToken takes a token and user ID
// validates the authenticity of the token followed by
// its validity based on expiration time and
// returns an error in case of any encountered issues
func verifyJWTToken(jwtToken string, id int) error {
	secretKey := os.Getenv("SECRET_KEY")
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
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

// Auth function is a middleware to authorize users to
// protected APIs before reaching the API handler
// It validates the path parameter, authorizes the user
// based on JWT token and verifies the ownership
// returns the API Handler Function if no error else
// writes back the response with the error message
func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.Param("id")
		if userID == "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": controllers.ErrMissingPathParam.Error()})
			return
		}

		id, err := strconv.Atoi(userID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": controllers.ErrInvalidPathParam.Error()})
			return
		}

		authHeaders := ctx.Request.Header["Authorization"]

		if len(authHeaders) == 0 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errors.New("missing Authorization Headers").Error()})
			return
		}

		authToken := strings.Split(authHeaders[0], " ")
		if len(authToken) != 2 {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": errors.New("invalid Authorization Headers").Error()})
			return
		}

		jwtToken := authToken[1]

		if err := verifyJWTToken(jwtToken, id); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx.Next()
	}
}
