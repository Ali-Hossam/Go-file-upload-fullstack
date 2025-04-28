package students

import (
	"file-uploader/config"
	"net/http"

	"github.com/labstack/echo/v4"
)

type StudentsFilter struct {
	Page      int               `query:"page"`
	Size      int               `query:"size"`
	SortBy    config.StudentCol `query:"sort_by"`
	SortOrder config.SortOrder  `query:"sort_order"`
}

func (h *Handler[T]) GetAll(c echo.Context) error {
	var filter StudentsFilter
	if err := c.Bind(&filter); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid query parameters: "+err.Error())
	}

	// set defaults
	if filter.Page <= 0 {
		filter.Page = 1
	}

	if filter.Size <= 0 || filter.Size > 100 {
		filter.Size = 20
	}

	// Validate sort bys
	if filter.SortBy != "" {
		validSortBys := map[config.StudentCol]bool{
			config.Id:      false,
			config.Name:    true,
			config.Grade:   true,
			config.Subject: true,
		}

		if !validSortBys[filter.SortBy] {
			return echo.NewHTTPError(http.StatusBadRequest, config.ErrInvalidFilterHttp)
		}
	}

	if filter.SortOrder != "" {
		// Validate sort orders
		validSortOrders := map[config.SortOrder]bool{
			config.SortAsc:  true,
			config.SortDesc: true,
		}

		if !validSortOrders[filter.SortOrder] {
			return echo.NewHTTPError(http.StatusBadRequest, config.ErrInvalidFilterHttp)
		}
	}

	records, err := h.Repo.GetAll(filter.SortBy, filter.SortOrder, filter.Page, filter.Size)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch records: "+err.Error())
	}

	return c.JSON(http.StatusOK, records)
}
