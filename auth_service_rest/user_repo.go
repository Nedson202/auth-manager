package auth_service_rest

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
)

const usersTable = "users"

func (app App) getUserByID(id string) (user User, err error) {
	query, args, err := queryBuilder.
		Select("id, username, email, role, created_at, updated_at").From(usersTable).
		Where(squirrel.Eq{"id": id}).
		ToSql()

	if err != nil {
		err = fmt.Errorf("queryBuilder error: %v", err)
		log.Println(err)
		return
	}

	rows, err := app.db.Queryx(query, args...)
	if err != nil {
		err = fmt.Errorf("db.Queryx error: %v", err)
		log.Println(err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err == sql.ErrNoRows {
			return user, nil
		}
	}

	return
}

func (app App) getUserByEmail(email string) (user UserReturnData, err error) {
	query, args, err := queryBuilder.
		Select("id, username, email, role, created_at, updated_at").From(usersTable).
		Where(squirrel.Eq{"email": email}).
		ToSql()

	if err != nil {
		err = fmt.Errorf("queryBuilder error: %v", err)
		log.Println(err)
		return
	}

	rows, err := app.db.Queryx(query, args...)
	if err != nil {
		err = fmt.Errorf("db.Queryx error: %v", err)
		log.Println(err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err == sql.ErrNoRows {
			return user, nil
		}
	}

	return
}

func (app App) addUser(newUer UserSchema) (user UserReturnData, err error) {
	now := time.Now().UTC()
	columns := []string{"id", "email", "password", "created_at", "updated_at"}
	values := []interface{}{
		&newUer.ID, &newUer.Email, &newUer.Password, now, now,
	}

	query, args, err := queryBuilder.
		Insert(usersTable).
		Columns(columns...).
		Values(values...).
		Suffix("ON CONFLICT DO NOTHING").
		Suffix("RETURNING \"id\", \"email\", \"role\", \"created_at\", \"updated_at\"").
		ToSql()
	if err != nil {
		err = fmt.Errorf("queryBuilder error: %v", err)
		log.Println(err)
		return
	}

	rows, err := app.db.Queryx(query, args...)
	if err != nil {
		log.Println(err.Error())
		return
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Email, &user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			log.Println(err.Error())
			return user, err
		}
	}

	return
}

func (app App) getUserForLogin(email string) (user UserSchema, err error) {
	query, args, err := queryBuilder.
		Select("id, email, password, created_at, updated_at").From(usersTable).
		Where(squirrel.Eq{"email": email}).
		ToSql()

	if err != nil {
		err = fmt.Errorf("queryBuilder error: %v", err)
		log.Println(err)
		return
	}

	err = app.db.QueryRow(query, args...).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func (app App) updateUser(id string, username string) (updated int64, err error) {
	query, args, err := queryBuilder.
		Update(usersTable).
		Set("username", username).
		Where(squirrel.Eq{"id": id}).
		Where(squirrel.NotEq{"username": username}).
		ToSql()

	if err != nil {
		log.Println(err)
		return
	}

	result, err := app.db.Exec(query, args...)
	if err != nil {
		log.Println(err)
		return
	}
	updated, err = result.RowsAffected()
	log.Println(err)

	return
}
