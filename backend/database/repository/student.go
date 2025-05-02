package repository

import (
	"file-uploader/config"
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
	for i := range maxConcurrent {
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

func (r *StudentRepo[T]) Query(
	opts []QueryOption,
	paginationOpt QueryOption,
) ([]*T, int64, error) {

	db := r.db
	for _, opt := range opts {
		db = opt(db)
	}

	// Get total count before pagiantion
	var totalCount int64
	if err := db.Model(new(T)).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	if paginationOpt != nil {
		db = paginationOpt(db)
	}

	var students []*T
	result := db.Find(&students)
	return students, totalCount, result.Error
}

func WithSubject(subject config.Course) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if subject != "" {
			return db.Where(string(config.Subject)+" = ?", subject)
		}
		return db
	}
}

func WithSort(sortedBy config.StudentCol, order config.SortOrder) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if sortedBy != "" {
			return db.Order(fmt.Sprintf("%s %s", sortedBy, order))
		}
		return db
	}
}

func WithPagination(page, size int) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if page > 0 && size > 0 {
			return db.Offset((page - 1) * size).Limit(size)
		}
		return db
	}
}

func WithNameFilter(name string) QueryOption {
	return func(db *gorm.DB) *gorm.DB {
		if name != "" {
			return db.Where(string(config.Name)+" ILIKE ?", name+"%")
		}
		return db
	}
}
