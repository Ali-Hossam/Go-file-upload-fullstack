package api

import (
	"file-uploader/database/model"
	"file-uploader/database/repository"
	"file-uploader/internal/api/handler/students"
	"file-uploader/internal/api/handler/upload"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB) {
	// Create repositories
	studentRepo := repository.NewStudentRepository[model.Student](db)

	// Create handlers
	uploadHandler := upload.NewUploadHandler(&studentRepo)
	studentsHandler := students.NewHandler[model.Student](studentRepo)

	// Register routes
	apiGroup := e.Group("/api")
	apiGroup.POST("/upload", uploadHandler.Handle)

	apiGroup.GET("/students", studentsHandler.GetAll)
	apiGroup.GET("/students/:name", studentsHandler.GetByName)
	apiGroup.GET("/students/:subject", studentsHandler.FilterBySubject)
}
