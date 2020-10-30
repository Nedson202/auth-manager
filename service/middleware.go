package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
)

func (app App) multipleMiddleware(h http.HandlerFunc, m ...Middleware) http.HandlerFunc {
	if len(m) < 1 {
		return h
	}

	wrapped := h

	// loop in reverse to preserve middleware order
	for i := len(m) - 1; i >= 0; i-- {
		wrapped = m[i](wrapped)
	}

	return wrapped
}

func (app App) validateToken(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		jwtSecret := os.Getenv("JWT_SECRET")
		authorizationHeader := req.Header.Get("authorization")

		bearerToken := strings.Split(authorizationHeader, " ")

		if authorizationHeader == "" || len(bearerToken) == 1 {
			app.respondWithError(w, http.StatusNotFound,
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
			app.respondWithError(w, http.StatusNotFound, "The token provided is invalid")
			return
		}
		if token.Valid {
			context.Set(req, "decoded", token.Claims)
			next(w, req)
		} else {
			app.respondWithError(w, http.StatusNotFound, "Authorization token is not valid")
			return
		}
	})
}

func (app App) verifyUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestMethod := req.Method

		userData := context.Get(req, "decoded")
		userInput := UserReturnData{}
		userSchema := UserSchema{}

		mapstructure.Decode(userData, &userSchema)

		userInput = app.getUser(userInput, userSchema.Email)

		if requestMethod == "DELETE" || requestMethod == "PUT" {
			if userInput.ID == 0 {
				app.respondWithError(w, http.StatusNotFound, "User does not exist")
				return
			}
		} else if userInput.ID != 0 {
			app.respondWithError(w, http.StatusNotFound, "User already exists")
			return
		}
		next(w, req)
	})
}

func (app App) validateAuthInput(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		routeName := ""
		if route := mux.CurrentRoute(req); route != nil {
			routeName = route.GetName()
		}
		var authSchema interface{}
		userLoginInput := &UserLoginSchema{}
		userSignupInput := &UserSignupSchema{}

		authSchema = userLoginInput
		if routeName == "AddUser" {
			authSchema = userSignupInput
		}

		json.NewDecoder(req.Body).Decode(&authSchema)

		_, err := govalidator.ValidateStruct(authSchema)
		if err != nil {
			app.respondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		context.Set(req, "decoded", authSchema)
		next(w, req)
	})
}

func (app App) validateUserUpdate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		userData := context.Get(req, "decoded")

		user := UserSchema{}
		mapstructure.Decode(userData, &user)

		userUpdateInput := &UsernameUpdateSchema{}

		json.NewDecoder(req.Body).Decode(&userUpdateInput)
		mapstructure.Decode(userUpdateInput, &user)

		_, err := govalidator.ValidateStruct(userUpdateInput)
		if err != nil {
			app.respondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		context.Set(req, "decoded", user)
		next(w, req)
	})
}
