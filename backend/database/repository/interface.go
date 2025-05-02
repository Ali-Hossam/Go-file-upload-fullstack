package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Query option defines filter/sort/pagination actions
type QueryOption func(*gorm.DB) *gorm.DB

type StudentRepository[T any] interface {
	Create(item *T) (uuid.UUID, error)
	CreateMany(item []*T) error
	Query(opts []QueryOption, paginationOpt QueryOption) ([]*T, int64, error)
}
