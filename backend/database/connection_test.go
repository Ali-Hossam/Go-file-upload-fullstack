package database_test

import (
	"file-uploader/config"
	"file-uploader/database"
	"file-uploader/database/model"
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

		db, repo, err := database.SetupDB[model.StudentTest](dsn, true)
		assert.NoError(t, err)
		assert.NotNil(t, db)
		assert.NotNil(t, repo)
	})

	t.Run("testing invalid database dsn, should return error", func(t *testing.T) {
		db, repo, err := database.SetupDB[model.StudentTest]("this is invalid dsn", true)
		assert.ErrorContains(t, err, config.ErrFailedDBConnection.Error())
		assert.Nil(t, db)
		assert.Nil(t, repo)
	})
}
