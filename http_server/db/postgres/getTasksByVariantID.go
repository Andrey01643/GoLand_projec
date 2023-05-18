package postgres

import (
	"database/sql"
	"httpServer/models"
)

// Получение всех заданий из БД для выбранного варианта

func GetTasksByVariantID(db *sql.DB, variantID int, t models.Task) (models.Task, error) {

	rows, err := db.Query("SELECT id, variant_id, task, correct_answer, answer_1, answer_2, answer_3, answer_4 FROM tasks WHERE variant_id = $1 ORDER BY id", variantID)
	if err != nil {
		return t, err
	}
	defer rows.Close()

	var tasks []models.Task
	//var t models.Task
	for rows.Next() {
		err = rows.Scan(&t.ID, &t.VariantID, &t.Task, &t.CorrectAnswer, &t.Answer1, &t.Answer2, &t.Answer3, &t.Answer4)
		if err != nil {
			return t, err
		}
		tasks = append(tasks, t)
	}
	err = rows.Err()
	if err != nil {
		return t, err
	}

	return t, nil
}

func GetTasks(db *sql.DB, variantID int, t models.Task) (models.Task, error) {
	rows, err := db.Query(`
		SELECT id, variant_id, task, correct_answer, answer_1, answer_2, answer_3, answer_4
		FROM tasks
		WHERE variant_id = $1 AND id NOT IN (
			SELECT task_id
			FROM answers
		)
		ORDER BY id`, variantID)
	if err != nil {
		return t, err
	}
	defer rows.Close()

	var tasks []models.Task
	//var t models.Task

	for rows.Next() {
		err = rows.Scan(&t.ID, &t.VariantID, &t.Task, &t.CorrectAnswer, &t.Answer1, &t.Answer2, &t.Answer3, &t.Answer4)
		if err != nil {
			return t, err
		}
		tasks = append(tasks, t)
	}

	err = rows.Err()
	if err != nil {
		return t, err
	}

	return t, nil
}
