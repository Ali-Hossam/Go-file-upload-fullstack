package processor

import (
	"context"
	"encoding/csv"
	"file-uploader/database/model"
	"file-uploader/database/repository"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type RecordMapper[T any] func([]string) (*T, error)

type CountingReader struct {
	R io.Reader
	N int64
}

type ProcessStatus struct {
	Id       int
	Percent  float64
	Timeleft float64
	Error    string
}

func (c *CountingReader) Read(p []byte) (int, error) {
	n, err := c.R.Read(p)
	c.N += int64(n)
	return n, err
}

// ToRecordMapper creates a RecordMapper from any function with the same signature
// This helps resolve type compatibility issues when using generic functions [AI]
func ToRecordMapper[T any](fn func([]string) (*T, error)) RecordMapper[T] {
	return RecordMapper[T](fn)
}

func ProcessCSV[T any](
	ctx context.Context,
	id int,
	file io.Reader,
	fileSize int64,
	batchSize int,
	studentRepo repository.StudentRepository[T],
	mapper RecordMapper[T],
	status chan ProcessStatus) error {

	// Send initial status update immediately
	status <- ProcessStatus{
		Id:       id,
		Percent:  0,
		Timeleft: 0,
	}

	// Handle empty files
	if fileSize == 0 {
		status <- ProcessStatus{
			Id:       id,
			Percent:  100,
			Timeleft: 0,
			Error:    "File is empty",
		}
		return nil
	}

	countingReader := &CountingReader{R: file} // a wrapper to count bytes read
	reader := csv.NewReader(countingReader)

	// Read and skip the header row
	_, err := reader.Read()
	if err != nil {
		return fmt.Errorf("error reading CSV header: %v", err)
	}

	buffer := make([]*T, 0, batchSize)
	startTime := time.Now()
	lastStatusUpdate := time.Now()
	recordCount := 0

	// Read the csv file line by line
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("error reading csv file : %v", err)
		}

		recordCount++

		// Map CSV record to struct
		entity, err := mapper(record)
		if err != nil {
			return fmt.Errorf("error mapping csv record :%v", err)
		}

		buffer = append(buffer, entity)

		// If buffer reaches batch size, insert into db
		if len(buffer) >= batchSize {
			// Check context before database operation
			if ctxErr := ctx.Err(); ctxErr != nil {
				return ctxErr
			}
			if err := studentRepo.CreateMany(buffer); err != nil {
				return fmt.Errorf("error inserting batch : %v", err)
			}
			buffer = buffer[:0] // clear buffer
		}

		updateInterval := 100 * time.Millisecond
		if recordCount < 10 {
			updateInterval = 50 * time.Millisecond
		}

		if time.Since(lastStatusUpdate) > updateInterval {

			percent := (float64(countingReader.N) / float64(fileSize)) * 100
			elapsed := time.Since(startTime).Seconds()
			speed := float64(countingReader.N) / elapsed
			timeLeft := float64(fileSize-countingReader.N) / speed

			status <- ProcessStatus{
				Id:       id,
				Percent:  percent,
				Timeleft: timeLeft,
			}
			lastStatusUpdate = time.Now()
		}
	}

	// Check for context cancellation
	select {
	case <-ctx.Done():
		log.Printf("CSV processing cancelled: %v", ctx.Err())
		return ctx.Err()
	default:
	}

	// Insert any remaining records in the buffer
	if len(buffer) > 0 {
		// Check context before final database operation
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		if err := studentRepo.CreateMany(buffer); err != nil {
			return fmt.Errorf("error inserting final batch : %v", err)
		}
	}

	status <- ProcessStatus{
		Id:       id,
		Percent:  100,
		Timeleft: 0,
	}
	return nil
}

// MapStudentData creates a Student or StudentTest struct from a CSV record
func MapStudentData(record []string) (uuid.UUID, string, string, uint, error) {
	studentID, err := uuid.Parse(record[0])
	if err != nil {
		return uuid.Nil, "", "", 0, fmt.Errorf("error parsing student_id: %v", err)
	}

	grade, err := strconv.Atoi(record[3])
	if err != nil {
		return uuid.Nil, "", "", 0, fmt.Errorf("error parsing grade: %v", err)
	}

	return studentID, record[1], record[2], uint(grade), nil
}

// StudentMapper converts a CSV record to a Student model
func StudentMapper(record []string) (*model.Student, error) {
	id, name, subject, grade, err := MapStudentData(record)
	if err != nil {
		return nil, err
	}

	return &model.Student{
		Student_id:   id,
		Student_name: name,
		Subject:      subject,
		Grade:        grade,
	}, nil
}

// StudentTestMapper converts a CSV record to a StudentTest model
func StudentTestMapper(record []string) (*model.StudentTest, error) {
	id, name, subject, grade, err := MapStudentData(record)
	if err != nil {
		return nil, err
	}

	return &model.StudentTest{
		Student_id:   id,
		Student_name: name,
		Subject:      subject,
		Grade:        grade,
	}, nil
}
