package repository

type StudentRepository[T any] interface {
	Create(item *T) (uint, error)
	GetByName(name string) ([]T, error)
	GetAll(sortedBy string) ([]T, error)
	Delete(id uint) error
	FilterBySubject(subject string) ([]T, error)
}
