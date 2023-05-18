package handlers

import (
	"httpServer/db/postgres"
	"httpServer/models"
	"net/http"
)

func checkUserCookie(w http.ResponseWriter, r *http.Request) (*models.User, *http.Cookie, error) {
	// Устанавливаем соединение с БД
	db, err := postgres.ConnectDb()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, nil, err
	}
	defer db.Close()

	// Проверяем наличие cookie с идентификатором пользователя
	cookie, err := r.Cookie("user_id")
	if err != nil {
		// Если cookie отсутствует, редиректим на страницу авторизации
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return nil, nil, err
	}

	// Выполняем запрос на выборку пользователя по идентификатору
	var user models.User
	err = db.QueryRow("SELECT * FROM auth WHERE id=$1", cookie.Value).Scan(&user.ID, &user.Login, &user.Password, &user.IsAuthorized, &user.LoginTime, &user.LogoutTime)
	if err != nil {
		// Если пользователь не найден в БД, редиректим на страницу авторизации
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return nil, nil, err
	}

	return &user, cookie, nil
}
