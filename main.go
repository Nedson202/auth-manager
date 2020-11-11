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

	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"github.com/subosito/gotenv"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"

	"github.com/nedson202/auth-manager/auth_service_rest"
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
	pgDbName := os.Getenv("PG_DATABASE")

	log.Println("Connecting to Database")
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		pgHost,
		pgPort,
		pgUser,
		pgPassword,
		pgDbName,
	))
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Connected to Database")
	return db
}

func initCachePool() *redis.Pool {
	cacheHost := os.Getenv("REDIS_HOST")
	cachePort := os.Getenv("REDIS_PORT")

	cacheAddress := fmt.Sprintf("%s:%s", cacheHost, cachePort)
	log.Println("Creating caching pool: PING...")
	pool := newPool(cacheAddress)

	s, err := pingCachePool(pool)
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.Println(fmt.Sprintf("Created caching pool: %s", s))
	return pool
}

func newPool(address string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 5 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", address) },
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func pingCachePool(pool *redis.Pool) (s string, err error) {
	conn := pool.Get()
	defer conn.Close()

	s, err = redis.String(conn.Do("PING"))
	if err != nil {
		return
	}

	return
}

func main() {
	jwtSecret := os.Getenv("JWT_SECRET")
	db := initDatabase()

	router := mux.NewRouter()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	})

	_, err := auth_service_rest.New(router, db, jwtSecret, nil)
	if err != nil {
		log.Println(err)
	}

	server := &http.Server{
		// launch server with CORS validations
		Handler:      c.Handler(router),
		Addr:         ":5000",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	// Start Server
	func() {
		log.Println("Server started on http://localhost:5000")

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
