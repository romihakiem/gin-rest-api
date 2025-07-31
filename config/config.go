package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	APPPort          string
	DBHost           string
	DBPort           string
	DBName           string
	DBUser           string
	DBPassword       string
	DBSSL            string
	RedisHost        string
	RedisPort        string
	RedisPassword    string
	RedisDB          int
	JWTSecret        string
	JWTAccessExpiry  int
	JWTRefreshExpiry int
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	return &Config{
		APPPort:          getEnv("APP_PORT", "8080"),
		DBHost:           getEnv("DB_HOST", "127.0.0.1"),
		DBPort:           getEnv("DB_PORT", "5432"),
		DBName:           getEnv("DB_NAME", "postgres"),
		DBUser:           getEnv("DB_USER", "postgres"),
		DBPassword:       getEnv("DB_PASSWORD", "postgres"),
		DBSSL:            getEnv("DB_SSL", "disable"),
		RedisHost:        getEnv("REDIS_HOST", "127.0.0.1"),
		RedisPort:        getEnv("REDIS_PORT", "6379"),
		RedisPassword:    getEnv("REDIS_PASSWORD", ""),
		RedisDB:          getEnvAsInt("REDIS_DB", 0),
		JWTSecret:        getEnv("JWT_SECRET", ""),
		JWTAccessExpiry:  getEnvAsInt("JWT_ACCESS_EXPIRY", 3600),
		JWTRefreshExpiry: getEnvAsInt("JWT_REFRESH_EXPIRY", 604800),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
