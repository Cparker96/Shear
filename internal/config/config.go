package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvValues(logger *slog.Logger) (string, string, string, string, error) {
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file", "Error", err)
		return "", "", "", "", err
	}

	token := os.Getenv("BOT_TOKEN")
	postgresPW := os.Getenv("POSTGRES_PW")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresDBName := os.Getenv("POSTGRES_DBNAME")

	return token, postgresPW, postgresUser, postgresDBName, nil
}
