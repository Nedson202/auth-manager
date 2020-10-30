package service

import (
	"fmt"
	"net/http"

	"github.com/gorilla/context"
	"github.com/mitchellh/mapstructure"
)

var users []UserSchema

func (app App) getUsersHandler(w http.ResponseWriter, req *http.Request) {
	var user UserSchema

	users = app.getUsers(user)

	if len(users) == 0 {
		app.respondWithError(w, http.StatusNotFound,
			"No user registed yet")
	} else {
		app.respondWithJSON(w, http.StatusOK,
			DataPayload{
				Error:   false,
				Message: "User list successfully retrieved",
				Data:    users,
			},
		)
	}
}

func (app App) getUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		userData := context.Get(req, "decoded")
		userSchema := UserSchema{}
		user := UserReturnData{}
		mapstructure.Decode(userData, &userSchema)

		fmt.Println(user)
		user = app.getUser(user, userSchema.ID)

		if user.ID == 0 {
			app.respondWithError(w, http.StatusNotFound, "User not found")
		} else {
			app.respondWithJSON(w, http.StatusOK,
				DataPayload{
					Error:   false,
					Message: "User successfully retrieved",
					Data:    user,
				},
			)
		}
	}
}

func (app App) addUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		userData := context.Get(req, "decoded")

		user := UserReturnData{}
		userInput := UserSchema{}
		var token string

		mapstructure.Decode(userData, &userInput)

		hashedPassword, err := app.hashPassword(userInput.Password)
		userInput.Password = hashedPassword

		user, err = app.addUser(userInput, user)

		if user.ID != 0 {
			token, err = app.generateToken(user)
		}

		user.Token = token

		if user.ID == 0 {
			app.respondWithError(w, http.StatusNotFound, err)
		} else {
			app.respondWithJSON(w, http.StatusCreated,
				DataPayload{
					Error:   false,
					Message: "User successfully created",
					Data:    user,
				},
			)
		}
	}
}

func (app App) loginUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		userData := context.Get(req, "decoded")

		user := UserSchema{}
		userInput := UserSchema{}

		userReturnData := UserReturnData{}

		var token string

		mapstructure.Decode(userData, &userInput)

		user, _ = app.loginUser(userInput, user)

		isPassword := app.checkPasswordHash(userInput.Password, user.Password)

		mapstructure.Decode(user, &userReturnData)

		if isPassword != true {
			app.respondWithError(w, http.StatusNotFound, "User credentials provided not valid")
			return
		}
		token, _ = app.generateToken(userReturnData)

		userReturnData.Token = token

		app.respondWithJSON(w, http.StatusCreated,
			DataPayload{
				Error:   false,
				Message: "User successfully logged in",
				Data:    userReturnData,
			},
		)
	}
}

func (app App) updateUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		userData := context.Get(req, "decoded")

		user := UserSchema{}
		mapstructure.Decode(userData, &user)

		rowsUpdated := app.updateUser(user)

		if rowsUpdated == 1 {
			app.respondWithJSON(w, http.StatusOK, RootPayload{Error: false, Payload: "User successfully updated"})
		} else {
			app.respondWithError(w, http.StatusInternalServerError,
				"An error occurred trying to update user")
		}
	}
}

func (app App) deleteUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		userData := context.Get(req, "decoded")

		user := UserSchema{}
		mapstructure.Decode(userData, &user)

		rowsDeleted := app.deleteUser(user.ID)
		if rowsDeleted == 1 {
			app.respondWithJSON(w, http.StatusOK, RootPayload{Error: false, Payload: "User successfully deleted"})
		} else {
			app.respondWithError(w, http.StatusInternalServerError,
				"An error occurred trying to delete user")
		}
	}
}
