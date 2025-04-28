package students

import "github.com/labstack/echo/v4"

type StudentsHandler interface {
	GetAll(c echo.Context) error
	GetByName(c echo.Context) error
	FilterBySubject(c echo.Context) error
}
