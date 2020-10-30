package service

import (
	"net/http"

	"github.com/jmoiron/sqlx"
)

// App _
type App struct {
	db *sqlx.DB
}

// Route defines a structure for routes
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// RootPayload _
type RootPayload struct {
	Error   bool        `json:"error"`
	Payload interface{} `json:"payload"`
}

// DataPayload structure for error responses
type DataPayload struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Middleware _
type Middleware func(http.HandlerFunc) http.HandlerFunc
