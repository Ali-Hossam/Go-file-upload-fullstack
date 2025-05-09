package config

type StudentCol string
type SortOrder string
type Course string

const (
	DBEnvVar   = "DB_DSN_LOCAL"
	PortEnvVar = "SERVER_PORT"

	Id      StudentCol = "Student_id"
	Name    StudentCol = "Student_name"
	Subject StudentCol = "Subject"
	Grade   StudentCol = "Grade"

	SortAsc  SortOrder = "asc"
	SortDesc SortOrder = "desc"

	Mathematics Course = "Mathematics"
	Physics     Course = "Physics"
	Chemistry   Course = "Chemistry"
	Biology     Course = "Biology"
	History     Course = "History"
	EnglishLit  Course = "English Literature"
	CompSci     Course = "Computer Science"
	Art         Course = "Art"
	Music       Course = "Music"
	Geography   Course = "Geography"

	StudentsTableHeader = "student_id,student_name,subject,grade"
)
