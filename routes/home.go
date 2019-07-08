package routes

var homeRoutes []Route

// GetHomeRoutes .
func GetHomeRoutes() []Route {
	homeRoutes = append(homeRoutes,
		Route{
			"Index",
			"GET",
			"/",
			controller.GetHome(),
		},
	)

	return homeRoutes
}
