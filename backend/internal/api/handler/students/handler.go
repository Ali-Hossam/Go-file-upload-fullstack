package students

import "file-uploader/database/repository"

type Handler[T any] struct {
	Repo repository.StudentRepository[T]
}

func NewHandler[T any](repo repository.StudentRepository[T]) StudentsHandler {
	return &Handler[T]{
		Repo: repo,
	}
}
