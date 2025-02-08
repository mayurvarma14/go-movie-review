package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Warning: Error loading .env file (this is OK in Container):", err)
	}

	required := []string{
		"MONGO_APP_USER",
		"MONGO_APP_PASSWORD",
		"MONGO_DOMAIN",
		"MONGO_INITDB_DATABASE",
		"MONGO_AUTH_SOURCE",
		"SECRET_KEY",
	}

	var missing []string
	for _, key := range required {
		if os.Getenv(key) == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) > 0 {
		log.Fatalf("Missing required environment variables: %v", missing)
	}
}
