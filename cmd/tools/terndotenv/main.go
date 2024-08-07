package main

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	fmt.Println("Starting migration...")

	cmd := exec.Command("tern", "migrate", "--migrations", "./internal/store/pgstore/migrations", "--config", "./internal/store/pgstore/migrations/tern.conf")

	if err := cmd.Run(); err != nil {
		fmt.Printf("Command failed with error: %v\n", err)
		return
	}

	fmt.Println("Migration completed successfully")
}
