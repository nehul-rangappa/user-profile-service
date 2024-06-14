package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nehul-rangappa/gigawrks-user-service/controllers"
	"github.com/nehul-rangappa/gigawrks-user-service/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Failed to load env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to establish a connection with database")
	}

	userStore := models.NewUserStore(db)

	userController := controllers.NewUserController(userStore)

	app := gin.New()

	app.POST("/signup", userController.Signup)
	app.POST("/login", userController.Login)
	app.GET("/users/:id", userController.Get)
	app.PUT("/users/:id", userController.Update)
	app.DELETE("/users/:id", userController.Delete)

	app.Run("localhost:8000")
}
