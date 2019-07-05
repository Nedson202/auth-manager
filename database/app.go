package database

import (
	"database/sql"
)

// MigrateDatabaseTables to add tables
func MigrateDatabaseTables(db *sql.DB) {
	CreateUserTable(db)
}
