package database

import (
	"database/sql"

	"github.com/nedson202/user-service/config"
)

// CreateUserTable with psql
func CreateUserTable(db *sql.DB) {
	var query = `
		CREATE TABLE IF NOT EXISTS users (
			id serial PRIMARY KEY,
			username text UNIQUE NOT NULL,
			email text UNIQUE NOT NULL,
			password text NOT NULL,
			role text NOT NULL DEFAULT 'User',
			created_at timestamp with time zone DEFAULT current_timestamp NOT NULL,
			updated_at timestamp with time zone DEFAULT current_timestamp NOT NULL,
			deleted boolean,
			deleted_at timestamp with time zone
		)
	`
	_, err := db.Query(query)
	config.LogFatal(err)

	return
}

// DropUserTable with psql
func DropUserTable(db *sql.DB) {
	var query = `
		DROP TABLE users
	`
	_, err := db.Query(query)
	config.LogFatal(err)

	return
}
