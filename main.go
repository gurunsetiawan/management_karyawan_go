//go:build ignore
// +build ignore

// This file is kept for backward compatibility.
// The main entry point is now cmd/server/main.go
//
// To run the application, use one of:
//   - go run cmd/server/main.go
//   - make run
//   - go build -o karyawan-app cmd/server/main.go && ./karyawan-app
//
// This file is marked with 'ignore' build tag to prevent compilation conflicts.
package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	fmt.Println("⚠️  DEPRECATED: This entry point is deprecated.")
	fmt.Println("   Please use: go run cmd/server/main.go")
	fmt.Println()
	fmt.Println("Redirecting to cmd/server/main.go...")
	fmt.Println()

	// Get the directory of the current executable
	cmd := exec.Command("go", "run", "cmd/server/main.go", os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
