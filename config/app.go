package config

import (
	"fmt"
	"os"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/subosito/gotenv"

	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
	"github.com/nedson202/user-service/models"	
)

func init() {
	gotenv.Load()
}

// LogFatal to handle logging errors
func LogFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// Logger function
func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		inner.ServeHTTP(w, r)
		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}

// RespondWithError handler for sending errors over http
func RespondWithError(w http.ResponseWriter, code int, errorData interface{}) {
	RespondWithJSON(w, code, RootPayload{Error: true, Payload: errorData})
}

// RespondWithJSON handler for sending responses over http
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// HashPassword encrypts password provided by user
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash compares password and hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateToken generate jwt from user data
func GenerateToken(user models.UserReturnData) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	fmt.Println(user)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": user.ID,
		"username": user.Username,
		"email": user.Email,
		"role": user.Role,
		"createdAt": user.CreatedAt,
	})
	tokenString, err := token.SignedString([]byte(jwtSecret))

	return tokenString, err
}