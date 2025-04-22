package database_test

import (
	"file-uploader/database"
	"file-uploader/database/config"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestSetupDB(t *testing.T) {
	t.Run("testing valid database dsn, should return no error", func(t *testing.T) {
		err := godotenv.Load("../.env")
		if err != nil {
			assert.Error(t, config.ErrDotEnvNotLoaded)
		}
		dsn, exist := os.LookupEnv(config.DBEnvVar)

		if !exist {
			assert.Error(t, config.ErrEnvVarNotFound)
		}

		db, err := database.SetupDB(dsn)
		assert.NoError(t, err)
		assert.NotNil(t, db)
	})

	t.Run("testing invalid database dsn, should return error", func(t *testing.T) {
		db, err := database.SetupDB("this is invalid dsn")
		assert.ErrorContains(t, err, config.ErrFailedDBConnection.Error())
		assert.Nil(t, db)
	})
}
