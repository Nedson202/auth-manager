package service

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func (app App) createUserTable(db *sqlx.DB) {
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
	if err != nil {
		log.Fatal(err)
	}

	return
}

func (app App) dropUserTable(db *sqlx.DB) {
	var query = `
		DROP TABLE users
	`
	_, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	return
}
