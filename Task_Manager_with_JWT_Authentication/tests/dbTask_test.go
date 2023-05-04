package tests

import (
	db2 "CRUD_REST/db"
	"CRUD_REST/models"
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestGetTasksFromDB(t *testing.T) {
	// Создание тестовой базы данных
	db, err := createTestDB()
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer db.Close()

	// Заполнение базы данных тестовыми данными
	_, err = db.Exec("INSERT INTO tasks (name, completed, completed_at, created_at, updated_at, user_id) VALUES ($1, $2, $3, $4, $5, $6)", "test task 1", true, time.Now(), time.Now(), time.Now(), 1)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// Вызов тестируемой функции
	tasks, err := db2.GetTasksFromDB(1)
	if err != nil {
		t.Fatalf("failed to get tasks from database: %v", err)
	}

	// Проверка полученных результатов
	if len(tasks) != 1 {
		t.Errorf("expected 1 task, but got %d", len(tasks))
	}

	expectedTask := models.Task{
		Name:      "test task 1",
		Completed: true,
		UserID:    1,
	}
	if !reflect.DeepEqual(tasks[0], expectedTask) {
		t.Errorf("expected task %+v, but got %+v", expectedTask, tasks[0])
	}
}

func TestGetTaskByID(t *testing.T) {
	// Создаем тестовую базу данных и таблицу
	db, err := createTestDB()
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer db.Close()

	// Вставляем тестовые данные в таблицу
	query := `INSERT INTO tasks (created_at, updated_at, user_id, name, completed, completed_at)
			  VALUES ($1, $2, $3, $4, $5, $6)
			  RETURNING id`
	var taskID int
	err = db.QueryRow(query, time.Now(), time.Now(), 1, "Test Task", false, nil).Scan(&taskID)
	if err != nil {
		t.Fatalf("failed to insert test task: %v", err)
	}

	// Тестирование функции при запросе существующей задачи
	task, err := db2.GetTaskByID(taskID)
	if err != nil {
		t.Fatalf("failed to get task by ID: %v", err)
	}
	if task.ID != taskID {
		t.Errorf("unexpected task ID: got %d, want %d", task.ID, taskID)
	}
	if task.Name != "Test Task" {
		t.Errorf("unexpected task name: got %s, want %s", task.Name, "Test Task")
	}

	// Тестирование функции при запросе несуществующей задачи
	task, err = db2.GetTaskByID(taskID + 1)
	if err == nil {
		t.Errorf("expected error, but got nil")
	}
	if task != nil {
		t.Errorf("unexpected task: got %+v, want nil", task)
	}
	if err.Error() != "Task not found" {
		t.Errorf("unexpected error message: got %s, want %s", err.Error(), "Task not found")
	}
}

func TestCreateTaskInDB(t *testing.T) {
	// Создаем тестовую базу данных
	testDB, err := createTestDB()
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer testDB.Close()

	// Создаем таблицу tasks в тестовой базе данных
	_, err = testDB.Exec(`CREATE TABLE IF NOT EXISTS tasks (
        id SERIAL PRIMARY KEY,
        user_id INTEGER NOT NULL,
        name TEXT NOT NULL
    )`)
	if err != nil {
		t.Fatalf("failed to create tasks table: %v", err)
	}

	// Подготавливаем данные для теста
	userID := 1
	taskName := "test task"

	// Вызываем функцию, которую хотим протестировать
	taskID, err := db2.CreateTaskInDB(userID, taskName)
	if err != nil {
		t.Fatalf("failed to create task in database: %v", err)
	}

	// Проверяем, что задача была успешно создана
	var count int
	err = testDB.QueryRow("SELECT COUNT(*) FROM tasks WHERE id = $1 AND user_id = $2 AND name = $3", taskID, userID, taskName).Scan(&count)
	if err != nil {
		t.Fatalf("failed to query task: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 row in tasks table, got %d", count)
	}
}

func InsertTaskIntoDB(db *sql.DB, task models.Task) error {
	// Создаем запрос на вставку новой записи в таблицу tasks
	_, err := db.Exec("INSERT INTO tasks (name, completed) VALUES ($1, $2)", task.Name, task.Completed)
	if err != nil {
		return fmt.Errorf("failed to insert task into database: %v", err)
	}
	return nil
}

func TestUpdateTaskInDB(t *testing.T) {
	// Создаем тестовую базу данных
	db, err := createTestDB()
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer db.Close()

	// Вставляем тестовые данные в таблицу tasks
	task := models.Task{Name: "test task", Completed: false}
	err = InsertTaskIntoDB(db, task)
	if err != nil {
		t.Fatalf("failed to insert task into database: %v", err)
	}

	// Обновляем данные в таблице tasks
	task.Name = "updated task"
	task.Completed = true
	err = db2.UpdateTaskInDB(1, task)
	if err != nil {
		t.Fatalf("failed to update task in database: %v", err)
	}

	// Проверяем, что данные были обновлены
	rows, err := db.Query("SELECT * FROM tasks WHERE id=1")
	if err != nil {
		t.Fatalf("failed to query tasks table: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		t.Fatalf("no rows returned from tasks table")
	}

	var updatedTask models.Task
	err = rows.Scan(&updatedTask.ID, &updatedTask.Name, &updatedTask.Completed, &updatedTask.CompletedAt, &updatedTask.CreatedAt, &updatedTask.UpdatedAt)
	if err != nil {
		t.Fatalf("failed to scan row from tasks table: %v", err)
	}

	if updatedTask.Name != "updated task" {
		t.Errorf("expected task name to be %q, but got %q", "updated task", updatedTask.Name)
	}

	if !updatedTask.Completed {
		t.Errorf("expected task to be completed")
	}

	if updatedTask.UpdatedAt.IsZero() {
		t.Errorf("expected updated_at field to be set")
	}
}

func TestDeleteTaskFromDB(t *testing.T) {
	// Создаем тестовую базу данных
	db, err := createTestDB()
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}
	defer db.Close()

	// Вставляем тестовые данные в таблицу tasks
	query := `INSERT INTO tasks (created_at, updated_at, user_id, name, completed, completed_at)
			  VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = db.Exec(query, time.Now(), time.Now(), 1, "test task", false, nil)
	if err != nil {
		t.Fatalf("failed to insert test data: %v", err)
	}

	// Удаляем задачу из таблицы tasks
	err = db2.DeleteTaskFromDB(1)
	if err != nil {
		t.Fatalf("failed to delete task: %v", err)
	}

	// Проверяем, что задача действительно удалена
	var exists bool
	err = db.QueryRow("SELECT EXISTS (SELECT 1 FROM tasks WHERE id=$1)", 1).Scan(&exists)
	if err != nil {
		t.Fatalf("failed to check if task exists: %v", err)
	}
	if exists {
		t.Fatalf("task was not deleted")
	}
}
