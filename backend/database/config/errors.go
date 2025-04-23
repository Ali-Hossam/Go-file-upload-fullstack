package config

import "errors"

var (
	ErrFailedDBConnection = errors.New("failed to connect to postgress db")
	ErrFailedMigration    = errors.New("failed migration")
	ErrEnvVarNotFound     = errors.New("environment variable doesn't exist")
	ErrDotEnvNotLoaded    = errors.New("error loading .env file")
	ErrFieldNotFound      = errors.New("field not found")
	ErrMissingStudentData = errors.New("student data are missing, required name, subject and grade")
	ErrStudentNotExist    = errors.New("student does not exist")
)
