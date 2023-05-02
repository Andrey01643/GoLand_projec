package db

import (
	"database/sql"
	"fmt"
	"time"

	"CRUD_REST/models"
)

func GetTasksFromDB(userID int) ([]models.Task, error) {
	dataBase, errors := OpenDB()
	defer dataBase.Close()

	rows, err := dataBase.Query("SELECT id, name, completed, completed_at, created_at, updated_at, user_id FROM tasks WHERE user_id=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []models.Task{}
	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Name, &task.Completed, &task.CompletedAt, &task.CreatedAt, &task.UpdatedAt, &task.UserID)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, errors
}

func GetTaskByID(id int) (*models.Task, error) {
	dataBase, errors := OpenDB()
	defer dataBase.Close()

	var task models.Task
	err := dataBase.QueryRow("SELECT * FROM tasks WHERE id=$1", id).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt, &task.UserID, &task.Name, &task.Completed, &task.CompletedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Task not found")
		}
		return nil, err
	}
	return &task, errors
}

func CreateTaskInDB(userID int, name string) (int, error) {
	dataBase, errors := OpenDB()
	defer dataBase.Close()
	stmt, err := dataBase.Prepare("INSERT INTO tasks (user_id, name) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	var taskID int
	err = stmt.QueryRow(userID, name).Scan(&taskID)
	if err != nil {
		return 0, err
	}
	return taskID, errors
}

func UpdateTaskInDB(id int, task models.Task) error {
	dataBase, errors := OpenDB()
	defer dataBase.Close()

	stmt, err := dataBase.Prepare("UPDATE tasks SET name=$1, completed=$2, completed_at=$3, updated_at=$4 WHERE id=$5")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	completedAt := sql.NullTime{}
	if task.Completed {
		completedAt = sql.NullTime{Time: time.Now(), Valid: true}
	}

	_, err = stmt.Exec(task.Name, task.Completed, completedAt, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %v", err)
	}

	return errors
}

func DeleteTaskFromDB(taskID int) error {
	dataBase, errors := OpenDB()
	defer dataBase.Close()

	// Удаление задачи из базы данных
	stmt, err := dataBase.Prepare("DELETE FROM tasks WHERE id=$1")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(taskID)
	if err != nil {
		return err
	}

	return errors
}
