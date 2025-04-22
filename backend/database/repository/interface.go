package repository

import (
	"file-uploader/database/config"
	"fmt"
	"reflect"
	"sync"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StudentRepo[T any] struct {
	db *gorm.DB
}

func NewStudentRepository[T any](db *gorm.DB) StudentRepository[T] {
	return &StudentRepo[T]{db: db}
}

func (r *StudentRepo[T]) Create(item *T) (uuid.UUID, error) {
	if item == nil {
		return uuid.Nil, config.ErrMissingStudentData
	}

	// Validate data
	value := reflect.ValueOf(item).Elem()

	// Get name field and check if it's empty (AI)
	nameField := value.FieldByName(config.StudentNameCol)
	if !nameField.IsValid() || nameField.String() == "" {
		return uuid.Nil, config.ErrMissingStudentData
	}

	// Get subject field and check if it's empty (AI)
	subjectField := value.FieldByName(config.StudentSubjectCol)
	if !subjectField.IsValid() || subjectField.String() == "" {
		return uuid.Nil, config.ErrMissingStudentData
	}

	// Get grade field and check if it's zero (AI)
	gradeField := value.FieldByName(config.StudentGradeCol)
	if !gradeField.IsValid() || gradeField.Uint() == 0 {
		return uuid.Nil, config.ErrMissingStudentData
	}

	// Get id field and check if it's nil
	idField := value.FieldByName(config.StudentIdCol)
	var studentId uuid.UUID

	if !idField.IsValid() || idField.Interface() == uuid.Nil {
		studentId = uuid.New()

		// Make sure the field is settable
		if idField.IsValid() && idField.CanSet() {
			idField.Set(reflect.ValueOf(studentId))
		} else {
			return uuid.Nil, fmt.Errorf("can't set ID field")
		}
	} else {
		studentId = idField.Interface().(uuid.UUID)
	}

	// Create student record
	result := r.db.Create(item)
	if result.Error != nil {
		return uuid.Nil, result.Error
	}

	return studentId, nil
}

func (r *StudentRepo[T]) CreateMany(items []*T) error {
	const (
		batchSize     = 500
		maxConcurrent = 10
	)

	if len(items) == 0 {
		return config.ErrMissingStudentData
	}

	// Calculate partition size for each worker
	partitionSize := (len(items) + maxConcurrent - 1) / maxConcurrent

	var wg sync.WaitGroup
	mu := sync.Mutex{}
	var errors []error

	// Assign paritions to workers
	for i := 0; i < maxConcurrent; i++ {
		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()

			// Calculate parition slice indices
			start := workerID * partitionSize
			if start >= len(items) {
				return
			}

			end := min(start+partitionSize, len(items))

			parition := items[start:end]

			// Process partition in batches
			result := r.db.CreateInBatches(parition, batchSize)

			if result.Error != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("worker %d batch error: %w", workerID, result.Error))
				mu.Unlock()
				return
			}

		}(i)
	}

	wg.Wait()
	if len(errors) > 0 {
		return fmt.Errorf("multiple errors occured: %v", errors)
	}
	return nil
}

func (r *StudentRepo[T]) GetByName(name string) ([]*T, error) {
	var students []*T
	query := fmt.Sprintf("%s = ?", config.StudentNameCol)
	result := r.db.Where(query, name).Find(&students)

	if len(students) == 0 {
		return students, config.ErrStudentNotExist
	}

	return students, result.Error
}
}
