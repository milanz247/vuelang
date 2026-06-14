package config

import (
	"log"
	"os"
)

type Config struct {
	Port       string
	Env        string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
}

// Load reads configuration from environment variables.
// Call this once at startup; the result is treated as immutable.
func Load() *Config {
	cfg := &Config{
		Port:       getEnv("PORT", "8080"),
		Env:        getEnv("ENV", "development"),
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "go_cloud_erp"),
		JWTSecret:  getEnv("JWT_SECRET", ""),
	}

	cfg.validate()
	return cfg
}

func (c *Config) validate() {
	if c.Env == "production" {
		// Hard-stop if JWT secret is empty or still set to the example value
		if c.JWTSecret == "" || c.JWTSecret == "change-me-before-deploying-to-production" {
			log.Fatal("FATAL: JWT_SECRET must be set to a strong random value in production.\n" +
				"Generate one with:  openssl rand -base64 48")
		}
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
