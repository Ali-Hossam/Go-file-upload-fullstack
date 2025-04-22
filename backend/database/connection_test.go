package database_test

import (
	"errors"
	"file-uploader/database"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

const (
	ErrEnvVarNotFound  = "environment variable doesn't exist"
	ErrDotEnvNotLoaded = "error loading .env file"

	DBEnvVar = "DB_DSN"
)

func TestSetupDB(t *testing.T) {
	t.Run("testing valid database dsn, should return no error", func(t *testing.T) {
		err := godotenv.Load("../.env")
		if err != nil {
			assert.Error(t, errors.New(ErrDotEnvNotLoaded))
		}
		dsn, exist := os.LookupEnv(DBEnvVar)

		if !exist {
			assert.Error(t, errors.New(ErrEnvVarNotFound))
		}

		db, err := database.SetupDB(dsn)
		assert.NoError(t, err)
		assert.NotNil(t, db)
	})

	t.Run("testing invalid database dsn, should return error", func(t *testing.T) {
		db, err := database.SetupDB("this is invalid dsn")
		assert.ErrorContains(t, err, database.ErrFailedDBConnection)
		assert.Nil(t, db)
	})
}
