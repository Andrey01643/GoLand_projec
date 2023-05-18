package handlers

import "net/http"

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Очищаем cookie с идентификатором пользователя
	cookie := &http.Cookie{
		Name:   "user_id",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)

	// Редирект на главную страницу
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
