package tests

import (
	db2 "CRUD_REST/db"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenDB(t *testing.T) {
	mockDB := &sql.DB{}
	mockErr := errors.New("some error")

	godotenvLoad = func(filename string) error {
		return nil
	}
	viperGetString = func(key string) string {
		switch key {
		case "db_user":
			return "postgres"
		case "db_password":
			return "200100"
		case "db_host":
			return "localhost"
		case "db_port":
			return "5432"
		case "db_name":
			return "crudRestAuth"
		default:
			return ""
		}
	}

	sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
		if driverName == "postgres" && dataSourceName == "postgres://postgres:200100@localhost:5432/crudRestAuth?sslmode=disable" {
			return mockDB, nil
		}
		return nil, mockErr
	}

	// Test
	db, err := db2.OpenDB()
	assert.Equal(t, mockDB, db)
	assert.Nil(t, err)

	// Test error case
	sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
		return nil, mockErr
	}
	db, err = db2.OpenDB()
	assert.Nil(t, db)
	assert.Equal(t, mockErr, err)
}
