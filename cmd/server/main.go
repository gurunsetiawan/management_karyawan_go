package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"karyawan-app/config"
	handler "karyawan-app/internal/handler"
	middleware "karyawan-app/internal/middleware"
	repo "karyawan-app/internal/repository"
	service "karyawan-app/internal/service"
)

// spaFileSystem wraps http.FileSystem to support SPA routing
type spaFileSystem struct {
	root http.FileSystem
}

func (fs *spaFileSystem) Open(name string) (http.File, error) {
	f, err := fs.root.Open(name)
	if os.IsNotExist(err) {
		return fs.root.Open("index.html")
	}
	return f, err
}

func main() {
	// Initialize database connection using shared config
	if err := config.ConnectDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if config.DB != nil {
			if err := config.DB.Close(); err != nil {
				log.Printf("Warning: Failed to close database connection: %v", err)
			}
		}
	}()

	// Initialize repositories
	employeeRepo := repo.NewEmployeeRepository(config.DB)
	userRepo := repo.NewUserRepository(config.DB)

	// Initialize services
	employeeService := service.NewEmployeeService(employeeRepo)
	authService := service.NewAuthService(userRepo)

	// Initialize handlers
	employeeHandler := handler.NewEmployeeHandler(employeeService)
	authHandler := handler.NewAuthHandler(authService)

	// Create router
	r := mux.NewRouter()

	// Apply global middleware to all routes
	r.Use(mux.MiddlewareFunc(handler.CORSMiddleware))
	r.Use(mux.MiddlewareFunc(handler.RateLimitMiddleware(getRateLimitConfig())))
	r.Use(mux.MiddlewareFunc(handler.LoggingMiddleware))

	// Public routes (no authentication required)
	publicRouter := r.PathPrefix("/api").Subrouter()
	publicRouter.Use(mux.MiddlewareFunc(handler.JSONContentTypeMiddleware))
	
	// Auth routes - login is public
	publicRouter.HandleFunc("/auth/login", authHandler.Login).Methods("POST")
	
	// Health check - public
	publicRouter.HandleFunc("/health", healthCheckHandler).Methods("GET")

	// Protected routes (require authentication)
	protectedRouter := r.PathPrefix("/api").Subrouter()
	protectedRouter.Use(mux.MiddlewareFunc(handler.JSONContentTypeMiddleware))
	protectedRouter.Use(mux.MiddlewareFunc(middleware.AuthMiddleware))
	
	// Auth routes - profile requires auth
	protectedRouter.HandleFunc("/auth/profile", authHandler.GetProfile).Methods("GET")
	
	// Employee routes - protected by default
	// Note: You can remove this line to make employee routes public
	// Currently employees API requires authentication
	protectedRouter.HandleFunc("/employees", employeeHandler.GetAllEmployees).Methods("GET")
	protectedRouter.HandleFunc("/employees/export", employeeHandler.ExportEmployees).Methods("GET")
	protectedRouter.HandleFunc("/employees/import", employeeHandler.ImportEmployees).Methods("POST")
	protectedRouter.HandleFunc("/employees/{id}", employeeHandler.GetEmployee).Methods("GET")
	protectedRouter.HandleFunc("/employees", employeeHandler.CreateEmployee).Methods("POST")
	protectedRouter.HandleFunc("/employees/{id}", employeeHandler.UpdateEmployee).Methods("PUT")
	protectedRouter.HandleFunc("/employees/{id}", employeeHandler.DeleteEmployee).Methods("DELETE")

	// Admin routes (require admin role)
	adminRouter := r.PathPrefix("/api/admin").Subrouter()
	adminRouter.Use(mux.MiddlewareFunc(handler.JSONContentTypeMiddleware))
	adminRouter.Use(mux.MiddlewareFunc(middleware.AuthMiddleware))
	adminRouter.Use(mux.MiddlewareFunc(middleware.RoleMiddleware("admin")))
	
	// Example admin route
	adminRouter.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"message":"Admin area accessed successfully"}`)); err != nil {
			log.Printf("Error writing admin response: %v", err)
		}
	}).Methods("GET")


	// Serve static files from the frontend/build directory (SPA fallback support)
	frontendDir := "./frontend/build"
	if _, err := os.Stat(frontendDir); !os.IsNotExist(err) {
		r.PathPrefix("/").Handler(http.FileServer(&spaFileSystem{http.Dir(frontendDir)}))
	} else {
		// Fallback to frontend directory for development
		frontendDevDir := "./frontend"
		if _, err := os.Stat(frontendDevDir); !os.IsNotExist(err) {
			r.PathPrefix("/").Handler(http.FileServer(&spaFileSystem{http.Dir(frontendDevDir)}))
		}
	}

	// Server configuration
	port := getEnv("PORT", "8080")
	host := getEnv("HOST", "127.0.0.1")
	serverAddr := host + ":" + port

	server := &http.Server{
		Addr:         serverAddr,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on http://%s", serverAddr)
		log.Println("")
		log.Println("===========================================")
		log.Println("API Endpoints:")
		log.Println("  POST /api/auth/login    - Login")
		log.Println("  GET  /api/auth/profile  - Get Profile (protected)")
		log.Println("  GET  /api/employees     - Get Employees (protected)")
		log.Println("  GET  /api/health        - Health Check")
		log.Println("===========================================")
		log.Println("")
		
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}

// healthCheckHandler returns a simple health check response
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := `{"status":"healthy","timestamp":"` + time.Now().Format(time.RFC3339) + `"}`
	if _, err := w.Write([]byte(response)); err != nil {
		log.Printf("Error writing health check response: %v", err)
	}
}

// getRateLimitConfig returns the rate limit configuration from environment
func getRateLimitConfig() int {
	if rateLimit := os.Getenv("RATE_LIMIT_REQUESTS"); rateLimit != "" {
		if rate, err := strconv.Atoi(rateLimit); err == nil && rate > 0 {
			return rate
		}
	}
	return 100 // Default: 100 requests per minute
}

// getEnv returns the value of the environment variable or the default value if not set
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
