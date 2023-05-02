package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

// OpenDB открывает соединение с базой данных PostgreSQL и возвращает объект sql.DB
func OpenDB() (*sql.DB, error) {
	err := godotenv.Load("config/config.yml")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Значения из файла конфигурации
	dbUser := viper.GetString("db_user")
	dbPassword := viper.GetString("db_password")
	dbHost := viper.GetString("db_host")
	dbPort := viper.GetString("db_port")
	dbName := viper.GetString("db_name")

	// Подключение к базе данных
	dbURI := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("postgres", dbURI)
	if err != nil {
		log.Fatal("Error connecting to db: ", err)
	}

	return db, err
}
