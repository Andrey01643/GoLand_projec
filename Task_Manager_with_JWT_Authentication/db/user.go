package db

func CreateUserInDB(username string, password string) error {
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	// Вставить из тасков
	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, password)
	if err != nil {
		return err
	}

	return nil
}
