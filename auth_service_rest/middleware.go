package auth_service_rest

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
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

func (app App) validateRequest(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		err := app.decodeJSONBody(w, req)
		if err != nil {
			app.processMalformedRequestError(err, w, req)
			return
		}
		next(w, req)
	})
}

func (app App) validateToken(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")

		bearerToken := strings.Split(authorizationHeader, " ")

		if authorizationHeader == "" || len(bearerToken) == 1 {
			app.respondWithJSON(w, http.StatusUnauthorized,
				DataPayload{
					Success: false,
					Data: map[string]interface{}{
						"message": "An authorization header is required",
					},
				},
			)
			return
		}

		token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("An error occured")
			}
			return app.jwtSecret, nil
		})

		if err != nil {
			app.respondWithJSON(w, http.StatusUnauthorized,
				DataPayload{
					Success: false,
					Data: map[string]interface{}{
						"message": "The token provided is invalid",
					},
				},
			)
			return
		}

		if token.Valid != true {
			app.respondWithJSON(w, http.StatusUnauthorized,
				DataPayload{
					Success: false,
					Data: map[string]interface{}{
						"message": "Authorization token is not valid",
					},
				},
			)
			return
		}

		context.Set(req, "decoded", token.Claims)
		next(w, req)
	})
}

func (app App) verifyUpdateDetails(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		body := context.Get(req, "body")
		userSchema := UserSchema{}
		mapstructure.Decode(body, &userSchema)

		tokenClaims := context.Get(req, "decoded")
		mapstructure.Decode(tokenClaims, &userSchema)

		_, err := app.getUserByID(userSchema.ID)
		if err != nil {
			app.respondWithJSON(w, http.StatusNotFound,
				DataPayload{
					Success: false,
					Data: map[string]interface{}{
						"message": "Could not find user details. Please contact support",
					},
				},
			)
			return
		}

		next(w, req)
	})
}

func (app App) validateAuthInput(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		body := context.Get(req, "body")
		authSchema := AuthSchema{}
		mapstructure.Decode(body, &authSchema)

		_, err := govalidator.ValidateStruct(authSchema)
		if err != nil {
			app.respondWithJSON(w, http.StatusBadRequest,
				DataPayload{
					Success: false,
					Data: map[string]interface{}{
						"message":       "Incomplete fields provided",
						"missingFields": err.Error(),
					},
				},
			)
			return
		}

		next(w, req)
	})
}

func (app App) validateUserUpdate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		body := context.Get(req, "body")
		user := UsernameUpdateSchema{}
		mapstructure.Decode(body, &user)

		_, err := govalidator.ValidateStruct(user)
		if err != nil {
			app.respondWithJSON(w, http.StatusBadRequest,
				DataPayload{
					Success: false,
					Data: map[string]interface{}{
						"message":       "Incomplete fields provided",
						"missingFields": err.Error(),
					},
				},
			)
			return
		}

		next(w, req)
	})
}
