package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"httpServer/db/postgres"
	"httpServer/models"
)

// Обработчик запроса на авторизацию
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем соединение с БД
	db, err := postgres.ConnectDb()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Получаем логин и пароль из POST запроса
	login := r.PostFormValue("login")
	password := r.PostFormValue("password")

	fmt.Println(login, password)
	// Проверяем, что логин и пароль не являются пустыми строками
	if login == "" || password == "" {
		http.Error(w, "Login and password can't be empty", http.StatusBadRequest)
		return
	}

	// Выполняем запрос на выборку пользователя
	var user models.User
	err = db.QueryRow("SELECT * FROM auth WHERE login=$1 AND password=$2", login, password).Scan(&user.ID, &user.Login, &user.Password, &user.IsAuthorized, &user.LoginTime, &user.LogoutTime)
	if err != nil {
		http.Error(w, "Invalid login or password", http.StatusUnauthorized)
		return
	}

	// Обновляем информацию о пользователе
	user.LoginTime = time.Now()
	user.IsAuthorized = true

	// Подготавливаем запрос на обновление информации о пользователе
	stmt, err := db.Prepare("UPDATE auth SET login_time=$1, is_authorized=$2 WHERE id=$3")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	// Выполняем запрос на обновление информации о пользователе
	_, err = stmt.Exec(user.LoginTime, user.IsAuthorized, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Устанавливаем cookie с идентификатором пользователя
	cookie := &http.Cookie{
		Name:  "user_id",
		Value: strconv.Itoa(user.ID),
	}
	http.SetCookie(w, cookie)

	// Редиректим на страницу с выбором варианта теста
	http.Redirect(w, r, "/test", http.StatusSeeOther)
}
