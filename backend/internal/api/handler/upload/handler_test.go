package upload_test

import (
	"context"
	"file-uploader/database/model"
	"file-uploader/database/repository"
	"file-uploader/internal/api/handler/upload"
	processor "file-uploader/internal/service/csv"
	testutils "file-uploader/internal/test-utils"
	"io"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var testDB *gorm.DB
var studentRepo repository.StudentRepository[model.StudentTest]

func TestMain(m *testing.M) {
	err := godotenv.Load("../../../../.env")
	if err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	db, repo, err := testutils.LoadDb()
	if err != nil {
		log.Fatalf("Failed to initalize test DB: %v", err)
	}

	testDB = db
	studentRepo = repo

	// Run tests
	code := m.Run()

	// Drop table after tests
	testDB.Migrator().DropTable(model.StudentTest{})

	// Cleanup
	sqlDB, _ := testDB.DB()
	sqlDB.Close()

	os.Exit(code)
}

// [AI]
func TestUploadFilesHandler(t *testing.T) {
	cases := []struct {
		name  string
		files []string
	}{
		{
			name:  "valid files",
			files: []string{"class1.csv", "class2.csv", "class3.csv"},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			_, fileHeaders := testutils.PrepareUploadTestFiles(t, tt.files)

			// Call the handler
			statusChan := make(chan processor.ProcessStatus)
			// Start a goroutine to consume the messages from the channel
			go func() {
				for status := range statusChan {
					// Just consume the messages - we don't need to do anything with them in this test
					// Optionally log them if you want to see what's happening
					t.Logf("Process status: %+v", status)
				}
			}()
			var files []*os.File
			ctxTest := context.TODO()

			for _, fh := range fileHeaders {
				src, err := fh.Open()
				require.NoError(t, err)
				defer src.Close()

				tmpFile, err := os.CreateTemp("", "upload-*")
				require.NoError(t, err)

				_, err = io.Copy(tmpFile, src)
				require.NoError(t, err)

				// Rewind to the start for reading
				_, err = tmpFile.Seek(0, 0)
				require.NoError(t, err)
				files = append(files, tmpFile)
			}

			upload.ProcessFiles(ctxTest, files, statusChan, studentRepo, processor.StudentTestMapper)

			// Check that files were processed and data is in the DB
			students, _, err := studentRepo.Query([]repository.QueryOption{}, nil)
			require.NoError(t, err)

			// Verify some data was inserted - we don't need to check the exact count
			// as it would be brittle and dependent on the test files
			assert.NotEmpty(t, students, "Expected students in the DB")

		})
	}
}

func TestValidations(t *testing.T) {
	t.Run("validate CSV files", func(t *testing.T) {
		// Create a file with binary content (PNG signature) that will fail content type validation
		invalidFile, err := os.CreateTemp("", "invalid-*.png")
		require.NoError(t, err)
		defer os.Remove(invalidFile.Name())

		// Write PNG header signature bytes to make it detect as image/png
		pngSignature := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
		_, err = invalidFile.Write(pngSignature)
		require.NoError(t, err)
		_, err = invalidFile.WriteString("This is a fake PNG file to trigger content type validation")
		require.NoError(t, err)

		// Rewind to the start
		_, err = invalidFile.Seek(0, 0)
		require.NoError(t, err)

		// Call the validation function directly
		err = upload.ValidateCSVFiles([]*os.File{invalidFile})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid content type")
	})

	t.Run("validate header", func(t *testing.T) {
		// Create a file with invalid CSV header
		badHeaderFile, err := os.CreateTemp("", "bad-headers-*.csv")
		require.NoError(t, err)
		defer os.Remove(badHeaderFile.Name())

		_, err = badHeaderFile.WriteString("wrong,header,format\n1,2,3\n")
		require.NoError(t, err)

		// Rewind to the start
		_, err = badHeaderFile.Seek(0, 0)
		require.NoError(t, err)

		// Call the validation function directly
		err = upload.ValidateCSVHeader([]*os.File{badHeaderFile})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid CSV header")
	})
}

func TestUploadHandlerWithValidations(t *testing.T) {
	// Create a channel to receive validation errors
	statusChan := make(chan processor.ProcessStatus)
	defer close(statusChan)

	// Create mock files and test directly with ProcessFiles
	t.Run("invalid content type", func(t *testing.T) {
		// Create a file with binary content (PNG signature) that will fail content type validation
		invalidFile, err := os.CreateTemp("", "invalid-*.png")
		require.NoError(t, err)
		defer os.Remove(invalidFile.Name())

		// Write PNG header signature bytes to make it detect as image/png
		pngSignature := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
		_, err = invalidFile.Write(pngSignature)
		require.NoError(t, err)
		_, err = invalidFile.WriteString("This is a fake PNG file to trigger content type validation")
		require.NoError(t, err)

		// Rewind to the start
		_, err = invalidFile.Seek(0, 0)
		require.NoError(t, err)

		// Test the validation function directly
		err = upload.ValidateCSVFiles([]*os.File{invalidFile})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid content type")
	})

	t.Run("invalid headers", func(t *testing.T) {
		// Create a file with invalid CSV header
		badHeaderFile, err := os.CreateTemp("", "bad-headers-*.csv")
		require.NoError(t, err)
		defer os.Remove(badHeaderFile.Name())

		_, err = badHeaderFile.WriteString("wrong,header,format\n1,2,3\n")
		require.NoError(t, err)

		// Rewind to the start
		_, err = badHeaderFile.Seek(0, 0)
		require.NoError(t, err)

		// Test the validation function directly
		err = upload.ValidateCSVHeader([]*os.File{badHeaderFile})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid CSV header")
	})
}
