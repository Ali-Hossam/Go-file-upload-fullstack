package Seeder

import (
	"encoding/csv"
	"file-uploader/config"
	"file-uploader/database/model"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func SeedStudentsCSV(filename, saveLocation string, length int) (string, error) {
	path := filepath.Join(saveLocation, filename)

	// Check if saving directory exists before saving
	_, err := os.Stat(saveLocation)
	if os.IsNotExist(err) {
		err := os.Mkdir(saveLocation, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create directory %w", err)
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"Student_id", "Student_name", "Subject", "Grade"}
	if err := writer.Write(header); err != nil {
		return "", fmt.Errorf("failed to write header: %w", err)
	}

	for range length {
		student := CreateStudentRecord()
		record := []string{
			student.Student_id.String(),
			student.Student_name,
			student.Subject,
			fmt.Sprintf("%d", student.Grade), // Convert uint to string
		}
		if err := writer.Write(record); err != nil {
			return "", fmt.Errorf("failed to write record: %w", err)
		}
	}

	return path, nil
}

func CreateStudentRecord() model.StudentTest {
	// Sample data for names and subjects
	names := []string{"Omar", "Ali", "Saad", "Mohamed", "Ahmed"}
	subjects := []config.Course{config.Art, config.Geography, config.Biology, config.History}

	// Pick random name and subject
	name := names[rand.Intn(len(names))]
	subject := subjects[rand.Intn(len(subjects))]

	// Generate random grade between 0 and 100 inclusive
	grade := uint(rand.Intn(101))

	// Generate new UUID
	id := uuid.New()

	return model.StudentTest{
		Student_id:   id,
		Student_name: name,
		Subject:      string(subject),
		Grade:        grade,
	}
}

func RemoveSeededCSVs(saveLocation string) error {
	return os.RemoveAll(saveLocation)
}
