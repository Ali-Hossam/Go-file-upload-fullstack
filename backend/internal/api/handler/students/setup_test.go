package students_test

import (
	"file-uploader/database/model"
	"file-uploader/database/repository"
	"file-uploader/internal/api/handler/students"
	testutils "file-uploader/internal/test-utils"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var testDB *gorm.DB
var testStudentsHandler students.StudentsHandler
var testStudentsRepo repository.StudentRepository[model.StudentTest]

func TestMain(m *testing.M) {
	err := godotenv.Load("../../../../.env")
	if err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	db, repo, err := testutils.LoadDb()
	if err != nil {
		log.Fatalf("Failed to load test DB: %v", err)
	}

	testDB = db
	handler := students.NewHandler[model.StudentTest](repo)
	testStudentsHandler = handler
	testStudentsRepo = repo
	// Run tests
	code := m.Run()

	// Drop table after tests
	testDB.Migrator().DropTable(model.StudentTest{})

	// Cleanup
	sqlDB, _ := testDB.DB()
	sqlDB.Close()

	os.Exit(code)
}
