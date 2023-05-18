package postgres

import (
	"database/sql"
	"time"
)

func InsertAnswer(db *sql.DB, authID string, variantID, taskID int, answer string) error {
	_, err := db.Exec("INSERT INTO answers (auth_id, variant_id, task_id, answer, answer_time) VALUES ($1, $2, $3, $4, $5)",
		authID, variantID, taskID, answer, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func InsertResult(db *sql.DB, variantID, taskID int, userID string, isCorrect bool) error {
	_, err := db.Exec("INSERT INTO results (variant_id, task_id, user_id, is_correct, answer_time) VALUES ($1, $2, $3, $4, $5)",
		variantID, taskID, userID, isCorrect, time.Now())
	if err != nil {
		return err
	}
	return nil
}
