package service

import (
	"github.com/jmoiron/sqlx"
)

func (app App) migrateDatabaseTables(db *sqlx.DB) {
	app.createUserTable(db)
}
