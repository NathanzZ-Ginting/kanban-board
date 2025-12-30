package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	Host          string
	Environment   string
	MySQLHost     string
	MySQLPort     string
	MySQLUser     string
	MySQLPassword string
	MySQLDatabase string
	JWTSecret     string
	JWTExpiry     time.Duration
	RefreshExpiry time.Duration
}

func LoadConfig() *Config {
	godotenv.Load()

	jwtExpiry, _ := time.ParseDuration(getEnv("JWT_EXPIRY", "24h"))
	refreshExpiry, _ := time.ParseDuration(getEnv("REFRESH_TOKEN_EXPIRY", "168h"))

	return &Config{
		Port:          getEnv("AUTH_SERVICE_PORT", "8001"),
		Host:          getEnv("AUTH_SERVICE_HOST", "0.0.0.0"),
		Environment:   getEnv("GO_ENV", "development"),
		MySQLHost:     getEnv("MYSQL_HOST", "localhost"),
		MySQLPort:     getEnv("MYSQL_PORT", "3306"),
		MySQLUser:     getEnv("MYSQL_USER", "root"),
		MySQLPassword: getEnv("MYSQL_PASSWORD", ""),
		MySQLDatabase: getEnv("MYSQL_DATABASE", "kanban_db"),
		JWTSecret:     getEnv("JWT_SECRET", "secret"),
		JWTExpiry:     jwtExpiry,
		RefreshExpiry: refreshExpiry,
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
