package v1

import (
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/mitchellh/mapstructure"
	v1 "github.com/nedson202/auth-manager/api/proto/v1"
	"github.com/nedson202/auth-manager/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var queryBuilder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Repository interface {
	GetUserByEmail(email string) (*v1.User, error)
	CreateUser(user *v1.User) (*v1.User, error)
}

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db}
}

const userTable = "identity.user"

func (r *PostgresRepository) CreateUser(newUer *v1.User) (*v1.User, error) {
	var userSchema UserSchema
	var user v1.User
	columns := []string{"id", "email", "password"}
	values := []interface{}{
		&newUer.Id, &newUer.Email, &newUer.Password,
	}

	query, args, err := queryBuilder.
		Insert(userTable).
		Columns(columns...).
		Values(values...).
		Suffix("ON CONFLICT DO NOTHING").
		Suffix("RETURNING \"id\", \"email\", \"created_at\", \"updated_at\"").
		ToSql()
	if err != nil {
		logger.Log.Error("repository:CreateUser:::queryBuilder error: " + err.Error())
		return &user, err
	}

	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		logger.Log.Error("repository:CreateUser::: " + err.Error())
		return &user, err
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&userSchema.Id, &userSchema.Email, &userSchema.CreatedAt, &userSchema.UpdatedAt)
		if err != nil {
			logger.Log.Error("repository:CreateUser::: " + err.Error())
			return &user, err
		}
	}

	mapstructure.Decode(userSchema, &user)
	return &user, nil
}

func (r *PostgresRepository) GetUserByEmail(email string) (*v1.User, error) {
	var userSchema UserSchema
	var user v1.User
	query, args, err := queryBuilder.
		Select("id", "username", "email", "password", "created_at", "updated_at").
		From(userTable).
		Where(squirrel.Eq{"email": email}).
		ToSql()

	if err != nil {
		err = fmt.Errorf("queryBuilder error: %v", err)
		logger.Log.Error("repository:GetUserByEmail:::queryBuilder error: " + err.Error())
		return &user, status.Error(codes.Unknown, "failed to retrieve user from identity.user-> "+err.Error())
	}

	row := r.db.QueryRowx(query, args...)
	if err := row.Scan(&userSchema.Id, &userSchema.Username, &userSchema.Email, &userSchema.Password, &userSchema.CreatedAt, &userSchema.UpdatedAt); err != nil {
		logger.Log.Error("repository:GetUserByEmail:::failed to retrieve field values from identity.user row-> " + err.Error())
		return nil, status.Error(codes.Unknown, "failed to retrieve field values from identity.user row-> "+err.Error())
	}

	mapstructure.Decode(userSchema, &user)
	return &user, nil
}
