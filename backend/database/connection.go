package database

import (
	"errors"

	"file-uploader/database/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	ErrFailedDBConnection = "failed to connect to postgress db"
	ErrFailedMigration    = "failed migration"
)

func SetupDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn, PreferSimpleProtocol: true}), &gorm.Config{})

	if err != nil {
		return nil, errors.New(ErrFailedDBConnection + ": " + err.Error())
	}

	err = db.AutoMigrate(&model.Student{})
	if err != nil {
		return nil, errors.New(ErrFailedMigration + " : " + err.Error())
	}

	return db, nil
}
