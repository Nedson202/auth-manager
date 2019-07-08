package middlewares

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	
	"github.com/asaskevich/govalidator"
	
	"github.com/nedson202/user-service/config"
	"github.com/nedson202/user-service/models"
	"github.com/mitchellh/mapstructure"
)

// ValidateAuthInput _
func ValidateAuthInput(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		routeName := ""
    if route := mux.CurrentRoute(req); route != nil {
			routeName = route.GetName()
    }
		var authSchema interface{}
		userLoginInput := &models.UserLoginSchema{}
		userSignupInput := &models.UserSignupSchema{}

		authSchema = userLoginInput
		if routeName == "AddUser" {
			authSchema = userSignupInput
		}
		
		json.NewDecoder(req.Body).Decode(&authSchema)

		_, err := govalidator.ValidateStruct(authSchema)
		if err != nil {
			config.RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		context.Set(req, "decoded", authSchema)
		next(w, req)
	})
}

// ValidateUserUpdate _
func ValidateUserUpdate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		userData := context.Get(req, "decoded")

		user := models.UserSchema{}
		mapstructure.Decode(userData, &user)

		userUpdateInput := &models.UsernameUpdateSchema{}

		json.NewDecoder(req.Body).Decode(&userUpdateInput)
		mapstructure.Decode(userUpdateInput, &user)

		_, err := govalidator.ValidateStruct(userUpdateInput)
		if err != nil {
			config.RespondWithError(w, http.StatusNotFound, err.Error())
			return
		}

		context.Set(req, "decoded", user)
		next(w, req)
	})
}
