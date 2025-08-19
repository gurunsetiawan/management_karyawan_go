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

var DB *sql.DB

func ConnectDB() error {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// First connect without database to create it if needed
	rootDB, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/", user, password, host, port))
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL: %w", err)
	}
	defer rootDB.Close()

	// Create database if not exists
	_, err = rootDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	// Now connect to the specific database
	DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		user, password, host, port, dbName))
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(5 * time.Minute)

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Run database migrations
	if err := runMigrations(); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Successfully connected to MySQL!")
	return nil
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

	log.Println("Database migrations completed successfully!")
	return nil
}
