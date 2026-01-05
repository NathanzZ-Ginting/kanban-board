package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Host        string
	Environment string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	JWTSecret   string
	JWTExpiry   time.Duration
}

func LoadConfig() *Config {
	// Try to load .env from current dir, then parent dirs
	godotenv.Load()
	godotenv.Load("../../.env")
	godotenv.Load(filepath.Join("..", "..", ".env"))

	jwtExpiry, _ := time.ParseDuration(getEnv("JWT_EXPIRY", "24h"))

	return &Config{
		Port:        getEnv("TASK_SERVICE_PORT", "8004"),
		Host:        getEnv("TASK_SERVICE_HOST", "0.0.0.0"),
		Environment: getEnv("GO_ENV", "development"),
		DBHost:      getEnv("SUPABASE_DB_HOST", ""),
		DBPort:      getEnv("SUPABASE_DB_PORT", "5432"),
		DBUser:      getEnv("SUPABASE_DB_USER", "postgres"),
		DBPassword:  getEnv("SUPABASE_DB_PASSWORD", ""),
		DBName:      getEnv("SUPABASE_DB_NAME", "postgres"),
		JWTSecret:   getEnv("JWT_SECRET", "secret"),
		JWTExpiry:   jwtExpiry,
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
