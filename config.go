package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvValues(logger *slog.Logger) (string, string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		logger.Error("Error loading .env file", "Error", err)
		return "", "", err
	}

	token := os.Getenv("BOT_TOKEN")
	postgresPW := os.Getenv("POSTGRES_PW")

	return token, postgresPW, nil
}
