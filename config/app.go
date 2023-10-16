package config

import (
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppPort string
}

func GetAppConfig(isUsingDotEnv bool) *AppConfig {
	if isUsingDotEnv {
		godotenv.Load()
	}

	return &AppConfig{
		AppPort: os.Getenv("APP_PORT"),
	}
}
