package testutils

import (
	"os"

	"github.com/P3rCh1/immersive-images/backend/internal/config"
)

// InitTestConfig replaces the config with the values required by acceptance
// tests, so they do not depend on any config file. Every value can be
// overridden with an env variable to point tests at another instance.
func InitTestConfig() {
	config.Config.DB.Host = envOrDefault("POSTGRES_HOST", "postgres")
	config.Config.DB.Port = envOrDefault("POSTGRES_PORT", "5432")
	config.Config.DB.Name = envOrDefault("POSTGRES_DB", "images")
	config.Config.DB.Password = envOrDefault("POSTGRES_PASSWORD", "postgres")
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}

	return fallback
}
