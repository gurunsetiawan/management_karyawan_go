package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	db, err := sql.Open("mysql", "karyawan_app:karyawan_app123@tcp(localhost:3306)/karyawan_db")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return
	}
	defer db.Close()

	hash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	_, err = db.Exec("UPDATE users SET password_hash = ? WHERE username = 'admin'", string(hash))
	if err != nil {
		fmt.Println("Error updating:", err)
	} else {
		fmt.Println("Password updated successfully via Go")
	}
}
