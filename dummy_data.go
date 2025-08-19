package main

import (
	"fmt"
	"log"

	"karyawan-app/config"

	"github.com/brianvoe/gofakeit/v6"
)

func GenerateDummyData() {
	gofakeit.Seed(0)

	for i := 0; i < 20; i++ {
		_, err := config.DB.Exec(
			"INSERT INTO employees (name, email, role, phone, alamat) VALUES (?, ?, ?, ?, ?)",
			gofakeit.Name(),
			gofakeit.Email(),
			gofakeit.JobTitle(),
			gofakeit.Phone(),
			gofakeit.Address().Address,
		)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("Berhasil menambahkan 20 data dummy")
}

// RunDummyDataGenerator can be called to generate dummy data
func RunDummyDataGenerator() {
	config.ConnectDB()
	GenerateDummyData()
}
