package config

type StudentCol string
type SortOrder string

const (
	DBEnvVar = "DB_DSN"

	Id      StudentCol = "Student_id"
	Name    StudentCol = "Student_name"
	Subject StudentCol = "Subject"
	Grade   StudentCol = "Grade"

	SortAsc  SortOrder = "asc"
	SortDesc SortOrder = "desc"
)
