package database

import (
	"errors"

	"file-uploader/config"
	"file-uploader/database/model"
	"file-uploader/database/repository"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDB[T any](dsn string, isTest bool) (*gorm.DB, repository.StudentRepository[T], error) {
	db, err := gorm.Open(
		postgres.New(postgres.Config{DSN: dsn, PreferSimpleProtocol: true}),
		&gorm.Config{},
	)

	if err != nil {
		return nil, nil, errors.New(config.ErrFailedDBConnection.Error() + ": " + err.Error())
	}

	err = db.AutoMigrate(&model.Student{}, &model.StudentTest{})
	if err != nil {
		return nil, nil, errors.New(config.ErrFailedMigration.Error() + " : " + err.Error())
	}

	StudentRepo := repository.NewStudentRepository[T](db)
	return db, StudentRepo, nil
}
