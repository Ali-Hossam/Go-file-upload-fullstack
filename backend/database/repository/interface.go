package repository

import (
	"file-uploader/database/config"
	"reflect"

	"gorm.io/gorm"
)

type StudentRepo[T any] struct {
	db *gorm.DB
}

func NewStudentRepository[T any](db *gorm.DB) StudentRepository[T] {
	return &StudentRepo[T]{db: db}
}

func (r *StudentRepo[T]) Create(item *T) (uint, error) {
	// Validate data
	value := reflect.ValueOf(item).Elem()

	// Get name field and check if it's empty (AI)
	nameField := value.FieldByName(config.StudentNameCol)
	if !nameField.IsValid() || nameField.String() == "" {
		return 0, config.ErrMissingStudentData
	}

	// Get subject field and check if it's empty (AI)
	subjectField := value.FieldByName(config.StudentSubjectCol)
	if !subjectField.IsValid() || subjectField.String() == "" {
		return 0, config.ErrMissingStudentData
	}

	// Get grade field and check if it's zero (AI)
	gradeField := value.FieldByName(config.StudentGradeCol)
	if !gradeField.IsValid() || gradeField.Uint() == 0 {
		return 0, config.ErrMissingStudentData
	}

	// Create student record
	result := r.db.Create(item)
	if result.Error != nil {
		return 0, result.Error
	}

	// Get student id
	idField := value.FieldByName(config.StudentIdCol)

	if !idField.IsValid() {
		return 0, config.ErrFieldNotFound
	}

	return uint(idField.Uint()), nil
}
