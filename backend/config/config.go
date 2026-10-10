package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Redis
	RedisURL      string
	RedisPassword string

	// JWT
	JWTSecret            string
	JWTExpiry            time.Duration
	RefreshTokenExpiry   time.Duration

	// Bootstrap - Single Operator
	OperatorUsername string
	OperatorPassword string

	// API
	APIPort string
	APIHost string
	Environment string

	// CORS
	FrontendURL string

	// Logging
	LogLevel string
}

func LoadConfig() *Config {
	jwtExpiry, _ := time.ParseDuration(getEnv("JWT_EXPIRY", "15m"))
	refreshExpiry, _ := time.ParseDuration(getEnv("REFRESH_TOKEN_EXPIRY", "720h"))

	return &Config{
		DBHost:               getEnv("DB_HOST", "localhost"),
		DBPort:               getEnv("DB_PORT", "5432"),
		DBUser:               getEnv("DB_USER", "postgres"),
		DBPassword:           getEnv("DB_PASSWORD", "postgres123"),
		DBName:               getEnv("DB_NAME", "toko_berkah"),
		DBSSLMode:            getEnv("DB_SSLMODE", "disable"),
		RedisURL:             getEnv("REDIS_URL", "localhost:6379"),
		RedisPassword:        getEnv("REDIS_PASSWORD", ""),
		JWTSecret:            getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-this-in-production"),
		JWTExpiry:            jwtExpiry,
		RefreshTokenExpiry:   refreshExpiry,
		OperatorUsername:     getEnv("OPERATOR_USERNAME", ""),
		OperatorPassword:     getEnv("OPERATOR_PASSWORD", ""),
		APIPort:              getEnv("API_PORT", "8080"),
		APIHost:              getEnv("API_HOST", "0.0.0.0"),
		Environment:          getEnv("ENVIRONMENT", "development"),
		FrontendURL:          getEnv("FRONTEND_URL", "http://localhost:3000"),
		LogLevel:             getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

// GetDSN returns PostgreSQL connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}
