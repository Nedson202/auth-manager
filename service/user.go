package service

import (
	"time"
)

// UserSchema defines structure of user information
type UserSchema struct {
	ID        int       `json:"id"`
	Username  string    `json:"username" valid:"required~Username is blank"`
	Email     string    `json:"email" valid:"email"`
	Password  string    `json:"password" valid:"length(5|16)~Password length should not be less than 5"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt"`
}

// UserSignupSchema defines structure of user information
type UserSignupSchema struct {
	Username string `json:"username" valid:"required~Username is blank"`
	Email    string `json:"email" valid:"email"`
	Password string `json:"password" valid:"length(5|16)~Password length should not be less than 5"`
}

// UserLoginSchema defines structure of user information
type UserLoginSchema struct {
	Email    string `json:"email" valid:"email~Email is required"`
	Password string `json:"password" valid:"length(5|16)~Password length should not be less than 5"`
}

// UsernameUpdateSchema defines structure of user information
type UsernameUpdateSchema struct {
	Username string `json:"username" valid:"required~Username is blank"`
}

// UserReturnData defines structure of user information to return
type UserReturnData struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Token     string    `json:"token"`
}
