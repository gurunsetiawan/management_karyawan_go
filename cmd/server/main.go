package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"

	handler "karyawan-app/internal/handler"
	repo "karyawan-app/internal/repository"
	service "karyawan-app/internal/service"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}

	// Initialize database connection
	db := initDB()
	defer db.Close()

	// Initialize repository, service, and handler
	employeeRepo := repo.NewEmployeeRepository(db)
	employeeService := service.NewEmployeeService(employeeRepo)
	employeeHandler := handler.NewEmployeeHandler(employeeService)

	// Create router
	r := mux.NewRouter()

	// Apply middleware
	middleware := handler.NewChain(
		handler.CORSMiddleware,
		handler.RateLimitMiddleware(100), // 100 requests per minute
		handler.LoggingMiddleware,
		handler.JSONContentTypeMiddleware,
	)

	// Register routes
	api := r.PathPrefix("/api").Subrouter()
	employeeHandler.RegisterRoutes(api)

	// Serve static files from the frontend directory
	frontendDir := "./frontend"
	if _, err := os.Stat(frontendDir); !os.IsNotExist(err) {
		r.PathPrefix("/").Handler(http.FileServer(http.Dir(frontendDir)))
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      middleware.Then(r),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Server starting on port %s\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}

func initDB() *sql.DB {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Default values if not set
	if dbUser == "" {
		dbUser = "root"
	}
	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "3306"
	}
	if dbName == "" {
		dbName = "karyawan_db"
	}

	dsn := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?parseTime=true"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Run migrations
	runMigrations(db)

	return db
}

func runMigrations(db *sql.DB) {
	// Create employees table if not exists
	query := `
	CREATE TABLE IF NOT EXISTS employees (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) NOT NULL UNIQUE,
		role VARCHAR(50) NOT NULL,
		phone VARCHAR(20) NOT NULL,
		alamat TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`

	if _, err := db.Exec(query); err != nil {
		log.Fatalf("Error creating employees table: %v", err)
	}
}
