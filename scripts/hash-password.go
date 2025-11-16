package main

import (
	"flag"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := flag.String("password", "", "Password to hash")
	cost := flag.Int("cost", 12, "Bcrypt cost (10-31)")
	flag.Parse()

	if *password == "" {
		fmt.Println("❌ Password is required")
		fmt.Println("Usage: go run hash-password.go -password 'YOUR_PASSWORD' [-cost 12]")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(*password), *cost)
	if err != nil {
		fmt.Printf("❌ Error generating hash: %v\n", err)
		return
	}

	fmt.Printf("✅ Password hash:\n%s\n", hash)
}
