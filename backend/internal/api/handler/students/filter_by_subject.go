package students

import (
	"file-uploader/config"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler[T]) FilterBySubject(c echo.Context) error {
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

	subject := c.Param("subject")
	if subject == "" {
		return echo.NewHTTPError(http.StatusBadRequest, config.ErrMissingSearchParamHttp)
	}

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

	if !validCourses[config.Course(subject)] {
		return echo.NewHTTPError(http.StatusBadRequest, config.ErrInvalidSearchParamHttp)
	}

	records, err := h.Repo.FilterBySubject(config.Course(subject), filter.Page, filter.Size)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch records: "+err.Error())
	}

	return c.JSON(http.StatusOK, records)

}
