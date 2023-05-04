package tests

import "database/sql"

var godotenvLoad = func(filename string) error {
	return nil
}

var viperGetString = func(key string) string {
	return "test_value"
}

var sqlOpen = func(driverName, dataSourceName string) (*sql.DB, error) {
	return &sql.DB{}, nil
}
