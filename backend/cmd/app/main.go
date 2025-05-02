package main

import (
	"file-uploader/config"
	"file-uploader/database"
	"file-uploader/database/model"
	"file-uploader/internal/api"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load environment variables

	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get database connection string
	dsn, exist := os.LookupEnv(config.DBEnvVar)
	if !exist {
		log.Fatalf("Database connection string not found in environment variables")
	}

	// Get server post
	port, exist := os.LookupEnv(config.PortEnvVar)
	if !exist {
		log.Fatalf("Server port not found in environment variables")
	}

	db, studentsRepository, err := database.SetupDB[model.Student](dsn, false)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	e := echo.New()

	// Add CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	api.RegisterRoutes(e, db, studentsRepository)

	if err := e.Start(":" + port); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
