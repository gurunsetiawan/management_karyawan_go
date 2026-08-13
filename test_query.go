package main

import (
	"database/sql"
	"fmt"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	LastLogin *time.Time
}

func main() {
	db, _ := sql.Open("mysql", "karyawan_app:karyawan_app123@tcp(localhost:3306)/karyawan_db?parseTime=true")
	defer db.Close()

	user := &User{}
	err := db.QueryRow("SELECT last_login FROM users WHERE username = 'admin'").Scan(&user.LastLogin)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Success:", user.LastLogin)
	}
}
