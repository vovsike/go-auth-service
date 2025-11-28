package main

import (
	"awesomeProject/internal/database"
	"awesomeProject/internal/user"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	openfga "github.com/openfga/go-sdk"
	openfgaClient "github.com/openfga/go-sdk/client"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// Create database connection pool
	pool, err := createDBPool()
	if err != nil {
		log.Fatalf("Failed to create database pool: %v", err)
	}
	defer pool.Close()

	// FGA setup
	fgaClient, err := openfgaClient.NewSdkClient(&openfgaClient.ClientConfiguration{
		ApiUrl:  os.Getenv("FGA_API_URL"),
		StoreId: os.Getenv("FGA_STORE_ID"),
	})

	if err != nil {
		log.Fatalf("Failed to create FGA client: %v", err)
	}

	// Run migrations
	ctx := context.Background()
	if err := database.RunMigrations(ctx, pool); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Logger
	textHandler := slog.NewTextHandler(os.Stdout, nil)
	logger := slog.New(textHandler)

	// Dependency Injection
	// Use PostgresStore for database persistence
	var store user.Store = user.NewPostgresStore(pool)
	// For in-memory storage (old implementation), use:
	// var store user.Store = user.NewInMemStore()

	var service user.Service = user.NewInMemoryUserService(store)
	var handler = user.Handler{Service: service, Logger: logger}

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Health check route
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	// User routes
	r.Post("/user", ErrorHandler(handler.CreateUser))
	r.Post("/authenticate", ErrorHandler(handler.Authenticate))
	r.Get("/user", ErrorHandler(handler.SearchUser))
	r.Get("/user/{id}", ErrorHandler(handler.GetUser))

	// Authentication routes
	r.Post("/authenticate", ErrorHandler(handler.Authenticate))

	log.Println("Server starting on :4001")
	err = http.ListenAndServe(":4000", r)
	if err != nil {
		panic(err)
	}
}

func createDBPool() (*pgxpool.Pool, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Build connection string
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)

	// Configure pool
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Create pool
	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}
