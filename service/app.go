package service

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

// New _
func New(router *mux.Router, db *sqlx.DB) (app App, err error) {
	app.db = db
	app.migrateDatabaseTables(db)
	app.newRouter(router)

	return
}
