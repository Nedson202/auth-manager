package service

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
)

func (app App) getUsers(user UserSchema) []UserSchema {
	var users []UserSchema

	query := `
		SELECT id, username, email, password, role, created_at, updated_at
		FROM users
		WHERE users.deleted IS NOT TRUE
	`
	rows, err := app.db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password,
			&user.Role, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			log.Fatal(err)
		}

		users = append(users, user)
	}

	return users
}

func (app App) getUser(user UserReturnData, id interface{}) UserReturnData {
	queryType := `users.id = $1`
	typeOfID := reflect.TypeOf(id).Kind()

	if typeOfID == reflect.String {
		queryType = `users.email = $1`
	}

	query := fmt.Sprintf(
		`
		SELECT id, username, email, role, created_at, updated_at
		FROM users
		WHERE %s and users.deleted IS NOT TRUE
		`, queryType,
	)
	rows := app.db.QueryRow(query, id)

	err := rows.Scan(&user.ID, &user.Username, &user.Email,
		&user.Role, &user.CreatedAt, &user.UpdatedAt)

	switch {
	case err == sql.ErrNoRows:
		log.Printf("User not found")
	case err != nil:
		log.Println(err)
	default:
		return user
	}

	return user
}

func (app App) addUser(userData UserSchema, user UserReturnData) (UserReturnData, error) {
	query := `
		INSERT into users (username, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, username, email, role, created_at, updated_at;
	`
	err := app.db.QueryRow(
		query,
		&userData.Username, &userData.Email, &userData.Password,
	).Scan(&user.ID, &user.Username, &user.Email,
		&user.Role, &user.CreatedAt, &user.UpdatedAt)

	switch {
	case err == sql.ErrNoRows:
		log.Printf("Error occurred")
	case err != nil:
		log.Println(err)
	default:
		return user, err
	}

	return user, err
}

func (app App) loginUser(userData UserSchema, user UserSchema) (UserSchema, error) {
	query := `
		SELECT id, username, email, password, role, created_at
		FROM users
		WHERE email = $1
	`
	err := app.db.QueryRow(
		query, &userData.Email,
	).Scan(&user.ID, &user.Username, &user.Email,
		&user.Password, &user.Role, &user.CreatedAt)

	switch {
	case err == sql.ErrNoRows:
		log.Printf("Error occurred")
	case err != nil:
		log.Println(err)
	default:
		return user, err
	}

	return user, err
}

func (app App) updateUser(user UserSchema) int64 {
	result, err := app.db.Exec("update users set username = $1 where id = $2 Returning id;",
		user.Username, user.ID)

	rowsUpdated, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	return rowsUpdated
}

func (app App) deleteUser(id int) int64 {
	query := `
		UPDATE users set deleted = TRUE, deleted_at = current_timestamp
		WHERE id = $1 and users.deleted IS NOT TRUE
	`
	result, err := app.db.Exec(query, id)

	rowsDeleted, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	return rowsDeleted
}
