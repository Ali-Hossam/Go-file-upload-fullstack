package students

import (
	"file-uploader/config"
	"file-uploader/database/repository"
	"net/http"

	"github.com/labstack/echo/v4"
)

type StudentsFilter struct {
	Page      int               `query:"page"`
	Size      int               `query:"size"`
	SortBy    config.StudentCol `query:"sort_by"`
	SortOrder config.SortOrder  `query:"sort_order"`
	Name      string            `query:"name"`
	Subject   config.Course     `query:"subject"`
}

const (
	DefaultPageSize = 100
	MaxPageSize     = 1000
)

// Validate subject
var validCourses = map[config.Course]bool{
	config.Mathematics: true,
	config.Physics:     true,
	config.Chemistry:   true,
	config.Biology:     true,
	config.History:     true,
	config.EnglishLit:  true,
	config.CompSci:     true,
	config.Art:         true,
	config.Music:       true,
	config.Geography:   true,
}

var validSortBys = map[config.StudentCol]bool{
	config.Id:      false,
	config.Name:    true,
	config.Grade:   true,
	config.Subject: true,
}

var validSortOrders = map[config.SortOrder]bool{
	config.SortAsc:  true,
	config.SortDesc: true,
}

// GetAll handles GET /students requests with filtering, sorting, and pagination
func (h *Handler[T]) GetAll(c echo.Context) error {
	var filter StudentsFilter
	if err := c.Bind(&filter); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid query parameters: "+err.Error())
	}

	// set defaults
	if filter.Page <= 0 {
		filter.Page = 1
	}

	if filter.Size <= 0 || filter.Size > MaxPageSize {
		filter.Size = DefaultPageSize
	}

	// Validate sort bys
	if filter.SortBy != "" {
		if !validSortBys[filter.SortBy] {
			return echo.NewHTTPError(http.StatusBadRequest, config.ErrInvalidFilterHttp)
		}
	}

	// Validate sort orders
	if filter.SortOrder != "" {
		if !validSortOrders[filter.SortOrder] {
			return echo.NewHTTPError(http.StatusBadRequest, config.ErrInvalidFilterHttp)
		}
	}

	if filter.Subject != "" {
		if !validCourses[filter.Subject] {
			return echo.NewHTTPError(http.StatusBadRequest, config.ErrInvalidFilterHttp)
		}
	}

	records, count, err := h.Repo.Query(
		[]repository.QueryOption{
			repository.WithNameFilter(filter.Name),
			repository.WithSubject(filter.Subject),
			repository.WithSort(filter.SortBy, filter.SortOrder),
		},
		repository.WithPagination(filter.Page, filter.Size),
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch records: "+err.Error())
	}

	return c.JSON(http.StatusOK, struct {
		Count   int64 `json:"count"`
		Records []*T  `json:"records"`
	}{
		Count:   count,
		Records: records,
	})
}
