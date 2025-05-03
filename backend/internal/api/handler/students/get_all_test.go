package students_test

import (
	"encoding/json"
	"file-uploader/config"
	"file-uploader/database/model"
	Seeder "file-uploader/internal/service/csv/seeder"
	testutils "file-uploader/internal/test-utils"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAll(t *testing.T) {
	// Create test data
	data := []model.StudentTest{
		Seeder.CreateStudentRecord(),
		Seeder.CreateStudentRecord(),
		Seeder.CreateStudentRecord(),
		Seeder.CreateStudentRecord(),
		Seeder.CreateStudentRecord(),
	}

	for _, student := range data {
		testStudentsRepo.Create(&student)
	}

	t.Run("all data without filters", func(t *testing.T) {
		c, rec := testutils.NewTestContext(http.MethodGet, "/students", nil)

		// Call the handler
		if assert.NoError(t, testStudentsHandler.GetAll(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)

			// unmarshal response
			var response struct {
				Count   int64               `json:"count"`
				Records []model.StudentTest `json:"records"`
			}

			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, data, response.Records)
			assert.Equal(t, int64(len(data)), response.Count)

		}
	})

	t.Run("all data with filters", func(t *testing.T) {
		cases := []struct {
			name          string
			filters       string
			expectedError string
		}{
			{
				name:    "valid columns filter, should have no errors",
				filters: fmt.Sprintf("?page=1&size=2&sort_by=Subject&sort_order=asc"),
			},
			{
				name:          "invalid filter value, should return error",
				filters:       fmt.Sprintf("?sort_by=not_a_valid_sort"),
				expectedError: config.ErrInvalidFilterHttp,
			},
			{
				name:    "invalid filter, should pass",
				filters: fmt.Sprintf("?iam_invalid=let_me_in"),
			},
		}

		for _, tt := range cases {
			t.Run(tt.name, func(t *testing.T) {
				c, _ := testutils.NewTestContext(http.MethodGet, "/students"+tt.filters, nil)

				// Call the handler
				err := testStudentsHandler.GetAll(c)
				if tt.expectedError != "" {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}
