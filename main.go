// main.go
package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"karyawan-app/config"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/time/rate"
)

type Employee struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Phone     string `json:"phone"`
	Alamat    string `json:"alamat"`
	CreatedAt string `json:"created_at"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Rate limiter: configurable requests per minute per IP
var limiter *rate.Limiter

func initRateLimiter() {
	requestsPerMinute := 100 // Default
	if envRateLimit := os.Getenv("RATE_LIMIT_REQUESTS"); envRateLimit != "" {
		if rate, err := strconv.Atoi(envRateLimit); err == nil && rate > 0 {
			requestsPerMinute = rate
		}
	}

	windowSeconds := 60 // Default 1 minute
	if envWindow := os.Getenv("RATE_LIMIT_WINDOW"); envWindow != "" {
		if window, err := strconv.Atoi(envWindow); err == nil && window > 0 {
			windowSeconds = window
		}
	}

	limiter = rate.NewLimiter(rate.Every(time.Duration(windowSeconds)*time.Second/time.Duration(requestsPerMinute)), requestsPerMinute)
}

// Email validation regex
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// Input sanitization function
func sanitizeInput(input string) string {
	// Remove HTML tags and dangerous characters
	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")
	input = strings.ReplaceAll(input, "\"", "&quot;")
	input = strings.ReplaceAll(input, "'", "&#39;")
	return input
}

// Email validation function
func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// CORS middleware
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		corsOrigin := os.Getenv("CORS_ORIGIN")
		if corsOrigin == "" {
			corsOrigin = "*" // Default to allow all origins
		}

		w.Header().Set("Access-Control-Allow-Origin", corsOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// Rate limiting middleware
func rateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if limiter != nil && !limiter.Allow() {
			errorResponse := ErrorResponse{
				Error:   "Rate limit exceeded",
				Message: "Too many requests. Please try again later.",
				Code:    http.StatusTooManyRequests,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
		next(w, r)
	}
}

// Error logging middleware
func errorLoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	}
}

func main() {
	// Initialize rate limiter
	initRateLimiter()

	// Initialize database connection
	if err := config.ConnectDB(); err != nil {
		log.Printf("Warning: Database connection failed: %v", err)
		log.Println("Starting server without database...")
		// Continue without database for demo purposes
	}

	// Serve static files with proper MIME types
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("frontend/css"))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.Handle("/", http.FileServer(http.Dir("frontend")))

	// API Routes with middleware
	http.HandleFunc("/api/employees", corsMiddleware(rateLimitMiddleware(errorLoggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			getEmployees(w, r)
		case "POST":
			createEmployee(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))))

	http.HandleFunc("/api/employees/", corsMiddleware(rateLimitMiddleware(errorLoggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Path[len("/api/employees/"):])
		if err != nil || id <= 0 {
			errorResponse := ErrorResponse{
				Error:   "Invalid ID",
				Message: "Employee ID must be a positive integer",
				Code:    http.StatusBadRequest,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}

		switch r.Method {
		case "GET":
			getEmployee(w, r, id)
		case "PUT":
			updateEmployee(w, r, id)
		case "DELETE":
			deleteEmployee(w, r, id)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))))

	// Get server configuration from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083" // Default port
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "127.0.0.1" // Default host
	}

	serverAddr := host + ":" + port
	log.Printf("Server berjalan di http://%s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, nil))
}

func getEmployees(w http.ResponseWriter, r *http.Request) {
	if config.DB == nil {
		errorResponse := ErrorResponse{
			Error:   "Database not available",
			Message: "Database connection is not available",
			Code:    http.StatusServiceUnavailable,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	rows, err := config.DB.Query("SELECT id, name, email, role, phone, alamat, created_at FROM employees WHERE deleted_at IS NULL")
	if err != nil {
		log.Printf("Database query error: %v", err)
		errorResponse := ErrorResponse{
			Error:   "Database error",
			Message: "Failed to fetch employees",
			Code:    http.StatusInternalServerError,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	defer rows.Close()

	var employees []Employee
	for rows.Next() {
		var e Employee
		if err := rows.Scan(&e.ID, &e.Name, &e.Email, &e.Role, &e.Phone, &e.Alamat, &e.CreatedAt); err != nil {
			log.Printf("Row scan error: %v", err)
			errorResponse := ErrorResponse{
				Error:   "Data processing error",
				Message: "Failed to process employee data",
				Code:    http.StatusInternalServerError,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
		employees = append(employees, e)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		errorResponse := ErrorResponse{
			Error:   "Data iteration error",
			Message: "Failed to iterate through employee data",
			Code:    http.StatusInternalServerError,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employees)
}

func getEmployee(w http.ResponseWriter, r *http.Request, id int) {
	row := config.DB.QueryRow("SELECT id, name, email, role, phone, alamat, created_at FROM employees WHERE id = ? AND deleted_at IS NULL", id)

	var e Employee
	err := row.Scan(&e.ID, &e.Name, &e.Email, &e.Role, &e.Phone, &e.Alamat, &e.CreatedAt)
	if err == sql.ErrNoRows {
		errorResponse := ErrorResponse{
			Error:   "Employee not found",
			Message: "Karyawan tidak ditemukan",
			Code:    http.StatusNotFound,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(e)
}

func createEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		errorResponse := ErrorResponse{
			Error:   "Invalid content type",
			Message: "Content-Type must be application/json",
			Code:    http.StatusUnsupportedMediaType,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	var emp Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		errorResponse := ErrorResponse{
			Error:   "Invalid JSON format",
			Message: "Request body must be valid JSON",
			Code:    http.StatusBadRequest,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Sanitize and validate input
	emp.Name = sanitizeInput(emp.Name)
	emp.Email = sanitizeInput(emp.Email)
	emp.Role = sanitizeInput(emp.Role)
	emp.Phone = sanitizeInput(emp.Phone)
	emp.Alamat = sanitizeInput(emp.Alamat)

	// Enhanced validation
	if emp.Name == "" || len(emp.Name) < 2 {
		errorResponse := ErrorResponse{
			Error:   "Invalid name",
			Message: "Name must be at least 2 characters long",
			Code:    http.StatusBadRequest,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if !isValidEmail(emp.Email) {
		errorResponse := ErrorResponse{
			Error:   "Invalid email",
			Message: "Please provide a valid email address",
			Code:    http.StatusBadRequest,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if emp.Alamat == "" || len(emp.Alamat) < 10 {
		errorResponse := ErrorResponse{
			Error:   "Invalid address",
			Message: "Address must be at least 10 characters long",
			Code:    http.StatusBadRequest,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	res, err := config.DB.Exec(
		"INSERT INTO employees (name, email, role, phone, alamat) VALUES (?, ?, ?, ?, ?)",
		emp.Name, emp.Email, emp.Role, emp.Phone, emp.Alamat,
	)
	if err != nil {
		log.Printf("Insert error: %v", err)
		errorResponse := ErrorResponse{
			Error:   "Database error",
			Message: "Failed to create employee",
			Code:    http.StatusInternalServerError,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("LastInsertId error: %v", err)
		errorResponse := ErrorResponse{
			Error:   "Database error",
			Message: "Failed to get created employee ID",
			Code:    http.StatusInternalServerError,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	emp.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(emp)
}

func updateEmployee(w http.ResponseWriter, r *http.Request, id int) {
	if r.Header.Get("Content-Type") != "application/json" {
		errorResponse := ErrorResponse{
			Error:   "Invalid content type",
			Message: "Content-Type must be application/json",
			Code:    http.StatusUnsupportedMediaType,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	var emp Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		errorResponse := ErrorResponse{
			Error:   "Invalid JSON format",
			Message: "Request body must be valid JSON",
			Code:    http.StatusBadRequest,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Sanitize and validate input
	emp.Name = sanitizeInput(emp.Name)
	emp.Email = sanitizeInput(emp.Email)
	emp.Role = sanitizeInput(emp.Role)
	emp.Phone = sanitizeInput(emp.Phone)
	emp.Alamat = sanitizeInput(emp.Alamat)

	// Enhanced validation
	if emp.Name == "" || len(emp.Name) < 2 {
		errorResponse := ErrorResponse{
			Error:   "Invalid name",
			Message: "Name must be at least 2 characters long",
			Code:    http.StatusBadRequest,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	if !isValidEmail(emp.Email) {
		errorResponse := ErrorResponse{
			Error:   "Invalid email",
			Message: "Please provide a valid email address",
			Code:    http.StatusBadRequest,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	_, err := config.DB.Exec(
		"UPDATE employees SET name=?, email=?, role=?, phone=?, alamat=? WHERE id=? AND deleted_at IS NULL",
		emp.Name, emp.Email, emp.Role, emp.Phone, emp.Alamat, id,
	)
	if err != nil {
		log.Printf("Update error: %v", err)
		errorResponse := ErrorResponse{
			Error:   "Database error",
			Message: "Failed to update employee",
			Code:    http.StatusInternalServerError,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	emp.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(emp)
}

func deleteEmployee(w http.ResponseWriter, r *http.Request, id int) {
	// Soft delete implementation
	_, err := config.DB.Exec("UPDATE employees SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL", id)
	if err != nil {
		log.Printf("Delete error: %v", err)
		errorResponse := ErrorResponse{
			Error:   "Database error",
			Message: "Failed to delete employee",
			Code:    http.StatusInternalServerError,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
