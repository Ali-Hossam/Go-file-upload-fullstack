package upload

import (
	"context"
	"encoding/csv"
	"errors"
	"file-uploader/config"
	"file-uploader/database/model"
	"file-uploader/database/repository"
	processor "file-uploader/internal/service/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type UploadHandler struct {
	repo           *repository.StudentRepository[model.Student]
	statusChannels map[uuid.UUID]chan processor.ProcessStatus
	mu             sync.Mutex
}

func NewUploadHandler(repo *repository.StudentRepository[model.Student]) *UploadHandler {
	return &UploadHandler{
		repo:           repo,
		statusChannels: make(map[uuid.UUID]chan processor.ProcessStatus),
	}
}

// UploadHandler handles file uploads with websocket support
func (uh *UploadHandler) HandleFileUpload(c echo.Context) error {
	ctx := c.Request().Context()

	uploadID := uuid.New()

	// Create status channel for this upload
	uh.mu.Lock()
	statusChan := make(chan processor.ProcessStatus)
	uh.statusChannels[uploadID] = statusChan
	uh.mu.Unlock()

	// Extract necessary data from the request
	form, err := c.MultipartForm()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, config.ErrFormParseFailureHttp)
	}

	files := form.File["files"]
	if len(files) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, config.ErrNoFilesProvidedHttp)
	}

	// Save templ files
	var tempFiles []*os.File

	for _, fh := range files {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		src, err := fh.Open()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		defer src.Close()

		tmp, err := os.CreateTemp("", "upload-*.csv")
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		if _, err := io.Copy(tmp, src); err != nil {
			tmp.Close()
			os.Remove(tmp.Name())
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		tmp.Seek(0, io.SeekStart)
		tempFiles = append(tempFiles, tmp)
	}

	// Process file in the background
	go func(tempFiles []*os.File, uploadID uuid.UUID) {
		// Create a new background context that won't be canceled when the HTTP request ends
		bgCtx := context.Background()

		// Reclean temp files to ensure they are closed and removed
		defer func() {
			for _, f := range tempFiles {
				f.Close()
				os.Remove(f.Name())
			}
		}()

		err = ValidateCSVFiles(tempFiles)
		if err != nil {
			statusChan <- processor.ProcessStatus{Error: err.Error()}
			return
		}

		err = ValidateCSVHeader(tempFiles)
		if err != nil {
			statusChan <- processor.ProcessStatus{Error: err.Error()}
			return
		}

		ProcessFiles(bgCtx, tempFiles, uh.statusChannels[uploadID], *uh.repo, processor.StudentMapper)

	}(tempFiles, uploadID)

	return c.JSON(http.StatusOK, map[string]string{
		"upload_id": uploadID.String(),
	})
}

func (uh *UploadHandler) HandleStatusUpdates(c echo.Context) error {
	uploadID, err := uuid.Parse(c.Param("uploadID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	wsUpgrader := websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
	ws, err := wsUpgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer ws.Close()

	uh.mu.Lock()
	statusChan, exists := uh.statusChannels[uploadID]
	uh.mu.Unlock()

	if !exists {
		return echo.NewHTTPError(http.StatusNotFound, "upload ID not found")
	}

	// Send an initial ping so the client sees activity
	if err := ws.WriteJSON(processor.ProcessStatus{Percent: 0}); err != nil {
		return nil
	}

	// Use a separate done channel to detect client disconnection
	clientClosed := make(chan struct{})
	go func() {
		// Read loop - waits for any sign that the browser disconnected
		for {
			_, _, err := ws.ReadMessage()
			if err != nil {
				close(clientClosed)
				return
			}
		}
	}()

	processComplete := false
	for {
		select {
		case status, ok := <-statusChan:
			if !ok {
				// Channel closed, processing is complete
				processComplete = true
				goto cleanup
			}

			// Try to send the status update
			err := ws.WriteJSON(status)
			if err != nil {
				// The connection might be closed, but we'll keep processing
				goto cleanup
			}

		case <-clientClosed:
			// The client disconnected, but we don't mark processing as complete
			goto cleanup
		}
	}

cleanup:
	// Only clean up the channel if processing is complete (to avoid breaking other listeners)
	if processComplete {
		uh.mu.Lock()
		delete(uh.statusChannels, uploadID)
		uh.mu.Unlock()
	}
	return nil
}

func ProcessFiles[T any](
	ctx context.Context,
	files []*os.File,
	statusChannel chan processor.ProcessStatus,
	studentRepo repository.StudentRepository[T],
	mapper func([]string) (*T, error),
) {
	defer close(statusChannel)
	const (
		batchSize             = 2000
		maxNumberOfGoRoutines = 10
	)

	var sem = make(chan struct{}, min(maxNumberOfGoRoutines, len(files))) // limit number of goroutines

	var wg sync.WaitGroup

	for i, file := range files {
		wg.Add(1)
		go func(i int, file *os.File) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			file.Seek(0, io.SeekStart)
			fileInfo, err := file.Stat()
			if err != nil {
				statusChannel <- processor.ProcessStatus{Id: i, Error: fmt.Sprintf("Error getting file info: %v", err)}
				return
			}

			err = processor.ProcessCSV(
				ctx,
				i,
				file,
				fileInfo.Size(),
				batchSize,
				studentRepo,
				mapper,
				statusChannel,
			)
			if err != nil && err != context.Canceled && err != context.DeadlineExceeded {
				statusChannel <- processor.ProcessStatus{Id: i, Percent: 0, Error: fmt.Sprintf("Processing failed: %v", err)}
			}
		}(i, file)
	}

	wg.Wait()
}

func ValidateCSVFiles(files []*os.File) error {
	for _, f := range files {
		buf := make([]byte, 512)
		n, _ := f.Read(buf)
		f.Seek(0, io.SeekStart)

		contentType := http.DetectContentType(buf[:n])
		if !strings.Contains(contentType, "csv") && !strings.Contains(contentType, "text/plain") {
			return fmt.Errorf("invalid content type: %s", contentType)
		}
	}
	return nil
}

func ValidateCSVHeader(files []*os.File) error {
	for _, f := range files {
		f.Seek(0, io.SeekStart)
		reader := csv.NewReader(f)
		header, err := reader.Read()
		if err != nil {
			return err
		}
		if strings.Join(header, ",") != config.StudentsTableHeader {
			return errors.New("invalid CSV header")
		}
	}
	return nil
}
