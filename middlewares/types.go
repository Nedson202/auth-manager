package middlewares

import (
	"net/http"
)

// Middleware _
type Middleware func( http.HandlerFunc) http.HandlerFunc
