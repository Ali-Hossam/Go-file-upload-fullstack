package repository_test

import (
	"file-uploader/database"
	"file-uploader/database/config"
	"file-uploader/database/model"
	"file-uploader/database/repository"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCreate(t *testing.T) {
	cases := []struct {
		name          string
		studentData   *model.StudentTest
		expectedError error
	}{
		{
			name:        "valid students data",
			studentData: &model.StudentTest{Student_name: "test", Subject: "Chemistry", Grade: 20},
		},
		{
			name:          "invalid students data, should return error",
			studentData:   &model.StudentTest{Student_name: "incomplete data"},
			expectedError: config.ErrMissingStudentData,
		},
		{
			name:          "no students data, should return error",
			studentData:   nil,
			expectedError: config.ErrMissingStudentData,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			studentRepo := repository.NewStudentRepository[model.StudentTest](testDB)
			studentId, err := studentRepo.Create(tt.studentData)

			if tt.expectedError == nil && err == nil {
				defer testDB.Delete(&model.StudentTest{}, studentId)
			}

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
			assert.NotZero(t, studentId)
		})
	}
}

func loadDb() (*gorm.DB, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		return nil, err
	}
	dsn, exist := os.LookupEnv(config.DBEnvVar)

	if !exist {
		return nil, config.ErrEnvVarNotFound
	}

	db, err := database.SetupDB(dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
