package config

import (
	"os"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	DBDriver   string
	DBHost     string
	DBName     string
	DBUsername string
	DBPassword string
	DBPort     string
}

func GetDBConfig(isUsingDotEnv bool) DBConfig {
	if isUsingDotEnv {
		godotenv.Load()
	}

	return DBConfig{
		DBDriver:   os.Getenv("DB_DRIVER"),
		DBHost:     os.Getenv("DB_HOST"),
		DBName:     os.Getenv("DB_NAME"),
		DBUsername: os.Getenv("DB_USERNAME"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBPort:     os.Getenv("DB_PORT"),
	}
}
