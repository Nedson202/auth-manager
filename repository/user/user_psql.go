package userrepository

import (
	"fmt"
	"reflect"
	"log"
	"database/sql"

	"github.com/nedson202/user-service/config"
	"github.com/nedson202/user-service/models"
	"github.com/nedson202/user-service/driver"
)

var db *sql.DB 

// GetUsers method
func (b *UserRepository) GetUsers(user models.UserSchema) []models.UserSchema {
	var users []models.UserSchema
	
	db = driver.DB
	query := `
		SELECT id, username, email, password, role, created_at, updated_at
		FROM users
		WHERE users.deleted IS NOT TRUE
	`
	rows, err := db.Query(query)
	config.LogFatal(err)

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password,
			&user.Role, &user.CreatedAt, &user.UpdatedAt)
		config.LogFatal(err)

		users = append(users, user)
	}

	return users
}

func (b *UserRepository) GetUser(user models.UserReturnData, id interface{}) (models.UserReturnData) {
	db = driver.DB
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
	rows := db.QueryRow(query, id)

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

func (b *UserRepository) AddUser(userData models.UserSchema, user models.UserReturnData) (models.UserReturnData, error) {
	db = driver.DB
	query := `
		INSERT into users (username, email, password)
		VALUES ($1, $2, $3)
		RETURNING id, username, email, role, created_at, updated_at;
	`
	err := db.QueryRow(
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

func (b *UserRepository) LoginUser(userData models.UserSchema, user models.UserSchema) (models.UserSchema, error) {
	db = driver.DB
	query := `
		SELECT id, username, email, password, role, created_at
		FROM users
		WHERE email = $1
	`
	err := db.QueryRow(
		query,  &userData.Email,
	).Scan(&user.ID, &user.Username, &user.Email,
		&user.Password, &user.Role, &user.CreatedAt,)

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

func (b *UserRepository) UpdateUser(user models.UserSchema) int64 {
	db = driver.DB
	result, err := db.Exec("update users set username = $1 where id = $2 Returning id;",
	user.Username, user.ID)

	rowsUpdated, err := result.RowsAffected()
	config.LogFatal(err)

	return rowsUpdated
}

func (b *UserRepository) DeleteUser(id int) int64 {
	db = driver.DB
	query := `
		UPDATE users set deleted = TRUE, deleted_at = current_timestamp
		WHERE id = $1 and users.deleted IS NOT TRUE
	`
	result, err := db.Exec(query, id)

	rowsDeleted, err := result.RowsAffected()
	config.LogFatal(err)

	return rowsDeleted
}
