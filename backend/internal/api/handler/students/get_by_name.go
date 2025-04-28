package students

import (
	"file-uploader/config"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler[T]) GetByName(c echo.Context) error {
	name := c.Param("name")
	if name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, config.ErrMissingPathParamHttp)
	}

	records, err := h.Repo.GetByName(name)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, records)
}
