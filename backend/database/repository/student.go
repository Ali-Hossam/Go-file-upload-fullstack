package repository

import "github.com/google/uuid"

type StudentRepository[T any] interface {
	Create(item *T) (uuid.UUID, error)
	CreateMany(item []*T) error
	GetByName(name string) ([]*T, error)
	GetAll(sortedBy string) ([]T, error)
	Delete(id uint) error
	FilterBySubject(subject string) ([]T, error)
}
