package auth_service_rest

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
	createUserHandler := app.multipleMiddleware(
		app.createUserHandler(),
		app.validateRequest,
		app.validateAuthInput,
	)

	loginUserHandler := app.multipleMiddleware(
		app.loginUserHandler(),
		app.validateRequest,
		app.validateAuthInput,
	)

	getUserHandler := app.multipleMiddleware(
		app.getUserHandler(),
		app.validateRequest,
		app.validateToken,
	)

	updateUserHandler := app.multipleMiddleware(
		app.updateUserHandler(),
		app.validateRequest,
		app.validateToken,
		app.verifyUpdateDetails,
		app.validateUserUpdate,
	)

	tokenRefreshHandler := app.multipleMiddleware(
		app.refreshToken(),
		app.validateRequest,
		app.validateToken,
	)

	userRoutes = append(userRoutes,
		Route{
			"GetUser",
			"GET",
			"/api/v1/users/profile",
			getUserHandler,
		},
		Route{
			"AddUser",
			"POST",
			"/api/v1/users",
			createUserHandler,
		},
		Route{
			"LoginUser",
			"POST",
			"/api/v1/users/login",
			loginUserHandler,
		},
		Route{
			"UpdateUser",
			"PATCH",
			"/api/v1/users",
			updateUserHandler,
		},
		Route{
			"RefreshToken",
			"POST",
			"/api/v1/users/refresh",
			tokenRefreshHandler,
		},
	)

	return userRoutes
}

func (app App) newRouter(router *mux.Router) *mux.Router {
	userRoutes := app.getUserRoutes()
	healthCheck := app.getHealthCheckHandler()

	routes = append(routes, healthCheck)
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
