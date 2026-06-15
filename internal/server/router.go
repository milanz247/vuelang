package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"vuelang/app/middleware"
	"vuelang/internal/framework/ratelimit"
	"vuelang/internal/framework/response"
)

// registerRoutes mounts all API routes onto the engine.
func registerRoutes(r *gin.Engine, s *Server) {
	// ── Health check ──────────────────────────────────────────────────────────
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{
			"status":    "ok",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"version":   "2.0.0",
		}, "healthy")
	})

	// ── API v1 ────────────────────────────────────────────────────────────────
	api := r.Group("/api/v1")

	// Rate limiter for general API traffic
	apiLimiter := ratelimit.New(s.cfg.RateLimitRequests,
		time.Duration(s.cfg.RateLimitWindowSecs)*time.Second)
	api.Use(middleware.RateLimit(apiLimiter))

	// Rate limiter for sensitive auth endpoints
	authLimiter := ratelimit.New(s.cfg.AuthRateLimitReqs,
		time.Duration(s.cfg.AuthRateLimitWinSecs)*time.Second)

	// ── Auth routes (public) ──────────────────────────────────────────────────
	auth := api.Group("/auth")
	auth.Use(middleware.RateLimit(authLimiter))
	{
		auth.POST("/register",        s.authCtrl.Register)
		auth.POST("/login",           s.authCtrl.Login)
		auth.POST("/refresh",         s.authCtrl.Refresh)
		auth.POST("/forgot-password", s.authCtrl.ForgotPassword)
		auth.POST("/reset-password",  s.authCtrl.ResetPassword)
	}

	// ── Protected routes (require valid JWT) ──────────────────────────────────
	protected := api.Group("/")
	protected.Use(s.authMw.Handle())
	{
		protected.POST("/auth/logout", s.authCtrl.Logout)
		protected.GET("/auth/me",      s.authCtrl.Me)

		// Users resource — admin only
		users := protected.Group("/users")
		users.Use(middleware.RequireRole("admin", "super_admin"))
		{
			users.GET("",        s.userCtrl.Index)
			users.GET("/:id",    s.userCtrl.Show)
			users.POST("",       s.userCtrl.Store)
			users.PUT("/:id",    s.userCtrl.Update)
			users.DELETE("/:id", s.userCtrl.Destroy)
		}
	}

	// ── 404 handler for /api routes ───────────────────────────────────────────
	r.NoRoute(func(c *gin.Context) {
		if len(c.Request.URL.Path) >= 4 && c.Request.URL.Path[:4] == "/api" {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "route not found"})
		}
	})
}
