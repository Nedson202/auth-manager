package auth_service_rest

import (
	"net/http"
)

func (app App) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	app.respondWithJSON(w, http.StatusOK,
		DataPayload{
			Success: true,
			Data: map[string]interface{}{
				"message": "Auth Management REST API running",
			},
		},
	)
}
