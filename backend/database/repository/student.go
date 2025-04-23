package repository

import (
	"file-uploader/database/config"

	"github.com/google/uuid"
)

type StudentRepository[T any] interface {
	Create(item *T) (uuid.UUID, error)
	CreateMany(item []*T) error
	GetByName(name string) ([]*T, error)
	GetAll(sortedBy config.StudentCol, sortOrder config.SortOrder, pageNumber, pageSize int) ([]*T, error)
	Delete(id uint) error
	FilterBySubject(subject string) ([]T, error)
}
