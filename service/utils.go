package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func (app App) logger(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf(
			"%s\t%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			r.UserAgent(),
			time.Since(start),
		)
		inner.ServeHTTP(w, r)
	})
}

func (app App) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (app App) respondWithError(w http.ResponseWriter, code int, errorData interface{}) {
	app.respondWithJSON(w, code, RootPayload{Error: true, Payload: errorData})
}

func (app App) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (app App) checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (app App) generateToken(user UserReturnData) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	fmt.Println(user)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":        user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"role":      user.Role,
		"createdAt": user.CreatedAt,
	})
	tokenString, err := token.SignedString([]byte(jwtSecret))

	return tokenString, err
}
