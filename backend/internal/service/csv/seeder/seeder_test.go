package Seeder_test

import (
	"encoding/csv"
	Seeder "file-uploader/internal/service/csv/seeder"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSeeder(t *testing.T) {
	t.Run("create one student record", func(t *testing.T) {
		student := Seeder.CreateStudentRecord()
		assert.NotEmpty(t, student.Student_name)
		assert.NotEmpty(t, student.Student_id)
		assert.NotEmpty(t, student.Grade)
		assert.NotEmpty(t, student.Subject)
	})

	t.Run("create students csv file with 10 records", func(t *testing.T) {
		const (
			testFilesDir  = "/tmp/testDir/"
			recordsLength = 10
		)

		filepath, err := Seeder.SeedStudentsCSV("test.csv", testFilesDir, recordsLength)
		require.NoError(t, err)

		// Check length of records in the saved file
		savedFile, err := os.Open(filepath)
		require.NoError(t, err)
		defer savedFile.Close()

		reader := csv.NewReader(savedFile)
		records, err := reader.ReadAll()
		require.NoError(t, err)

		assert.Equal(t, recordsLength+1, len(records)) // header row

		// Clear
		err = Seeder.RemoveSeededCSVs(testFilesDir)
		require.NoError(t, err)
	})
}
