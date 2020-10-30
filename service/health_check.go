package service

import (
	"net/http"
)

func (app App) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	app.respondWithJSON(w, http.StatusOK,
		RootPayload{
			Error:   false,
			Payload: "User management API running",
		},
	)
}
