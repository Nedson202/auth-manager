package controllers

import (
	"net/http"
	"github.com/nedson202/user-service/config"
)

// GetHome handler
func (c *Controller) GetHome() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		config.RespondWithJSON(w, http.StatusOK,
			config.RootPayload{
				Error: false,
				Payload: "Welcome to the user management service",
			},
		)
	}
}
