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

	// Проверяем наличие cookie с идентификатором пользователя
	cookie, err := r.Cookie("user_id")
	if err != nil {
		// Если cookie отсутствует, редиректим на страницу авторизации
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	// Выполняем запрос на выборку пользователя по идентификатору
	var user models.User
	err = db.QueryRow("SELECT * FROM auth WHERE id=$1", cookie.Value).Scan(&user.ID, &user.Login, &user.Password, &user.IsAuthorized, &user.LoginTime, &user.LogoutTime)
	if err != nil {
		// Если пользователь не найден в БД, редиректим на страницу авторизации
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
