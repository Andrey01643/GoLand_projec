package models

import "time"

// Структура для хранения информации о пользователе
type User struct {
	ID           int       `json:"id"`
	Login        string    `json:"login"`
	Password     string    `json:"password"`
	IsAuthorized bool      `json:"is_authorized"`
	LoginTime    time.Time `json:"login_time"`
	LogoutTime   time.Time `json:"logout_time"`
}
