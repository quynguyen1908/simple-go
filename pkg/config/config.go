package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Port string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	DBTimezone string

	JWTSecretKey string

	SMTPHost string
	SMTPPort int
	SMTPUser string
	SMTPPass string
	AppURL   string
}

func LoadConfig() *AppConfig {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	port, _ := strconv.Atoi(getEnvOrDefault("SMTP_PORT", "587"))

	return &AppConfig{
		Port: getEnvOrDefault("PORT", "8080"),

		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  os.Getenv("DB_SSLMODE"),
		DBTimezone: os.Getenv("DB_TIMEZONE"),

		JWTSecretKey: os.Getenv("JWT_SECRET_KEY"),

		SMTPHost: getEnvOrDefault("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort: port,
		SMTPUser: os.Getenv("SMTP_USER"),
		SMTPPass: os.Getenv("SMTP_PASS"),
		AppURL:   getEnvOrDefault("APP_URL", "http://localhost:8080"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultValue
}
