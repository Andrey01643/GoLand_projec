package handlers

import (
	"html/template"
	"log"
	"net/http"

	"httpServer/db/postgres"
	"httpServer/models"
)

//Обработчик главной страницы с выбором варианта теста
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем соединение с БД
	db, err := postgres.ConnectDb()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	user, cookie, err := checkUserCookie(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if cookie == nil {
		http.Error(w, "User cookie not found", http.StatusNotFound)
		return
	}

	// Получение списка вариантов теста из БД
	rows, err := db.Query("SELECT id, name FROM variants")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var variants []models.Variant
	for rows.Next() {
		var v models.Variant
		if err := rows.Scan(&v.ID, &v.Name); err != nil {
			log.Fatal(err)
		}
		variants = append(variants, v)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// Отрисовка HTML шаблона
	t, err := template.ParseFiles("static/variantTests.html")
	if err != nil {
		log.Fatal(err)
	}
	data := models.TemplateData{Variants: variants}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal(err)
	}

}
