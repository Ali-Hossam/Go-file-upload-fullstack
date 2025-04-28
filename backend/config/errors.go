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

const (
	ErrFormParseFailureHttp    = "Failed to parse multipart form"
	ErrNoFilesProvidedHttp     = "No files were provided for upload"
	ErrDBConfigNotFoundHttp    = "Database configuration not found"
	ErrDBConnectionFailureHttp = "Failed to connect to database"
	ErrFileOpenFailureHttp     = "Failed to open uploaded file"
	ErrProcessingFailureHttp   = "Failed to process CSV data"
	ErrInvalidFileTypeHttp     = "Invalid File type"
	ErrInvalidCSVCols          = "Invalid CSV columns"
	ErrInvalidFilterHttp       = "Invalid filter"
	ErrMissingPathParamHttp    = "Missing path parameter"
	ErrMissingSearchParamHttp  = "Missing search parameter"
	ErrInvalidSearchParamHttp  = "Invalid search parameter"
)
