package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"github.com/subosito/gotenv"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"

	"github.com/nedson202/user-service/service"
)

func init() {
	err := gotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	govalidator.SetFieldsRequiredByDefault(false)
}

func initDatabase() *sqlx.DB {
	pgHost := os.Getenv("PG_HOST")
	pgPort := os.Getenv("PG_PORT")
	pgUser := os.Getenv("PG_USER")
	pgPassword := os.Getenv("PG_PASSWORD")

	log.Println("Connecting to Database")
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		pgHost,
		pgPort,
		pgUser,
		pgPassword,
	))
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Connected to Database")
	return db
}

func main() {
	db := initDatabase()

	router := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	})

	port := os.Getenv("PORT")

	_, err := service.New(router, db)
	if err != nil {
		log.Println(err)
	}

	combineServerAddress := fmt.Sprintf("%s%s", ":", port)
	server := &http.Server{
		// launch server with CORS validations
		Handler:      c.Handler(router),
		Addr:         combineServerAddress,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Start Server
	func() {
		startMessage := fmt.Sprintf("%s%s", "Server started on http://localhost:", port)
		log.Println(startMessage)

		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
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
