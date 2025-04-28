package students_test

import (
	"encoding/json"
	"file-uploader/config"
	"file-uploader/database/model"
	testutils "file-uploader/internal/test-utils"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// [AI]
func TestGetByNameHandler(t *testing.T) {
	// Test cases
	tests := []struct {
		name           string
		studentName    string
		setupData      []*model.StudentTest
		expectedStatus int
		expectError    bool
		expectedCount  int
	}{
		{
			name:        "successful request - single student",
			studentName: "TestStudent123",
			setupData: []*model.StudentTest{
				{Student_id: uuid.New(), Student_name: "TestStudent123", Subject: string(config.Mathematics), Grade: 80},
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
			expectedCount:  1,
		},
		{
			name:        "successful request - multiple students with same name",
			studentName: "DuplicateName",
			setupData: []*model.StudentTest{
				{Student_id: uuid.New(), Student_name: "DuplicateName", Subject: string(config.Mathematics), Grade: 80},
				{Student_id: uuid.New(), Student_name: "DuplicateName", Subject: string(config.Physics), Grade: 90},
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
			expectedCount:  2,
		},
		{
			name:           "empty name parameter",
			studentName:    "",
			setupData:      []*model.StudentTest{},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear previous test data
			testDB.Where("1 = 1").Delete(&model.StudentTest{})

			// Setup test data
			if len(tt.setupData) > 0 {
				err := testStudentsRepo.CreateMany(tt.setupData)
				require.NoError(t, err)
			}

			// Create request context
			c, rec := testutils.NewTestContext(http.MethodGet, "/api/v1/students/"+tt.studentName, nil)
			c.SetParamNames("name")
			c.SetParamValues(tt.studentName)

			// Execute handler
			err := testStudentsHandler.GetByName(c)

			t.Log(rec.Body)
			t.Log(err)
			// Check response status
			if tt.expectError {
				// Write the proper status code and error message to the response recorder.
				c.Echo().DefaultHTTPErrorHandler(err, c)
				assert.Error(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)

				// Parse response body for successful requests
				var response []model.StudentTest
				err := json.Unmarshal(rec.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response, tt.expectedCount)

				// Verify student names in response
				for _, student := range response {
					assert.Equal(t, tt.studentName, student.Student_name)
				}
			}
		})
	}
}
