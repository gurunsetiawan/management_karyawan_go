package main

import (
	"fmt"
	"log"

	"karyawan-app/config"
	"karyawan-app/internal/auth"
)

func main() {
	// Connect to database
	if err := config.ConnectDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer config.DB.Close()

	// Generate password hash for admin123
	password := "admin123"
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	fmt.Printf("Generated hash for 'admin123': %s\n\n", hashedPassword)

	// Update admin user with correct password hash
	_, err = config.DB.Exec(`
		UPDATE users 
		SET password_hash = ?, email = 'admin@example.com'
		WHERE username = 'admin'
	`, hashedPassword)
	if err != nil {
		log.Printf("Warning: Failed to update admin password: %v", err)
	} else {
		fmt.Println("✓ Admin user password updated successfully!")
	}

	// Or create if not exists
	_, err = config.DB.Exec(`
		INSERT INTO users (username, email, password_hash, role_id, is_active)
		SELECT 'admin', 'admin@example.com', ?, 1, TRUE
		WHERE NOT EXISTS (SELECT 1 FROM users WHERE username = 'admin')
	`, hashedPassword)
	if err != nil {
		log.Printf("Warning: Failed to create admin user: %v", err)
	} else {
		fmt.Println("✓ Admin user created successfully!")
	}

	fmt.Println("\n===========================================")
	fmt.Println("Default Admin User:")
	fmt.Println("  Username: admin")
	fmt.Println("  Password: admin123")
	fmt.Println("===========================================")
}
