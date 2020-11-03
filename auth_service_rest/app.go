package auth_service_rest

import (
	"github.com/Masterminds/squirrel"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

var queryBuilder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

// New _
func New(router *mux.Router, db *sqlx.DB, jwtSecret string, cachePool *redis.Pool) (app App, err error) {
	app.db = db
	app.jwtSecret = []byte(jwtSecret)
	app.cachePool = cachePool

	app.newRouter(router)
	return
}
