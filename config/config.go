package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all the configuration settings for the application.
type Config struct {
	DBURL          string
	JWTSecret      string
	Port           string
	AllowedOrigins string
}

// LoadConfig loads the configuration from the environment variables or a .env file.
func LoadConfig() *Config {
	// Load .env file if it exists (useful for local development)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading configurations from system environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Println("Warning: DB_URL environment variable is not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "supersecretkey" // Default secret for safety/fallback
		log.Println("Warning: JWT_SECRET environment variable is not set, using default fallback")
	}

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "*" // Default fallback
	}

	return &Config{
		DBURL:          dbURL,
		JWTSecret:      jwtSecret,
		Port:           port,
		AllowedOrigins: allowedOrigins,
	}
}
