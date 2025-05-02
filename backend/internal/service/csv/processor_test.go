package processor_test

import (
	"context"
	"file-uploader/database/model"
	"file-uploader/database/repository"
	processor "file-uploader/internal/service/csv"
	Seeder "file-uploader/internal/service/csv/seeder"
	testutils "file-uploader/internal/test-utils"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var testDB *gorm.DB
var studentRepo repository.StudentRepository[model.StudentTest]

func TestMain(m *testing.M) {
	err := godotenv.Load("../../../.env")
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

func TestCSVProcessor(t *testing.T) {
	const (
		testFilesDir  = "/tmp/testDir/"
		recordsLength = 10
		batchSize     = 1000
	)

	ctxTest := context.TODO()

	t.Run("a valid local csv file", func(t *testing.T) {

		filepath, err := Seeder.SeedStudentsCSV("test.csv", testFilesDir, recordsLength)
		require.NoError(t, err)

		f, err := os.Open(filepath)
		require.NoError(t, err)
		defer f.Close()

		fileStat, err := f.Stat()
		require.NoError(t, err)

		status := make(chan processor.ProcessStatus)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = processor.ProcessCSV(ctxTest, 0, f, fileStat.Size(), batchSize, studentRepo, StudentTestMapper, status)
			require.NoError(t, err)
			close(status)
		}()

		for stat := range status {
			t.Logf("Progress: %.2f%%, Time left: %.2fs", stat.Percent, stat.Timeleft)
		}
		wg.Wait()
	})

	t.Run("a valid file sent by a post request", func(t *testing.T) {
		_, fileHeaders := testutils.PrepareUploadTestFiles(t, []string{"class1.csv"})
		src, err := fileHeaders[0].Open()
		require.NoError(t, err)
		defer src.Close()

		size := fileHeaders[0].Size
		t.Logf("file size: %d", size)

		status := make(chan processor.ProcessStatus)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			err = processor.ProcessCSV(ctxTest, 0, src, size, batchSize, studentRepo, StudentTestMapper, status)
			require.NoError(t, err)
			close(status)
		}()

		for stat := range status {
			t.Logf("Progress: %.2f%%, Time left: %.2fs", stat.Percent, stat.Timeleft)
		}
		wg.Wait()
	})
}

func StudentTestMapper(record []string) (*model.StudentTest, error) {
	studentID, err := uuid.Parse(record[0])
	if err != nil {
		return nil, fmt.Errorf("error parsing student_id: %v", err)
	}
	grade, err := strconv.Atoi(record[3])
	if err != nil {
		return nil, fmt.Errorf("error parsing grade: %v", err)
	}
	return &model.StudentTest{
		Student_id:   studentID,
		Student_name: record[1],
		Subject:      record[2],
		Grade:        uint(grade),
	}, nil
}
