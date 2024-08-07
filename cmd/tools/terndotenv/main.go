package main

import (
	"log/slog"
	"os/exec"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file", "error", err)
		return
	}

	slog.Info("Starting migration...")

	cmd := exec.Command("tern", "migrate", "--migrations", "./internal/store/pgstore/migrations", "--config", "./internal/store/pgstore/migrations/tern.conf")

	if err := cmd.Run(); err != nil {
		slog.Error("Command failed with error", "error", err)
		return
	}

	slog.Info("Migration completed successfully")
}
