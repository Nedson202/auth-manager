package routes

import (
	"github.com/nedson202/user-service/controllers"
	"github.com/nedson202/user-service/middlewares"
)

var controller = &controllers.Controller{}

var userRoutes []Route

// GetUserRoutes .
func GetUserRoutes() []Route {
	addUserHandlers := middlewares.MultipleMiddleware(
		controller.AddUser(),
		middlewares.ValidateAuthInput,
		middlewares.VerifyUser,
	)

	loginUserHandlers := middlewares.MultipleMiddleware(
		controller.LoginUser(),
		middlewares.ValidateAuthInput,
	)

	getUserHandlers := middlewares.MultipleMiddleware(
		controller.GetUser(),
		middlewares.ValidateToken,
	)

	updateUserHandlers := middlewares.MultipleMiddleware(
		controller.UpdateUser(),
		middlewares.ValidateToken,
		middlewares.VerifyUser,
		middlewares.ValidateUserUpdate,
	)

	deleteUserHandlers := middlewares.MultipleMiddleware(
		controller.DeleteUser(),
		middlewares.ValidateToken,
		middlewares.VerifyUser,
	)

	userRoutes = append(userRoutes,
		Route{
			"GetUsers",
			"GET",
			"/api/v1/users",
			controller.GetUsers(),
		},
		Route{
			"GetUser",
			"GET",
			"/api/v1/users/profile",
			getUserHandlers,
		},
		Route{
			"AddUser",
			"POST",
			"/api/v1/users/signup",
			addUserHandlers,
		},
		Route{
			"LoginUser",
			"POST",
			"/api/v1/users/login",
			loginUserHandlers,
		},
		Route{
			"UpdateBook",
			"PUT",
			"/api/v1/users",
			updateUserHandlers,
		},
		Route{
			"DeleteUser",
			"DELETE",
			"/api/v1/users/{id}",
			deleteUserHandlers,
		},
	)

	return userRoutes
}
