package upload

import (
	"encoding/csv"
	"errors"
	"file-uploader/config"
	"file-uploader/database/model"
	"file-uploader/database/repository"
	processor "file-uploader/internal/service/csv"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all connections
	},
}

type UploadHandler struct {
	repo *repository.StudentRepository[model.Student]
}

func NewUploadHandler(repo *repository.StudentRepository[model.Student]) *UploadHandler {
	return &UploadHandler{
		repo: repo,
	}
}

// UploadHandler handles file uploads with websocket support
func (uh *UploadHandler) Handle(c echo.Context) error {
	ws, err := wsUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer ws.Close()

	return processUpload(c, ws, *uh.repo, processor.StudentMapper)
}

// Generic processUpload function that works with both Student and StudentTest models
func processUpload[T any](c echo.Context, ws *websocket.Conn, studentRepo repository.StudentRepository[T], mapper func([]string) (*T, error)) error {
	const (
		batchSize             = 1000
		maxNumberOfGoRoutines = 10
	)

	form, err := c.MultipartForm()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, config.ErrFormParseFailureHttp)
	}

	files := form.File["files"]
	if len(files) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, config.ErrNoFilesProvidedHttp)
	}

	err = validateCSVType(files)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = validateHeader(files)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var sem = make(chan struct{}, min(maxNumberOfGoRoutines, len(files))) // limit number of goroutines
	var status = make(chan processor.ProcessStatus)

	var wg sync.WaitGroup
	var statusWg sync.WaitGroup

	// Start a separate goroutine to handle status messages
	statusWg.Add(1)
	go func() {
		defer statusWg.Done()
		for st := range status {
			if ws != nil {
				err := ws.WriteJSON(st)
				if err != nil {
					log.Println("Websocket Write error:", err)
					break
				}
				log.Println(st)
			}
		}
	}()

	for i, file := range files {

		wg.Add(1)
		go func(i int, file *multipart.FileHeader) {
			defer wg.Done()
			src, err := file.Open()
			if err != nil {
				log.Printf("file open error: %v", err)
				return
			}
			defer src.Close()

			sem <- struct{}{}
			defer func() { <-sem }()

			processor.ProcessCSV(
				i,
				src,
				file.Size,
				batchSize,
				studentRepo,
				mapper,
				status,
			)
		}(i, file)
	}

	wg.Wait()
	close(status)
	statusWg.Wait()

	return c.JSON(http.StatusOK, map[string]any{
		"message": "Files processed",
	})
}

// UploadHandlerWithoutWebsocket is used for testing
func UploadHandlerWithoutWebsocket(c echo.Context, repo repository.StudentRepository[model.StudentTest]) error {
	return processUpload(c, nil, repo, processor.StudentTestMapper)
}

func validateCSVType(files []*multipart.FileHeader) error {
	invalidFiles := make([]string, 0, len(files))
	for _, file := range files {
		mimeType := file.Header.Get("Content-Type")
		log.Print(mimeType)
		if mimeType != "text/csv" && mimeType != "application/vnd.ms-excel" {
			invalidFiles = append(invalidFiles, file.Filename)
		}
	}

	if len(invalidFiles) > 0 {
		errMsg := config.ErrInvalidFileTypeHttp + " : " + strings.Join(invalidFiles, ",")
		return errors.New(errMsg)
	}

	return nil
}

func validateHeader(files []*multipart.FileHeader) error {
	invalidFiles := make([]string, 0, len(files))

	for _, fh := range files {
		file, err := fh.Open()
		if err != nil {
			return err
		}

		reader := csv.NewReader(file)
		headerRow, err := reader.Read()
		file.Close()
		if err != nil {
			return err
		}
		if strings.Join(headerRow, ",") != config.StudentsTableHeader {
			invalidFiles = append(invalidFiles, fh.Filename)
		}
	}

	if len(invalidFiles) > 0 {
		errMsg := config.ErrInvalidCSVCols + " : " + strings.Join(invalidFiles, ",")
		return errors.New(errMsg)
	}

	return nil
}
