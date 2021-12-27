package v1

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	// pg driver
	_ "github.com/lib/pq"
	"github.com/nedson202/auth-manager/pkg/logger"
)

var (
	GRPC_PORT   string
	HTTP_PORT   string
	PG_HOST     string
	PG_PORT     string
	PG_USER     string
	PG_PASSWORD string
	JWT_SECRET  string
)

// Config is configuration for Server
type Config struct {
	// gRPC server start parameters section
	// gRPC is TCP port to listen by gRPC server
	GRPCPort string

	// HTTP/REST gateway start parameters section
	// HTTPPort is TCP port to listen by HTTP/REST gateway
	HTTPPort string

	// DB Datastore parameters section
	// DatastoreDBHost is host of database
	DatastoreDBHost string
	// DatastoreDBPort is port to connect to database
	DatastoreDBPort string
	// DatastoreDBUser is username to connect to database
	DatastoreDBUser string
	// DatastoreDBPassword password to connect to database
	DatastoreDBPassword string

	JwtSecret string

	// Log parameters section
	// LogLevel is global log level: Debug(-1), Info(0), Warn(1), Error(2), DPanic(3), Panic(4), Fatal(5)
	LogLevel int
	// LogTimeFormat is print time format for logger e.g. 2006-01-02T15:04:05Z07:00
	LogTimeFormat string
}

func initDatabase(cfg Config) *sqlx.DB {
	log.Println("Connecting to Database")
	db, err := sqlx.Connect("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		cfg.DatastoreDBHost,
		cfg.DatastoreDBPort,
		cfg.DatastoreDBUser,
		cfg.DatastoreDBPassword,
	))
	if err != nil {
		log.Fatalln(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to Database")
	db.SetMaxOpenConns(20)
	return db
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	GRPC_PORT = os.Getenv("GRPC_PORT")
	HTTP_PORT = os.Getenv("HTTP_PORT")
	PG_HOST = os.Getenv("PG_HOST")
	PG_PORT = os.Getenv("PG_PORT")
	PG_USER = os.Getenv("PG_USER")
	PG_PASSWORD = os.Getenv("PG_PASSWORD")
	JWT_SECRET = os.Getenv("JWT_SECRET")

	// get configuration
	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", GRPC_PORT, "gRPC port to bind")
	flag.StringVar(&cfg.HTTPPort, "http-port", HTTP_PORT, "HTTP port to bind")
	flag.StringVar(&cfg.DatastoreDBHost, "db-host", PG_HOST, "Database host")
	flag.StringVar(&cfg.DatastoreDBPort, "db-port", PG_PORT, "Database port")
	flag.StringVar(&cfg.DatastoreDBUser, "db-user", PG_USER, "Database user")
	flag.StringVar(&cfg.DatastoreDBPassword, "db-password", PG_PASSWORD, "Database password")
	flag.StringVar(&cfg.JwtSecret, "jwt-secret", JWT_SECRET, "JWT secret")
	flag.IntVar(&cfg.LogLevel, "log-level", 0, "Global log level")
	flag.StringVar(&cfg.LogTimeFormat, "log-time-format", "", "Print time format for logger e.g. 2006-01-02T15:04:05Z07:00")
	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}

	if len(cfg.HTTPPort) == 0 {
		return fmt.Errorf("invalid TCP port for HTTP gateway: '%s'", cfg.HTTPPort)
	}

	// initialize logger
	if err := logger.Init(cfg.LogLevel, cfg.LogTimeFormat); err != nil {
		return fmt.Errorf("failed to initialize logger: %v", err)
	}

	db := initDatabase(cfg)
	defer db.Close()

	tokenService := NewTokenService(cfg)
	repo := NewPostgresRepository(db)
	v1API := NewAuthServiceServer(repo, tokenService)

	// run HTTP gateway
	go func() {
		_ = RunRestServer(ctx, cfg.GRPCPort, cfg.HTTPPort)
	}()

	return RunGrpcServer(ctx, v1API, cfg.GRPCPort)
}
