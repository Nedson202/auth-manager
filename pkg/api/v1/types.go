package v1

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

type UserSchema struct {
	Id        string     `json:"id"`
	Username  NullString `json:"username"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

type UserWithoutPassword struct {
	Id        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
