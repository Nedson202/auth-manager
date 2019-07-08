package middlewares

import (
	"os"
	"net/http"
	"strings"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"

	"github.com/nedson202/user-service/config"
)

// ValidateToken to validate authorization header
func ValidateToken(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		jwtSecret := os.Getenv("JWT_SECRET")
		authorizationHeader := req.Header.Get("authorization")

		bearerToken := strings.Split(authorizationHeader, " ")

		if authorizationHeader == "" || len(bearerToken) == 1 {
			config.RespondWithError(w, http.StatusNotFound,
				"An authorization header is required",
			)
			return
		}
		token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
			}
			return []byte(jwtSecret), nil
		})
		if err != nil {
			config.RespondWithError(w, http.StatusNotFound, "The token provided is invalid")
			return
		}
		if token.Valid {
			context.Set(req, "decoded", token.Claims)
			next(w, req)
		} else {
			config.RespondWithError(w, http.StatusNotFound, "Authorization token is not valid")
			return
		}
	})
}