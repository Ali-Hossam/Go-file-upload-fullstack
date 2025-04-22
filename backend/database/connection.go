package database

import (
	"errors"

	"file-uploader/database/config"
	"file-uploader/database/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn, PreferSimpleProtocol: true}), &gorm.Config{})

	if err != nil {
		return nil, errors.New(config.ErrFailedDBConnection.Error() + ": " + err.Error())
	}

	err = db.AutoMigrate(&model.Student{}, &model.StudentTest{})
	if err != nil {
		return nil, errors.New(config.ErrFailedMigration.Error() + " : " + err.Error())
	}

	return db, nil
}
