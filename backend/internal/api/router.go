package api

import (
	"file-uploader/database/model"
	"file-uploader/database/repository"
	"file-uploader/internal/api/handler/students"
	"file-uploader/internal/api/handler/upload"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, studentsRepo repository.StudentRepository[model.Student]) {

	// Create handlers
	uploadHandler := upload.NewUploadHandler(&studentsRepo)
	studentsHandler := students.NewHandler[model.Student](studentsRepo)

	// Register routes
	apiGroup := e.Group("/api")
	apiGroup.POST("/upload", uploadHandler.HandleFileUpload)
	apiGroup.GET("/upload/status/:uploadID", uploadHandler.HandleStatusUpdates)

	apiGroup.GET("/students", studentsHandler.GetAll)
}
