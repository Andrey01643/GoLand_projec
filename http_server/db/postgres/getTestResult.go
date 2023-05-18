package postgres

import "database/sql"

func GetTestResults(db *sql.DB, userID string, variantID int) (float64, error) {

	// Выполняем запрос на выборку количества правильных ответов для данного пользователя и выбранного варианта
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM results JOIN tasks ON results.task_id=tasks.id WHERE results.user_id=$1 AND results.is_correct=true AND tasks.variant_id=$2", userID, variantID).Scan(&count)
	if err != nil {
		return 0, err
	}

	// Выполняем запрос на выборку общего количества заданий
	var total int
	err = db.QueryRow("SELECT COUNT(*) FROM tasks WHERE tasks.variant_id = $1 AND tasks.variant_id IN (SELECT variant_id FROM answers)", variantID).Scan(&total)
	if err != nil {
		return 0, err
	}

	// Вычисляем процент правильных ответов
	result := float64(count) / float64(total) * 100

	return result, nil
}
