package controllers

import (
	"fmt"
	// "log"
	"net/http"

	"github.com/gorilla/context"
	// "github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"

	"github.com/nedson202/user-service/config"
	"github.com/nedson202/user-service/models"
	"github.com/nedson202/user-service/repository/user"
)

var users []models.UserSchema

var userRepo = &userrepository.UserRepository{}

// GetUsers handler
func (c *Controller) GetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var user models.UserSchema

		users = userRepo.GetUsers(user)

		if len(users) == 0 {
			config.RespondWithError(w, http.StatusNotFound, 
			"No user registed yet")
		} else {
			config.RespondWithJSON(w, http.StatusOK,
				config.DataPayload{
					Error: false,
					Message: "User list successfully retrieved",
					Data: users,
				},
			)
		}
	}
}

// GetUser controller to retrieve a user from db
func (c *Controller) GetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		userData := context.Get(req, "decoded")
		userSchema := models.UserSchema{}
		user := models.UserReturnData{}
		mapstructure.Decode(userData, &userSchema)

		fmt.Println(user)
		user = userRepo.GetUser(user, userSchema.ID)

		if user.ID == 0 {
			config.RespondWithError(w, http.StatusNotFound, "User not found")
		} else {
			config.RespondWithJSON(w, http.StatusOK,
				config.DataPayload{
					Error: false,
					Message: "User successfully retrieved",
					Data: user,
				},
			)
		}
	}
}

// AddUser controller to add a book to db
func (c *Controller) AddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		userData := context.Get(req, "decoded")

		user := models.UserReturnData{}
		userInput  := models.UserSchema{}
		var token string

		mapstructure.Decode(userData, &userInput)

		hashedPassword, err := config.HashPassword(userInput.Password)
		userInput.Password = hashedPassword

		user, err = userRepo.AddUser(userInput, user)

		if user.ID != 0 {
			token, err = config.GenerateToken(user)
		}

		user.Token = token;

		if user.ID == 0 {
			config.RespondWithError(w, http.StatusNotFound, err)
		} else {
			config.RespondWithJSON(w, http.StatusCreated,
				config.DataPayload{
					Error: false,
					Message: "User successfully created",
					Data: user,
				},
			)
		}
	}
}

// LoginUser controller to add a book to db
func (c *Controller) LoginUser() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		userData := context.Get(req, "decoded")

		user := models.UserSchema{}
		userInput  := models.UserSchema{}

		userReturnData := models.UserReturnData{}

		var token string

		mapstructure.Decode(userData, &userInput)

		user, _ = userRepo.LoginUser(userInput, user)

		isPassword := config.CheckPasswordHash(userInput.Password, user.Password)
		
		mapstructure.Decode(user, &userReturnData)

		if isPassword != true {
			config.RespondWithError(w, http.StatusNotFound, "User credentials provided not valid")
			return
		}
		token, _ = config.GenerateToken(userReturnData)

		userReturnData.Token = token;
			
		config.RespondWithJSON(w, http.StatusCreated,
			config.DataPayload{
				Error: false,
				Message: "User successfully logged in",
				Data: userReturnData,
			},
		)
	}
}

// UpdateUser controller to update user record on db
func (c *Controller) UpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		userData := context.Get(req, "decoded")

		user := models.UserSchema{}
		mapstructure.Decode(userData, &user)

		rowsUpdated := userRepo.UpdateUser(user)

		if rowsUpdated == 1 {
			config.RespondWithJSON(w, http.StatusOK,
				config.RootPayload{Error: false, Payload: "User successfully updated"})
		} else {
			config.RespondWithError(w, http.StatusInternalServerError,
				"An error occurred trying to update user")
		}
	}
}

// DeleteUser controller to delete book from db
func (c *Controller) DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		userData := context.Get(req, "decoded")

		user := models.UserSchema{}
		mapstructure.Decode(userData, &user)

		rowsDeleted := userRepo.DeleteUser(user.ID)
		if rowsDeleted == 1 {
			config.RespondWithJSON(w, http.StatusOK,
				config.RootPayload{Error: false, Payload: "User successfully deleted"})
		} else {
			config.RespondWithError(w, http.StatusInternalServerError,
				"An error occurred trying to delete user")
		}
	}
}
