package main

import "database/sql"

func searchUsers(db *sql.DB, q string) (*sql.Rows, error) {
	return db.Query("SELECT * FROM users WHERE name = '" + q + "'")
}
