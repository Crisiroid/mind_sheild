package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	App AppConfig
	DB  DBConfig
	JWT JWTConfig
}

type AppConfig struct {
	Port string
	Env  string
}

type DBConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

type JWTConfig struct {
	Secret            string
	AccessExpiration  string
	RefreshExpiration string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		App: AppConfig{
			Port: GetEnv("APP_PORT", "8080"),
			Env:  GetEnv("APP_ENV", "development"),
		},
		DB: DBConfig{
			Host:     GetEnv("DB_HOST", "localhost"),
			Port:     GetEnv("DB_PORT", "5432"),
			Name:     GetEnv("DB_NAME", "psychology_app"),
			User:     GetEnv("DB_USER", "postgres"),
			Password: GetEnv("DB_PASSWORD", "postgres"),
			SSLMode:  GetEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:            GetEnv("JWT_SECRET", "your-secret-key-change-in-production"),
			AccessExpiration:  GetEnv("JWT_ACCESS_EXPIRATION", "24h"),
			RefreshExpiration: GetEnv("JWT_REFRESH_EXPIRATION", "720h"),
		},
	}
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
