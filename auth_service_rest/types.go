package auth_service_rest

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gomodule/redigo/redis"

	"github.com/jmoiron/sqlx"
)

// App _
type App struct {
	db        *sqlx.DB
	jwtSecret []byte
	cachePool *redis.Pool
}

// Route defines a structure for routes
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// DataPayload structure for error responses
type DataPayload struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
}

// Middleware _
type Middleware func(http.HandlerFunc) http.HandlerFunc

// MalformedRequest _
type MalformedRequest struct {
	status int
	msg    string
}

func (mr *MalformedRequest) Error() string {
	return mr.msg
}

type JwtClaims struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	ExpiresAt time.Time `json:"expirationTime"`
	jwt.StandardClaims
}
