package postgres

import (
	"database/sql"
	"httpServer/models"
)

// Получение информации о варианте из БД
func GetVariantByID(db *sql.DB, id int) (models.Variant, error) {

	var v models.Variant
	err := db.QueryRow("SELECT id, name FROM variants WHERE id = $1", id).Scan(&v.ID, &v.Name)
	if err != nil {
		return models.Variant{}, err
	}

	return v, nil
}
