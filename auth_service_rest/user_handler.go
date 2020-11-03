package auth_service_rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/context"
	"github.com/mitchellh/mapstructure"
)

var users []UserSchema

func (app App) getUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		userSchema := UserSchema{}
		userData := context.Get(req, "body")
		tokenClaims := context.Get(req, "decoded")
		mapstructure.Decode(tokenClaims, &userSchema)
		mapstructure.Decode(userData, &userSchema)

		user, err := app.getUserByID(userSchema.ID)
		if err != nil || user.ID == "" {
			app.respondWithJSON(w, http.StatusNotFound,
				DataPayload{
					Success: false,
					Data: map[string]interface{}{
						"message": "User not found",
					},
				},
			)
			return
		}

		userReturnData := UserReturnData{}
		byte, err := json.Marshal(&user)
		err = json.Unmarshal(byte, &userReturnData)

		token, _ := app.generateToken(userReturnData)
		userReturnData.Token = token
		app.respondWithJSON(w, http.StatusOK,
			DataPayload{
				Success: false,
				Data: map[string]interface{}{
					"message": "User retrieved",
					"users":   userReturnData,
				},
			},
		)
	}
}

func (app App) createUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		userData := context.Get(req, "body")

		user := UserReturnData{}
		userInput := UserSchema{}
		var token string

		mapstructure.Decode(userData, &userInput)

		hashedPassword, _ := app.hashPassword(userInput.Password)
		userInput.Password = hashedPassword
		userInput.ID = app.getUUID()

		user, err := app.addUser(userInput)
		if err != nil {
			app.respondWithJSON(w, http.StatusInternalServerError,
				DataPayload{
					Success: false,
					Data: map[string]interface{}{
						"message": "Failed to create user",
					},
				},
			)
			return
		}

		if user.ID == "" {
			app.respondWithJSON(w, http.StatusConflict,
				DataPayload{
					Success: false,
					Data: map[string]interface{}{
						"message": "User already exists. Please provide another email",
					},
				},
			)
			return
		}

		token, _ = app.generateToken(user)
		user.Token = token
		app.respondWithJSON(w, http.StatusCreated,
			DataPayload{
				Success: true,
				Data: map[string]interface{}{
					"message": "User created",
					"user":    user,
				},
			},
		)
		return
	}
}

func (app App) loginUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		userData := context.Get(req, "body")

		var token string
		userSchema := UserSchema{}
		userReturnData := UserReturnData{}

		mapstructure.Decode(userData, &userSchema)

		user, _ := app.getUserForLogin(userSchema.Email)
		isPassword := app.checkPasswordHash(userSchema.Password, user.Password)

		mapstructure.Decode(user, &userReturnData)
		if isPassword != true {
			app.respondWithJSON(w, http.StatusUnauthorized,
				DataPayload{
					Success: false,
					Data: map[string]interface{}{
						"message": "Invalid auth credentials provided",
					},
				},
			)
			return
		}

		token, _ = app.generateToken(userReturnData)
		app.respondWithJSON(w, http.StatusOK,
			DataPayload{
				Success: true,
				Data: map[string]interface{}{
					"message": "Login successful",
					"user": map[string]interface{}{
						"id":    user.ID,
						"email": user.Email,
						"token": token,
					},
				},
			},
		)
		return
	}
}

func (app App) updateUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		user := UserReturnData{}
		tokenClaims := context.Get(req, "decoded")
		mapstructure.Decode(tokenClaims, &user)

		userData := context.Get(req, "body")
		mapstructure.Decode(userData, &user)

		rowsUpdated, err := app.updateUser(user.ID, user.Username)
		if err == nil && rowsUpdated == 1 {
			token, _ := app.generateToken(user)
			user.Token = token
			app.respondWithJSON(w, http.StatusOK,
				DataPayload{
					Success: true,
					Data: map[string]interface{}{
						"message": "User updated",
						"user":    user,
					},
				},
			)
			return
		}

		if err == nil && rowsUpdated == 0 {
			app.respondWithJSON(w, http.StatusConflict,
				DataPayload{
					Success: false,
					Data: map[string]interface{}{
						"message": "Username provided is no longer available",
					},
				},
			)
			return
		}

		app.respondWithJSON(w, http.StatusInternalServerError,
			DataPayload{
				Success: false,
				Data: map[string]interface{}{
					"message": "Failed to update username",
				},
			},
		)
		return
	}
}

func (app App) refreshToken() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		user := UserReturnData{}
		tokenClaims := context.Get(req, "decoded")
		mapstructure.Decode(tokenClaims, &user)

		token, err := app.generateToken(user)
		if err != nil {
			app.respondWithJSON(w, http.StatusInternalServerError,
				DataPayload{
					Success: false,
					Data: map[string]interface{}{
						"message": "Token refresh failed",
					},
				},
			)
			return
		}

		user.Token = token
		app.respondWithJSON(w, http.StatusOK,
			DataPayload{
				Success: true,
				Data: map[string]interface{}{
					"message": "Token refreshed",
					"user":    user,
				},
			},
		)
	}
}
