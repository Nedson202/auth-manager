package service

import (
	"net/http"

	"github.com/gorilla/mux"
)

var routes []Route
var healthCheckRoute Route
var userRoutes []Route

func (app App) getHealthCheckHandler() Route {
	healthCheckRoute = Route{
		"HealthCheck",
		"GET",
		"/api/v1/health",
		app.healthCheckHandler,
	}

	return healthCheckRoute
}

func (app App) getUserRoutes() []Route {
	addUserHandlers := app.multipleMiddleware(
		app.addUserHandler(),
		app.validateAuthInput,
		app.verifyUser,
	)

	loginUserHandlers := app.multipleMiddleware(
		app.loginUserHandler(),
		app.validateAuthInput,
	)

	getUserHandlers := app.multipleMiddleware(
		app.getUserHandler(),
		app.validateToken,
	)

	updateUserHandlers := app.multipleMiddleware(
		app.updateUserHandler(),
		app.validateToken,
		app.verifyUser,
		app.validateUserUpdate,
	)

	deleteUserHandlers := app.multipleMiddleware(
		app.deleteUserHandler(),
		app.validateToken,
		app.verifyUser,
	)

	userRoutes = append(userRoutes,
		Route{
			"GetUsers",
			"GET",
			"/api/v1/users",
			app.getUsersHandler,
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

func (app App) newRouter(router *mux.Router) *mux.Router {
	userRoutes := app.getUserRoutes()
	homeRoute := app.getHealthCheckHandler()

	routes = append(routes, homeRoute)
	routes = append(routes, userRoutes...)

	combineRoutes := router.StrictSlash(true)

	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		combineRoutes.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	router.Use(app.logger)

	return router
}
