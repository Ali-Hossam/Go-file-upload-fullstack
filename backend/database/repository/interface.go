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
	nameField := value.FieldByName(string(config.Name))
	if !nameField.IsValid() || nameField.String() == "" {
		return uuid.Nil, config.ErrMissingStudentData
	}

	// Get subject field and check if it's empty (AI)
	subjectField := value.FieldByName(string(config.Subject))
	if !subjectField.IsValid() || subjectField.String() == "" {
		return uuid.Nil, config.ErrMissingStudentData
	}

	// Get grade field and check if it's zero (AI)
	gradeField := value.FieldByName(string(config.Grade))
	if !gradeField.IsValid() || gradeField.Uint() == 0 {
		return uuid.Nil, config.ErrMissingStudentData
	}

	// Get id field and check if it's nil
	idField := value.FieldByName(string(config.Id))
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
	result := r.db.Where(string(config.Name)+" = ?", name).Find(&students)

	if len(students) == 0 {
		return students, config.ErrStudentNotExist
	}

	return students, result.Error
}

func (r *StudentRepo[T]) GetAll(sortedBy config.StudentCol,
	sortOrder config.SortOrder,
	pageNumber int,
	pageSize int) ([]*T, error) {

	var students []*T
	db := r.db

	// Apply sorting if specified
	if sortedBy != "" {
		query := fmt.Sprintf("%s %s", sortedBy, sortOrder)
		db = db.Order(query)
	}

	// Apply pagination
	if pageNumber > 0 && pageSize > 0 {
		db = db.Limit(pageSize).Offset((pageNumber - 1) * pageSize)
	}

	// Execute query
	result := db.Find(&students)
	return students, result.Error
}

func (r *StudentRepo[T]) FilterBySubject(subject config.Course, pageNumber int, pageSize int) ([]*T, error) {
	var students []*T
	db := r.db

	if pageNumber > 0 && pageSize > 0 {
		db = db.Limit(pageSize).Offset((pageNumber - 1) * pageSize)
	}

	result := db.Where(string(config.Subject)+" = ?", subject).Find(&students)

	return students, result.Error
}

func (r *StudentRepo[T]) Delete(id uint) error {
	return nil
}
