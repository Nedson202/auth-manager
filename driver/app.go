package driver

import (
	"database/sql"
	"os"

	"github.com/lib/pq"
	"github.com/nedson202/user-service/config"
)

var db *sql.DB

// SetupDatabaseConnection with sql and postgres
func SetupDatabaseConnection() *sql.DB {
	dbURL, err := pq.ParseURL(os.Getenv("DATABASE_URL"))
	config.LogFatal(err)

	db, err = sql.Open("postgres", dbURL)
	config.LogFatal(err)

	err = db.Ping()
	config.LogFatal(err)

	return db
}

// DB connection instance
var DB = SetupDatabaseConnection()
