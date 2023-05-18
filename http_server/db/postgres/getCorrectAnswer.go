package postgres

import "database/sql"

func GetCorrectAnswer(db *sql.DB, taskID int) (string, error) {
	var correctAnswer string
	err := db.QueryRow("SELECT correct_answer FROM tasks WHERE id = $1", taskID).Scan(&correctAnswer)
	if err != nil {
		return "", err
	}
	return correctAnswer, nil
}
