package testutils

import (
	"os"

	"github.com/P3rCh1/immersive-images/backend/internal/config"
)

func InitTestConfig() {
	config.Config.DB.Host = env("POSTGRES_HOST", "postgres")
	config.Config.DB.Port = env("POSTGRES_PORT", "5432")
	config.Config.DB.Name = env("POSTGRES_DB", "images")
	config.Config.DB.Password = env("POSTGRES_PASSWORD", "postgres")
}

func env(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}

	return fallback
}
