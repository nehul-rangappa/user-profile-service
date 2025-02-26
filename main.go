package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nehul-rangappa/gigawrks-user-service/controllers"
	"github.com/nehul-rangappa/gigawrks-user-service/middleware"
	"github.com/nehul-rangappa/gigawrks-user-service/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables from config file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load env file")
	}

	// Read environment variables
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Open connection to MySQL with GORM
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=Local", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to establish a connection with database")
	}

	// Migrations can be used if needed but applying DDL from schema.sql is suggested for control on data types and attributes in schema
	// db.AutoMigrate(&models.Country{})
	// db.AutoMigrate(&models.User{})

	userStore := models.NewUserStore(db)
	countryStore := models.NewCountryStore(db)

	userController := controllers.NewUserController(userStore)
	countryController := controllers.NewCountryController(countryStore)

	// Initiate the app using GIN framework with default configuration
	app := gin.Default()

	// User APIs
	app.POST("/signup", userController.Signup)
	app.POST("/login", userController.Login)

	// Protected User APIs
	app.GET("/users/:id", middleware.Auth(), userController.Get)
	app.PUT("/users/:id", middleware.Auth(), userController.Update)
	app.DELETE("/users/:id", middleware.Auth(), userController.Delete)

	// Rest Country API
	app.GET("/rest-countries", countryController.GetMetaCountries)

	// Country API with filter support using query parameters id, code or name
	app.GET("/countries", countryController.GetCountries)

	// Start the server on port 8000
	app.Run("localhost:8000")
}
