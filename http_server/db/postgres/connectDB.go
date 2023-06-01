package postgres

import (
	"database/sql"
	"fmt"
)

func ConnectDb() (*sql.DB, error) {
	dbUser := "postgres"
	dbPassword := "200100"
	dbHost := "localhost"
	dbPort := "5432"
	dbName := "postgres"
	sslMode := "disable"

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbUser, dbPassword, dbHost, dbPort, dbName, sslMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}
