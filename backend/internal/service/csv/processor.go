package processor

import (
	"encoding/csv"
	"file-uploader/config"
	"file-uploader/database/model"
	"file-uploader/database/repository"
	"fmt"
	"io"
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
	id int,
	file io.Reader,
	fileSize int64,
	batchSize int,
	studentRepo repository.StudentRepository[T],
	mapper RecordMapper[T],
	status chan ProcessStatus) error {

	countingReader := &CountingReader{R: file} // a wrapper to count bytes read
	reader := csv.NewReader(countingReader)

	buffer := make([]*T, 0, batchSize)
	startTime := time.Now()
	lastStatusUpdate := time.Now()

	// Read the csv file line by line
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		// Skip the header row
		if record[0] == string(config.Id) {
			continue
		}

		if err != nil {
			return fmt.Errorf("error reading csv file : %v", err)
		}

		// Map CSV record to struct
		entity, err := mapper(record)
		if err != nil {
			return fmt.Errorf("error mapping csv record :%v", err)
		}

		buffer = append(buffer, entity)

		// If buffer reaches batch size, insert into db
		if len(buffer) >= batchSize {
			if err := studentRepo.CreateMany(buffer); err != nil {
				return fmt.Errorf("error inserting batch : %v", err)
			}
			buffer = buffer[:0] // clear buffer
		}

		if time.Since(lastStatusUpdate) > 100*time.Millisecond {

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

	// Insert any remaining records in the buffer
	if len(buffer) > 0 {
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
