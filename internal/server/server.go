package server

import (
	"database/sql"
	"io/fs"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"

	"go-cloud-erp/internal/config"
	"go-cloud-erp/internal/platform/logger"
)

// Server holds the HTTP engine and all shared dependencies.
type Server struct {
	router *gin.Engine
	cfg    *config.Config
	db     *sql.DB
}

func NewServer(cfg *config.Config, db *sql.DB) *Server {
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(requestLogger())
	r.Use(securityHeaders())

	return &Server{router: r, cfg: cfg, db: db}
}

// Start serves the embedded production build (ENV=production).
// staticFS must be the embed.FS containing ui/dist.
func (s *Server) Start(staticFS fs.FS) error {
	s.setupAPI()
	s.setupStaticProd(staticFS)

	port := s.port()
	logger.Log.Info("production server started", slog.String("addr", "http://localhost:"+port))
	return s.router.Run(":" + port)
}

// StartDev proxies non-API requests to the Vite dev server on viteURL.
func (s *Server) StartDev(viteURL string) error {
	s.setupAPI()
	s.setupStaticDev(viteURL)

	port := s.port()
	logger.Log.Info("dev server started",
		slog.String("addr", "http://localhost:"+port),
		slog.String("vite", viteURL),
	)
	return s.router.Run(":" + port)
}

func (s *Server) port() string {
	if s.cfg.Port != "" {
		return s.cfg.Port
	}
	return "8080"
}

// ── Middleware ────────────────────────────────────────────────────────────────

func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Log.Info("request",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("latency", time.Since(start)),
		)
	}
}

func securityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		// No-store cache for all API responses
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			c.Header("Cache-Control", "no-store")
		}
		c.Next()
	}
}
