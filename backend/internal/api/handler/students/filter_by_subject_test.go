package students_test

import (
	"encoding/json"
	"file-uploader/config"
	"file-uploader/database/model"
	Seeder "file-uploader/internal/service/csv/seeder"
	testutils "file-uploader/internal/test-utils"
	"net/http"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// [AI]
func TestFilterBySubject(t *testing.T) {
	// Clear previous test data
	testDB.Where("1 = 1").Delete(&model.StudentTest{})

	// Create test data with different subjects
	mathStudent := Seeder.CreateStudentRecord()
	mathStudent.Subject = string(config.Mathematics)

	physicsStudent := Seeder.CreateStudentRecord()
	physicsStudent.Subject = string(config.Physics)

	chemistryStudent := Seeder.CreateStudentRecord()
	chemistryStudent.Subject = string(config.Chemistry)

	biologyStudent := Seeder.CreateStudentRecord()
	biologyStudent.Subject = string(config.Biology)

	// Add test data to the database
	testStudentsRepo.Create(&mathStudent)
	testStudentsRepo.Create(&physicsStudent)
	testStudentsRepo.Create(&chemistryStudent)
	testStudentsRepo.Create(&biologyStudent)

	t.Run("successful filter by subject", func(t *testing.T) {
		// Create a request with subject as a path parameter
		c, rec := testutils.NewTestContext(http.MethodGet, "/students/subject/Mathematics", nil)
		c.SetParamNames("subject")
		c.SetParamValues(string(config.Mathematics))

		// Call the handler
		err := testStudentsHandler.FilterBySubject(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		// Unmarshal response
		var students []model.StudentTest
		err = json.Unmarshal(rec.Body.Bytes(), &students)
		assert.NoError(t, err)

		// Check that all returned students have Mathematics as subject
		for _, student := range students {
			assert.Equal(t, string(config.Mathematics), student.Subject)
		}
	})

	t.Run("error with missing subject parameter", func(t *testing.T) {
		// Create a request without the subject path parameter
		c, _ := testutils.NewTestContext(http.MethodGet, "/students/subject/", nil)

		// No subject parameter set

		// Call the handler
		err := testStudentsHandler.FilterBySubject(c)
		assert.Error(t, err)

		// Check that the error is due to missing path parameter
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
		assert.Equal(t, config.ErrMissingSearchParamHttp, httpError.Message)
	})

	t.Run("error with invalid subject parameter", func(t *testing.T) {
		// Create a request with an invalid subject path parameter
		c, _ := testutils.NewTestContext(http.MethodGet, "/students/subject/InvalidSubject", nil)
		c.SetParamNames("subject")
		c.SetParamValues("InvalidSubject")

		// Call the handler
		err := testStudentsHandler.FilterBySubject(c)
		assert.Error(t, err)

		// Check that the error is due to invalid search parameter
		httpError, ok := err.(*echo.HTTPError)
		assert.True(t, ok)
		assert.Equal(t, http.StatusBadRequest, httpError.Code)
		assert.Equal(t, config.ErrInvalidSearchParamHttp, httpError.Message)
	})

	t.Run("pagination works correctly", func(t *testing.T) {
		// Add more math students to test pagination
		for range 5 {
			student := Seeder.CreateStudentRecord()
			student.Subject = string(config.Mathematics)
			testStudentsRepo.Create(&student)
		}

		// Test first page
		c1, rec1 := testutils.NewTestContext(http.MethodGet, "/students/subject/Mathematics", nil)
		c1.SetParamNames("subject")
		c1.SetParamValues(string(config.Mathematics))
		q1 := c1.Request().URL.Query()
		q1.Add("page", "1")
		q1.Add("size", "2")
		c1.Request().URL.RawQuery = q1.Encode()

		err := testStudentsHandler.FilterBySubject(c1)
		assert.NoError(t, err)

		var firstPage []model.StudentTest
		err = json.Unmarshal(rec1.Body.Bytes(), &firstPage)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(firstPage), 1)

		// Test second page
		c2, rec2 := testutils.NewTestContext(http.MethodGet, "/students/subject/Mathematics", nil)
		c2.SetParamNames("subject")
		c2.SetParamValues(string(config.Mathematics))
		q2 := c2.Request().URL.Query()
		q2.Add("page", "2")
		q2.Add("size", "2")
		c2.Request().URL.RawQuery = q2.Encode()

		err = testStudentsHandler.FilterBySubject(c2)
		assert.NoError(t, err)

		var secondPage []model.StudentTest
		err = json.Unmarshal(rec2.Body.Bytes(), &secondPage)
		assert.NoError(t, err)

		// If we have enough data, pages should be different
		if len(firstPage) > 0 && len(secondPage) > 0 {
			assert.NotEqual(t, firstPage[0].Student_id, secondPage[0].Student_id)
		}
	})
}
