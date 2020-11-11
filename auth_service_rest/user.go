package auth_service_rest

import (
	"database/sql"
	"encoding/json"
	"time"
)

type NullString struct {
	sql.NullString
}

// override MarshalJSON method for custom type NullString
func (value *NullString) MarshalJSON() ([]byte, error) {
	if !value.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(value.String)
}

// override UnmarshalJSON method for custom type NullString
func (value *NullString) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &value.String)
	value.Valid = (err == nil)

	return err
}

type User struct {
	ID        string     `json:"id"`
	Username  NullString `json:"username"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	Role      string     `json:"role"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// UserSchema defines structure of user information
type UserSchema struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// AuthSchema defines structure of auth information
type AuthSchema struct {
	Email    string `json:"email" valid:"required~Email is required"`
	Password string `json:"password" valid:"length(4|16)~Password length should not be less than 4"`
}

// UsernameUpdateSchema defines structure of user information
type UsernameUpdateSchema struct {
	Username string `json:"username" valid:"required~Username is blank"`
}

// UserReturnData defines structure of user information to return
type UserReturnData struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Token     string    `json:"token"`
}
