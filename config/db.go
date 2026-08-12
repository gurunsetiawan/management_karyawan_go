// config/db.go
package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

const (
	defaultDBHost     = "127.0.0.1"
	defaultDBPort     = "3306"
	defaultDBUser     = "root"
	defaultDBPassword = "rootpassword"
	defaultDBName     = "karyawan_app"
)

var DB *sql.DB

func ConnectDB() error {
	// Load environment variables from .env file in current directory
	// Try multiple possible paths
	envPaths := []string{".env", "./.env", "../.env"}
	for _, path := range envPaths {
		if err := godotenv.Load(path); err == nil {
			log.Printf("Loaded .env from: %s", path)
			break
		}
	}

	host := getEnv("DB_HOST", defaultDBHost)
	port := getEnv("DB_PORT", defaultDBPort)
	user := getEnv("DB_USER", defaultDBUser)
	password := getEnv("DB_PASSWORD", defaultDBPassword)
	dbName := getEnv("DB_NAME", defaultDBName)

	log.Printf("Connecting to database: %s@%s:%s/%s", user, host, port, dbName)

	// Connect to the database with retry logic
	var db *sql.DB
	var err error
	maxRetries := 3
	retryDelay := 2 * time.Second

	// Connection string with parameters for MariaDB
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local&timeout=5s&readTimeout=30s&writeTimeout=30s&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		user, password, host, port, dbName)

	// Try to connect with retries
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("mysql", dsn)
		if err == nil {
			err = db.Ping()
			if err == nil {
				DB = db
				break
			}
			db.Close()
		}
		
		if i < maxRetries-1 {
			log.Printf("Connection attempt %d failed: %v. Retrying in %v...", i+1, err, retryDelay)
			time.Sleep(retryDelay)
		} else {
			return fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
		}
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)
	DB.SetConnMaxIdleTime(5 * time.Minute)

	// Run database migrations
	if err := runMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Printf("Successfully connected to MariaDB at %s:%s/%s", host, port, dbName)
	return nil
}

// getEnv returns the value of the environment variable or the default value if not set
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func runMigrations() error {
	// Migration 1: Create employees table with soft delete
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS employees (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		role VARCHAR(100) NOT NULL,
		phone VARCHAR(20) UNIQUE NOT NULL,
		alamat TEXT NOT NULL,
		position VARCHAR(255) DEFAULT '',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP NULL DEFAULT NULL
	)`)
	if err != nil {
		return fmt.Errorf("failed to create employees table: %w", err)
	}

	// Migration 1.1: Add deleted_at column if it doesn't exist
	_, err = DB.Exec(`ALTER TABLE employees ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP NULL DEFAULT NULL`)
	if err != nil {
		log.Printf("Warning: Could not add deleted_at column: %v", err)
	}

	// Migration 1.2: Add updated_at column if it doesn't exist
	_, err = DB.Exec(`ALTER TABLE employees ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP`)
	if err != nil {
		log.Printf("Warning: Could not add updated_at column: %v", err)
	}

	// Migration 1.3: Add position column if it doesn't exist
	_, err = DB.Exec(`ALTER TABLE employees ADD COLUMN IF NOT EXISTS position VARCHAR(255) DEFAULT ''`)
	if err != nil {
		log.Printf("Warning: Could not add position column: %v", err)
	}

	// Migration 2: Add indexes for better performance
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_employees_email ON employees(email)",
		"CREATE INDEX IF NOT EXISTS idx_employees_created_at ON employees(created_at)",
		"CREATE INDEX IF NOT EXISTS idx_employees_role ON employees(role)",
	}

	for _, index := range indexes {
		_, err := DB.Exec(index)
		if err != nil {
			log.Printf("Warning: Failed to create index: %v", err)
		}
	}

	// Add deleted_at index only if column exists
	_, err = DB.Exec("CREATE INDEX IF NOT EXISTS idx_employees_deleted_at ON employees(deleted_at)")
	if err != nil {
		log.Printf("Warning: Could not create deleted_at index: %v", err)
	}

	// Migration 3: Create migration tracking table
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS migrations (
		id INT AUTO_INCREMENT PRIMARY KEY,
		version VARCHAR(50) NOT NULL UNIQUE,
		description TEXT,
		applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Migration 4: Insert migration records
	migrations := []struct {
		version     string
		description string
	}{
		{"001", "Create employees table with soft delete"},
		{"002", "Add database indexes"},
		{"003", "Create migration tracking table"},
	}

	for _, migration := range migrations {
		_, err := DB.Exec("INSERT IGNORE INTO migrations (version, description) VALUES (?, ?)",
			migration.version, migration.description)
		if err != nil {
			log.Printf("Warning: Failed to insert migration record: %v", err)
		}
	}

	// Migration 5: Create roles table
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS roles (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(50) NOT NULL UNIQUE,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return fmt.Errorf("failed to create roles table: %w", err)
	}

	// Migration 6: Create users table
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(50) NOT NULL UNIQUE,
		email VARCHAR(255) NOT NULL UNIQUE,
		password_hash VARCHAR(255) NOT NULL,
		role_id INT NOT NULL,
		is_active BOOLEAN DEFAULT TRUE,
		last_login TIMESTAMP NULL DEFAULT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE RESTRICT
	)`)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Migration 7: Insert default roles
	roles := []struct {
		name        string
		description string
	}{
		{"admin", "Administrator with full access"},
		{"manager", "Manager with limited access"},
		{"user", "Regular user with read-only access"},
	}

	for _, role := range roles {
		_, err := DB.Exec("INSERT IGNORE INTO roles (name, description) VALUES (?, ?)",
			role.name, role.description)
		if err != nil {
			log.Printf("Warning: Failed to insert role: %v", err)
		}
	}

	// Migration 8: Create default admin user placeholder
	// Password will be set by seed command: go run cmd/seed/seed.go
	// Default credentials: admin / admin123
	_, err = DB.Exec(`
		INSERT INTO users (username, email, password_hash, role_id, is_active)
		SELECT 'admin', 'admin@example.com', 'pending_seed', 1, TRUE
		WHERE NOT EXISTS (SELECT 1 FROM users WHERE username = 'admin')
	`)
	if err != nil {
		log.Printf("Warning: Failed to create default admin user placeholder: %v", err)
	}

	// Migration 9: Insert migration records for auth
	authMigrations := []struct {
		version     string
		description string
	}{
		{"004", "Create roles table"},
		{"005", "Create users table"},
		{"006", "Insert default roles"},
		{"007", "Create default admin user"},
	}

	for _, migration := range authMigrations {
		_, err := DB.Exec("INSERT IGNORE INTO migrations (version, description) VALUES (?, ?)",
			migration.version, migration.description)
		if err != nil {
			log.Printf("Warning: Failed to insert migration record: %v", err)
		}
	}

	log.Println("Database migrations completed successfully!")
	return nil
}
