package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/nedson202/user-service/config"
)

var routes []Route

//NewRouter configures a new router to the API
func NewRouter() *mux.Router {
	userRoutes := GetUserRoutes()
	homeRoutes := GetHomeRoutes()

	routes = append(routes, userRoutes...)
	routes = append(routes, homeRoutes...)
	router := mux.NewRouter().StrictSlash(true)
	
	for _, route := range routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = config.Logger(handler, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
