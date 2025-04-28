package upload_test

import (
	"bytes"
	"file-uploader/config"
	"file-uploader/database/model"
	"file-uploader/database/repository"
	"file-uploader/internal/api/handler/upload"
	testutils "file-uploader/internal/test-utils"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
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
			e := echo.New()
			_, fileHeaders := testutils.PrepareUploadTestFiles(t, tt.files)

			// Create a buffer to write our multipart form data
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			// Add each file to the multipart form data
			for _, fileHeader := range fileHeaders {
				src, err := fileHeader.Open()
				require.NoError(t, err)

				// Create form file field with content type header
				h := make(textproto.MIMEHeader)
				h.Set("Content-Disposition",
					fmt.Sprintf(`form-data; name="files"; filename="%s"`, fileHeader.Filename))
				h.Set("Content-Type", "text/csv")

				part, err := writer.CreatePart(h)
				require.NoError(t, err)

				// Copy file content to form file field
				_, err = io.Copy(part, src)
				require.NoError(t, err)
				src.Close()

			}

			// Close the multipart writer
			err := writer.Close()
			require.NoError(t, err)

			// Create a request with the multipart form data
			req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
			req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())

			// Create a response recorder
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Call the handler
			err = upload.UploadHandlerWithoutWebsocket(c, studentRepo)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, rec.Result().StatusCode, "Expected status code to be %d", http.StatusOK)

		})
	}
}

// [AI]
func TestUploadFileHandlerWithInvalidTypes(t *testing.T) {
	e := echo.New()

	t.Run("invalid mime types", func(t *testing.T) {
		// Create a buffer to write our multipart form data
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Add a file with incorrect mime type
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", `form-data; name="files"; filename="invalid.txt"`)
		h.Set("Content-Type", "text/plain") // Invalid mime type

		part, err := writer.CreatePart(h)
		require.NoError(t, err)

		// Write some dummy content
		part.Write([]byte("This is not a CSV file"))

		// Close the writer
		writer.Close()

		// Create a request with the multipart form data
		req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())

		// Create a response recorder
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Call the handler
		err = upload.UploadHandlerWithoutWebsocket(c, studentRepo)

		// Expect an error for invalid file type
		if assert.Error(t, err) {
			httpErr, ok := err.(*echo.HTTPError)
			assert.True(t, ok)
			assert.Equal(t, http.StatusBadRequest, httpErr.Code)
			assert.Contains(t, httpErr.Message, config.ErrInvalidFileTypeHttp)
		}
	})

	t.Run("invalid headers", func(t *testing.T) {
		// Create a buffer to write our multipart form data
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// Add file with valid mime type but invalid headers
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", `form-data; name="files"; filename="bad_headers.csv"`)
		h.Set("Content-Type", "text/csv")

		part, err := writer.CreatePart(h)
		require.NoError(t, err)

		// Write invalid CSV content (with wrong headers)
		part.Write([]byte("wrong,header,format\n1,2,3\n"))

		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
		req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Call the handler
		err = upload.UploadHandlerWithoutWebsocket(c, studentRepo)

		// Expect an error for invalid headers
		if assert.Error(t, err) {
			httpErr, ok := err.(*echo.HTTPError)
			assert.True(t, ok)
			assert.Equal(t, http.StatusBadRequest, httpErr.Code)
			assert.Contains(t, httpErr.Message, config.ErrInvalidCSVCols)
		}
	})

}
