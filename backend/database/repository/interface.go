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
}
