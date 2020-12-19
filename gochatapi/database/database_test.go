package database

import (
	"testing"
	"api/utils"

	"github.com/stretchr/testify/assert"
)

func TestDatabase(t *testing.T) {
	db := NewDatabase(&Config{
		DSN: utils.GetEnv("DSN",
			"user=postgres password=1234 dbname=postgres port=5432 host=localhost sslmode=disable"),
	})

	err := db.Open()
	assert.NoError(t, err)
	if err == nil {
		defer db.Close()
	}

	assert.NoError(t, db.Migrate("./migrations"))
}
