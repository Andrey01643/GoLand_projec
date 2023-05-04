package tests

import (
	"database/sql"
	"fmt"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

// GetDSN возвращает строку подключения к базе данных.
func (c *DBConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port, c.User, c.Password, c.Name)
}

// OpenDB открывает соединение с базой данных, указанной в конфигурации.
func (c *DBConfig) OpenDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", c.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}
	return db, nil
}

// Тестовая конфигурация базы данных.
var testDBConfig = DBConfig{
	Host:     "localhost",
	Port:     "5432",
	User:     "postgres",
	Password: "200100",
	Name:     "crudRestAuth",
}

func createTasksTable(db *sql.DB) error {
	// Проверяем, существует ли таблица tasks
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT FROM pg_tables WHERE schemaname = 'public' AND tablename = 'tasks')").Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if tasks table exists: %v", err)
	}
	if exists {
		return nil // таблица уже существует, ничего не делаем
	}

	// Создаем таблицу tasks
	query := `CREATE TABLE tasks (
                id SERIAL PRIMARY KEY,
                created_at TIMESTAMP NOT NULL,
                updated_at TIMESTAMP NOT NULL,
                user_id INTEGER NOT NULL,
                name TEXT NOT NULL,
                completed BOOLEAN NOT NULL,
                completed_at TIMESTAMP
            )`
	_, err = db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create tasks table: %v", err)
	}
	return nil
}

// Создаем тестовую базу данных.
func createTestDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", testDBConfig.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	// Удаляем тестовую базу данных, если она уже существует
	_, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", testDBConfig.Name))
	if err != nil {
		return nil, fmt.Errorf("failed to drop database: %v", err)
	}

	// Создаем тестовую базу данных
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", testDBConfig.Name))
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %v", err)
	}

	// Связываем тестовую базу данных с конфигурацией
	testDBConfig.Name = "postgres"
	db, err = testDBConfig.OpenDB()
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	// Создаем таблицу tasks, если ее еще нет
	err = createTasksTable(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create tasks table: %v", err)
	}

	return db, nil
}
