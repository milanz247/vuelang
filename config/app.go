package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// App is the central configuration struct for the entire framework.
// Call Load() once at startup; treat the result as immutable.
type App struct {
	// Server
	Port    string
	Env     string
	AppName string
	AppURL  string

	// Database
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	DBMaxOpenConns int
	DBMaxIdleConns int
	DBSeed         bool

	// JWT
	JWTSecret        string
	JWTAccessTTL     time.Duration
	JWTRefreshTTLDay int

	// CORS
	CORSAllowedOrigins []string

	// Rate Limiting
	RateLimitRequests    int
	RateLimitWindowSecs  int
	AuthRateLimitReqs    int
	AuthRateLimitWinSecs int

	// Mail
	MailDriver      string
	MailHost        string
	MailPort        int
	MailUsername    string
	MailPassword    string
	MailFromAddress string
	MailFromName    string
}

func Load() *App {
	cfg := &App{
		Port:    getEnv("PORT", "8080"),
		Env:     getEnv("ENV", "development"),
		AppName: getEnv("APP_NAME", "Vuelang"),
		AppURL:  getEnv("APP_URL", "http://localhost:8080"),

		DBHost:         getEnv("DB_HOST", "127.0.0.1"),
		DBPort:         getEnv("DB_PORT", "3306"),
		DBUser:         getEnv("DB_USER", "root"),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		DBName:         getEnv("DB_NAME", "vuelang"),
		DBMaxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 50),
		DBMaxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 25),
		DBSeed:         getEnv("DB_SEED", "false") == "true",

		JWTSecret:        getEnv("JWT_SECRET", ""),
		JWTAccessTTL:     time.Duration(getEnvInt("JWT_ACCESS_TTL_MINUTES", 15)) * time.Minute,
		JWTRefreshTTLDay: getEnvInt("JWT_REFRESH_TTL_DAYS", 7),

		CORSAllowedOrigins: strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173,http://localhost:8080"), ","),

		RateLimitRequests:    getEnvInt("RATE_LIMIT_REQUESTS", 60),
		RateLimitWindowSecs:  getEnvInt("RATE_LIMIT_WINDOW_SECONDS", 60),
		AuthRateLimitReqs:    getEnvInt("AUTH_RATE_LIMIT_REQUESTS", 5),
		AuthRateLimitWinSecs: getEnvInt("AUTH_RATE_LIMIT_WINDOW_SECONDS", 60),

		MailDriver:      getEnv("MAIL_DRIVER", "smtp"),
		MailHost:        getEnv("MAIL_HOST", ""),
		MailPort:        getEnvInt("MAIL_PORT", 587),
		MailUsername:    getEnv("MAIL_USERNAME", ""),
		MailPassword:    getEnv("MAIL_PASSWORD", ""),
		MailFromAddress: getEnv("MAIL_FROM_ADDRESS", "noreply@vuelang.dev"),
		MailFromName:    getEnv("MAIL_FROM_NAME", "Vuelang"),
	}

	cfg.validate()
	return cfg
}

func (c *App) IsProd() bool { return c.Env == "production" }
func (c *App) IsDev() bool  { return c.Env == "development" }

func (c *App) validate() {
	if c.IsProd() {
		if c.JWTSecret == "" || c.JWTSecret == "change-me-before-deploying-to-production" {
			log.Fatal("FATAL: JWT_SECRET must be a strong random value in production.\n" +
				"Generate one with:  openssl rand -base64 64")
		}
		if len(c.JWTSecret) < 32 {
			log.Fatal("FATAL: JWT_SECRET must be at least 32 characters long.")
		}
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
