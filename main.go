package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"database/sql"

	"github.com/subosito/gotenv"

	"github.com/gorilla/handlers"
	"github.com/asaskevich/govalidator"

	"github.com/nedson202/user-service/driver"
	"github.com/nedson202/user-service/config"
	"github.com/nedson202/user-service/routes"
	"github.com/nedson202/user-service/database"
)

var db *sql.DB

func init() {
	err := gotenv.Load()
	config.LogFatal(err)
}

func init() {
  govalidator.SetFieldsRequiredByDefault(false)
}

func main() {
	db = driver.DB
	database.MigrateDatabaseTables(db)

	router := routes.NewRouter()
	
	allowedOrigins := handlers.AllowedOrigins([]string{"*"}) 
  allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})
	
	port := os.Getenv("PORT")

	server := &http.Server{
		// launch server with CORS validations
		Handler:      handlers.CORS(allowedOrigins, allowedMethods)(router),
		Addr:         port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Start Server
	func() {
		log.Println("Starting Server on http://localhost:7000")
		if err := server.ListenAndServe(); err != nil {
			config.LogFatal(err)
		}
	}()

	handleShutdown(server)
}

// Handle graceful shutdown
func handleShutdown(server *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-interruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	server.Shutdown(ctx)

	log.Println("Shutting down")
	os.Exit(0)
}
