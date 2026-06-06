package config

import "os"

type Config struct {
	ServerPort string
	DBDriver   string
	DBSource   string
	JWTSecret  string
}

func LoadConfig() (*Config, error) {
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DBDriver:   getEnv("DB_DRIVER", "postgres"),
		DBSource:   getEnv("DATABASE_URL", "postgresql://postgres:postgres@localhost:5432/paygo?sslmode=disable"),
		JWTSecret:  getEnv("JWT_SECRET", "super-secret-paygo-key-change-in-production"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
