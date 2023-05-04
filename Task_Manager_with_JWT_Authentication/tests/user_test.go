package tests

import (
	"testing"

	db2 "CRUD_REST/db"
)

func TestCreateUserInDB(t *testing.T) {
	// Создаем тестовую базу данных
	db, err := createTestDB()
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer db.Close()

	// Создаем пользователя в базе данных
	username := "testuser"
	password := "testpassword"
	err = db2.CreateUserInDB(username, password)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Проверяем, что пользователь действительно добавлен в базу данных
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = $1", username).Scan(&count)
	if err != nil {
		t.Fatalf("failed to count users: %v", err)
	}
	if count != 1 {
		t.Fatalf("unexpected number of users: %d", count)
	}
}
