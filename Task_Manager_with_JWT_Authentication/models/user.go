package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) Authenticate(username, password string) bool {
	return u.Username == username && u.Password == password
}
