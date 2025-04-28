package testutils

import (
	"bytes"
	"file-uploader/config"
	"file-uploader/database"
	"file-uploader/database/model"
	"file-uploader/database/repository"
	Seeder "file-uploader/internal/service/csv/seeder"
	"fmt"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"net/textproto"
	"os"
	"path/filepath"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func PrepareUploadTestFiles(t *testing.T, files []string) (uploadDir string, fileHeaders []*multipart.FileHeader) {
	// Create a temp directory for our seed files
	seedDir, err := os.MkdirTemp("", "upload_test_seed")
	require.NoError(t, err)
	defer os.RemoveAll(seedDir)

	// Create a separate temp directory for the upload destination
	uploadDir, err = os.MkdirTemp("", "upload_test_dest")
	require.NoError(t, err)

	// Don't remove the upload directory so we can inspect the files
	// But log where they are
	t.Logf("Upload directory (for inspection): %s", uploadDir)

	for _, fileName := range files {
		// Seed a test CSV file in the seed directory
		path, err := Seeder.SeedStudentsCSV(fileName, uploadDir, 20)
		require.NoError(t, err)

		// Create a real multipart file header
		header, err := createMultipartFileHeader(path)
		require.NoError(t, err)
		fileHeaders = append(fileHeaders, header)
	}

	return
}

func createMultipartFileHeader(filePath string) (*multipart.FileHeader, error) {
	// Get the filename from the path
	filename := filepath.Base(filePath)

	// Open and read the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Create a buffer to write our multipart data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create a custom header with content-Type = text/csv
	header := textproto.MIMEHeader{}
	header.Set("Content-Disposition", fmt.Sprintf("form-data; name=\"file\"; filename=\"%s\"", filename))
	header.Set("Content-Type", "text/csv")

	// Create a form file
	part, err := writer.CreatePart(header)
	if err != nil {
		return nil, err
	}

	// Write the file content
	_, err = part.Write(fileContent)
	if err != nil {
		return nil, err
	}

	// Close the writer
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	// Parse the multipart form to get actual FileHeader
	mr := multipart.NewReader(body, writer.Boundary())
	form, err := mr.ReadForm(int64(body.Len()))
	if err != nil {
		return nil, err
	}

	// Use the real FileHeader created by the multipart package
	fileHeaders := form.File["file"]
	if len(fileHeaders) == 0 {
		return nil, os.ErrNotExist
	}

	fileHeaders[0].Header.Set("Content-Type", "text/csv")
	return fileHeaders[0], nil
}

func LoadDb() (*gorm.DB, repository.StudentRepository[model.StudentTest], error) {
	dsn, exist := os.LookupEnv(config.DBEnvVar)

	if !exist {
		return nil, nil, config.ErrEnvVarNotFound
	}

	db, repository, err := database.SetupDB[model.StudentTest](dsn, true)
	if err != nil {
		return nil, nil, err
	}

	return db, repository, nil
}

func NewTestContext(method, path string, body io.Reader) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, path, body)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}
