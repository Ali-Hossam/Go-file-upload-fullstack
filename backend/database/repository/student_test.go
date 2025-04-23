package repository_test

import (
	"file-uploader/database"
	"file-uploader/database/config"
	"file-uploader/database/model"
	"file-uploader/database/repository"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var testDB *gorm.DB

func TestMain(m *testing.M) {
	db, err := loadDb()
	if err != nil {
		log.Fatalf("Failed to initalize test DB: %v", err)
	}

	testDB = db

	// Run tests
	code := m.Run()

	// Drop table after tests
	testDB.Migrator().DropTable(model.StudentTest{})

	// Cleanup
	sqlDB, _ := testDB.DB()
	sqlDB.Close()

	os.Exit(code)
}

func TestCreate(t *testing.T) {
	cases := []struct {
		name          string
		studentData   *model.StudentTest
		expectedError error
	}{
		{
			name:        "valid students data",
			studentData: &model.StudentTest{Student_name: "test", Subject: string(config.Chemistry), Grade: 20},
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

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				return
			}

			assert.NoError(t, err)
			assert.NotZero(t, studentId)
		})
	}
}

func TestCreateMany(t *testing.T) {
	cases := []struct {
		name          string
		studentsData  []*model.StudentTest
		expectedError error
	}{
		{
			name: "valid students data",
			studentsData: []*model.StudentTest{
				{Student_name: "Omar", Subject: string(config.Chemistry), Grade: 10},
				{Student_name: "Ali", Subject: string(config.CompSci), Grade: 20},
				{Student_name: "Saeed", Subject: string(config.EnglishLit), Grade: 30},
			},
		},
		{
			name:          "no students data, should return error",
			expectedError: config.ErrMissingStudentData,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			studentRepo := repository.NewStudentRepository[model.StudentTest](testDB)
			err := studentRepo.CreateMany(tt.studentsData)

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestGetByName(t *testing.T) {
	cases := []struct {
		name          string
		studentsData  []*model.StudentTest
		studentName   string
		expectedError error
	}{
		{
			name:         "get an existing student",
			studentsData: []*model.StudentTest{{Student_id: uuid.New(), Student_name: "unique_test_name", Subject: string(config.Mathematics), Grade: 80}},
			studentName:  "unique_test_name",
		},
		{
			name:          "get an non-existing student, should return error",
			studentsData:  []*model.StudentTest{{Student_id: uuid.New(), Student_name: "unique_test_name", Subject: string(config.Mathematics), Grade: 80}},
			studentName:   "iam not here",
			expectedError: config.ErrStudentNotExist,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			studentRepo := repository.NewStudentRepository[model.StudentTest](testDB)
			err := studentRepo.CreateMany(tt.studentsData)
			require.NoError(t, err)

			students, err := studentRepo.GetByName(tt.studentName)
			if tt.expectedError != nil {
				assert.Equal(t, tt.expectedError, err)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, students, len(tt.studentsData))

			// Verify all returned students have the expected name
			for _, student := range students {
				assert.Equal(t, tt.studentName, student.Student_name)
			}
		})
	}
}

// [AI]
func TestGetAllPagination(t *testing.T) {
	studentRepo := setupTestData(t)

	// Test cases for pagination
	testCases := []struct {
		name          string
		pageNumber    int
		pageSize      int
		sortBy        config.StudentCol
		sortOrder     config.SortOrder
		expectedCount int
		expectedNames []string
	}{
		{
			name:          "first page with 3 items",
			pageNumber:    1,
			pageSize:      3,
			sortBy:        config.StudentCol("Student_name"),
			sortOrder:     config.SortAsc,
			expectedCount: 3,
			expectedNames: []string{"Student01", "Student02", "Student03"},
		},
		{
			name:          "second page with 3 items",
			pageNumber:    2,
			pageSize:      3,
			sortBy:        config.StudentCol("Student_name"),
			sortOrder:     config.SortAsc,
			expectedCount: 3,
			expectedNames: []string{"Student04", "Student05", "Student06"},
		},
		{
			name:          "third page with 3 items",
			pageNumber:    3,
			pageSize:      3,
			sortBy:        config.StudentCol("Student_name"),
			sortOrder:     config.SortAsc,
			expectedCount: 3,
			expectedNames: []string{"Student07", "Student08", "Student09"},
		},
		{
			name:          "fourth page with 1 item (remainder)",
			pageNumber:    4,
			pageSize:      3,
			sortBy:        config.StudentCol("Student_name"),
			sortOrder:     config.SortAsc,
			expectedCount: 1,
			expectedNames: []string{"Student10"},
		},
		{
			name:          "first page with 5 items",
			pageNumber:    1,
			pageSize:      5,
			sortBy:        config.StudentCol("Student_name"),
			sortOrder:     config.SortAsc,
			expectedCount: 5,
			expectedNames: []string{"Student01", "Student02", "Student03", "Student04", "Student05"},
		},
		{
			name:          "sort by grade desc and paginate",
			pageNumber:    1,
			pageSize:      3,
			sortBy:        config.Grade,
			sortOrder:     config.SortDesc,
			expectedCount: 3,
			// Highest grades first
			expectedNames: []string{"Student10", "Student07", "Student09"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Get the specified page
			students, err := studentRepo.GetAll(tc.sortBy, tc.sortOrder, tc.pageNumber, tc.pageSize)
			assert.NoError(t, err)

			// Check page size
			assert.Len(t, students, tc.expectedCount)

			// Check student names in the page
			names := make([]string, len(students))
			for i, student := range students {
				names[i] = student.Student_name
			}
			assert.Equal(t, tc.expectedNames, names)
		})
	}
}

// Let's also refactor TestGetAll to use the same setup
func TestGetAll(t *testing.T) {
	studentRepo := setupTestData(t)

	cases := []struct {
		name          string
		expectedNames []string
		expectedCount int
		sortBy        config.StudentCol
		sortOrder     config.SortOrder
		checkOrder    bool
	}{
		{
			name:          "get all data without sorting",
			expectedNames: []string{"Student01", "Student02", "Student03", "Student04", "Student05", "Student06", "Student07", "Student08", "Student09", "Student10"},
			expectedCount: 10,
			sortBy:        "",
			sortOrder:     "",
			checkOrder:    false, // Don't check order for unsorted results
		},
		{
			name:          "get all data sorted by student name asc",
			expectedNames: []string{"Student01", "Student02", "Student03", "Student04", "Student05", "Student06", "Student07", "Student08", "Student09", "Student10"},
			expectedCount: 10,
			sortBy:        config.StudentCol("Student_name"),
			sortOrder:     config.SortAsc,
			checkOrder:    true,
		},
		{
			name:          "get all data sorted by subject desc",
			expectedNames: []string{"Student02", "Student08", "Student01", "Student05", "Student06", "Student09", "Student10", "Student03", "Student04", "Student07"},
			expectedCount: 10,
			sortBy:        config.Subject,
			sortOrder:     config.SortDesc,
			checkOrder:    true,
		},
		{
			name:          "get all data sorted by grade asc",
			expectedNames: []string{"Student06", "Student01", "Student04", "Student05", "Student02", "Student08", "Student03", "Student09", "Student07", "Student10"},
			expectedCount: 10,
			sortBy:        config.Grade,
			sortOrder:     config.SortAsc,
			checkOrder:    true,
		},
		{
			name:          "get all data sorted by grade desc",
			expectedNames: []string{"Student10", "Student07", "Student09", "Student03", "Student08", "Student02", "Student05", "Student04", "Student01", "Student06"},
			expectedCount: 10,
			sortBy:        config.Grade,
			sortOrder:     config.SortDesc,
			checkOrder:    true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			allStudents, err := studentRepo.GetAll(tt.sortBy, tt.sortOrder, 1, 20) // Get all students
			assert.NoError(t, err)
			assert.Len(t, allStudents, tt.expectedCount)

			// Verify expected order of names if order matters [AI]
			if len(tt.expectedNames) > 0 {
				names := make([]string, len(allStudents))
				for i, student := range allStudents {
					names[i] = student.Student_name
				}

				if tt.checkOrder {
					assert.Equal(t, tt.expectedNames, names)
				} else {
					// Check that all expected names are present, regardless of order
					assert.ElementsMatch(t, tt.expectedNames, names)
				}
			}
		})
	}
}

func TestFilterBySubject(t *testing.T) {
	studentRepo := setupTestData(t)

	testCases := []struct {
		name          string
		subject       config.Course
		pageNumber    int
		pageSize      int
		expectedCount int
		expectedNames []string
	}{
		{
			name:          "filter by Mathematics",
			subject:       config.Mathematics,
			pageNumber:    1,
			pageSize:      10,
			expectedCount: 1,
			expectedNames: []string{"Student01"},
		},
		{
			name:          "filter by EnglishLit",
			subject:       config.EnglishLit,
			pageNumber:    1,
			pageSize:      10,
			expectedCount: 1,
			expectedNames: []string{"Student09"},
		},
		{
			name:          "filter by Physics with pagination",
			subject:       config.Physics,
			pageNumber:    1,
			pageSize:      1,
			expectedCount: 1,
			expectedNames: []string{"Student02"},
		},
		{
			name:          "filter by non-existing subject",
			subject:       config.Course("Non-Existing"),
			pageNumber:    1,
			pageSize:      10,
			expectedCount: 0,
			expectedNames: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			students, err := studentRepo.FilterBySubject(tc.subject, tc.pageNumber, tc.pageSize)
			assert.NoError(t, err)

			// Check result count
			assert.Len(t, students, tc.expectedCount)

			// Extract names for comparison
			names := make([]string, len(students))
			for i, student := range students {
				names[i] = student.Student_name
			}

			// Check student names in results
			assert.ElementsMatch(t, tc.expectedNames, names)
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

// [AI]
func setupTestData(t *testing.T) repository.StudentRepository[model.StudentTest] {
	// Clear any previous test data
	testDB.Where("1=1").Delete(&model.StudentTest{})

	// Create test data - 10 students with ordered names
	studentRepo := repository.NewStudentRepository[model.StudentTest](testDB)
	testData := []*model.StudentTest{
		{Student_name: "Student01", Subject: string(config.Mathematics), Grade: 70},
		{Student_name: "Student02", Subject: string(config.Physics), Grade: 85},
		{Student_name: "Student03", Subject: string(config.Chemistry), Grade: 90},
		{Student_name: "Student04", Subject: string(config.Biology), Grade: 75},
		{Student_name: "Student05", Subject: string(config.History), Grade: 80},
		{Student_name: "Student06", Subject: string(config.Geography), Grade: 65},
		{Student_name: "Student07", Subject: string(config.Art), Grade: 95},
		{Student_name: "Student08", Subject: string(config.Music), Grade: 88},
		{Student_name: "Student09", Subject: string(config.EnglishLit), Grade: 92},
		{Student_name: "Student10", Subject: string(config.CompSci), Grade: 98},
	}
	err := studentRepo.CreateMany(testData)
	require.NoError(t, err, "Failed to create test data")

	return studentRepo
}
