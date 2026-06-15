package server

import (
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"vuelang/app/controllers"
	"vuelang/app/middleware"
	"vuelang/config"
	"vuelang/internal/platform/logger"
)

// Server wires together the HTTP engine and all registered controllers.
type Server struct {
	cfg      *config.App
	authCtrl *controllers.AuthController
	userCtrl *controllers.UserController
	authMw   *middleware.AuthMiddleware
}

func New(
	cfg *config.App,
	authCtrl *controllers.AuthController,
	userCtrl *controllers.UserController,
	authMw *middleware.AuthMiddleware,
) *Server {
	return &Server{
		cfg:      cfg,
		authCtrl: authCtrl,
		userCtrl: userCtrl,
		authMw:   authMw,
	}
}

// Build wires routes and returns an *http.Server ready for ListenAndServe.
func (s *Server) Build(staticFS fs.FS) *http.Server {
	if s.cfg.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(requestLogger())
	r.Use(securityHeaders())
	r.Use(corsMiddleware(s.cfg.CORSAllowedOrigins))

	registerRoutes(r, s)

	if s.cfg.IsProd() {
		s.mountStaticProd(r, staticFS)
	} else {
		s.mountStaticDev(r, "http://localhost:5173")
	}

	logger.Log.Info("server built",
		slog.String("env", s.cfg.Env),
		slog.String("addr", "http://localhost:"+s.cfg.Port),
	)

	return &http.Server{
		Addr:         ":" + s.cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

// ── Middleware ─────────────────────────────────────────────────────────────────

func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Log.Info("request",
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("latency", time.Since(start)),
			slog.String("ip", c.ClientIP()),
		)
	}
}

func securityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "0") // disabled in favour of CSP
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		c.Header("Content-Security-Policy",
			"default-src 'self'; "+
				"script-src 'self' 'unsafe-inline'; "+
				"style-src 'self' 'unsafe-inline'; "+
				"img-src 'self' data: https:; "+
				"connect-src 'self'")
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.Header("Cache-Control", "no-store")
		}
		c.Next()
	}
}

func corsMiddleware(allowedOrigins []string) gin.HandlerFunc {
	originSet := make(map[string]bool, len(allowedOrigins))
	for _, o := range allowedOrigins {
		originSet[strings.TrimSpace(o)] = true
	}
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if originSet[origin] || originSet["*"] {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Requested-With")
			c.Header("Access-Control-Max-Age", "86400")
		}
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

// ── Static serving ────────────────────────────────────────────────────────────

func (s *Server) mountStaticProd(r *gin.Engine, staticFS fs.FS) {
	publicFS, err := fs.Sub(staticFS, "ui/dist")
	if err != nil {
		panic("ui/dist not in binary — run 'make build' first: " + err.Error())
	}
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/" || path == "" {
			serveIndex(c, publicFS)
			return
		}
		f, err := publicFS.Open(path[1:])
		if err == nil {
			defer f.Close()
			if st, err := f.Stat(); err == nil && !st.IsDir() {
				c.FileFromFS(path[1:], http.FS(publicFS))
				return
			}
		}
		serveIndex(c, publicFS)
	})
}

func serveIndex(c *gin.Context, publicFS fs.FS) {
	data, err := fs.ReadFile(publicFS, "index.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "index.html missing")
		return
	}
	c.Header("Cache-Control", "no-cache")
	c.Data(http.StatusOK, "text/html; charset=utf-8", data)
}

func (s *Server) mountStaticDev(r *gin.Engine, viteAddr string) {
	target, err := url.Parse(viteAddr)
	if err != nil {
		panic("invalid vite URL: " + err.Error())
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Director = func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host
	}
	r.NoRoute(func(c *gin.Context) {
		proxy.ServeHTTP(c.Writer, c.Request)
	})
}
