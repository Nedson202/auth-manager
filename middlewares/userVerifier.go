package middlewares

import (
	"github.com/mitchellh/mapstructure"
	"net/http"
	"github.com/gorilla/context"

	"github.com/nedson202/user-service/config"
	"github.com/nedson202/user-service/models"
	"github.com/nedson202/user-service/repository/user"
)

var userRepo = userrepository.UserRepository{}

// VerifyUser _
func VerifyUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestMethod := req.Method

		userData := context.Get(req, "decoded")
		userInput := models.UserReturnData{}
		userSchema := models.UserSchema{}
		
		mapstructure.Decode(userData, &userSchema)
		
		userInput = userRepo.GetUser(userInput, userSchema.Email)

		if requestMethod == "DELETE" || requestMethod == "PUT" {
			if userInput.ID == 0 {
				config.RespondWithError(w, http.StatusNotFound, "User does not exist")
				return
			}
		} else if userInput.ID != 0 {
			config.RespondWithError(w, http.StatusNotFound, "User already exists")
			return
		}
		next(w, req)
	})
}