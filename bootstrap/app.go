// Package bootstrap wires together all framework components.
// It is the single place where the dependency graph is constructed.
package bootstrap

import (
	"database/sql"
	"io/fs"
	"net/http"

	"vuelang/app/controllers"
	"vuelang/app/middleware"
	"vuelang/app/repositories"
	"vuelang/app/services"
	"vuelang/config"
	"vuelang/database/migrations"
	"vuelang/database/seeders"
	"vuelang/internal/framework/hash"
	jwtpkg "vuelang/internal/framework/jwt"
	"vuelang/internal/platform/database"
	"vuelang/internal/platform/logger"
	"vuelang/internal/server"
)

// Application holds all top-level dependencies.
type Application struct {
	cfg *config.App
	db  *sql.DB
	srv *server.Server
}

// New bootstraps the entire application, runs migrations, and returns
// a ready-to-start Application. It is called exactly once in main().
func New() *Application {
	cfg := config.Load()
	logger.Init(cfg.Env)

	// ── Database ──────────────────────────────────────────────────────────────
	db, err := database.NewMySQL(cfg)
	if err != nil {
		logger.Log.Warn("MySQL unavailable, running without database: " + err.Error())
	} else {
		logger.Log.Info("MySQL connected")
		if err := migrations.Run(db); err != nil {
			logger.Log.Error("migration failed: " + err.Error())
		}
		if cfg.DBSeed && cfg.IsDev() {
			if err := seeders.RunAll(db); err != nil {
				logger.Log.Error("seeder failed: " + err.Error())
			}
		}
	}

	// ── Framework services ────────────────────────────────────────────────────
	hasher := hash.NewBcrypt()
	jwtSvc := jwtpkg.NewService(cfg.JWTSecret, cfg.JWTAccessTTL, cfg.JWTRefreshTTLDay)

	// ── Repositories ──────────────────────────────────────────────────────────
	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)

	// ── Services ──────────────────────────────────────────────────────────────
	authSvc := services.NewAuthService(userRepo, roleRepo, hasher, jwtSvc, cfg)
	userSvc := services.NewUserService(userRepo, hasher)

	// ── Controllers ───────────────────────────────────────────────────────────
	authCtrl := controllers.NewAuthController(authSvc)
	userCtrl := controllers.NewUserController(userSvc)

	// ── Middleware ────────────────────────────────────────────────────────────
	authMiddleware := middleware.NewAuth(jwtSvc)

	// ── Server ────────────────────────────────────────────────────────────────
	srv := server.New(cfg, authCtrl, userCtrl, authMiddleware)

	return &Application{cfg: cfg, db: db, srv: srv}
}

// HTTPServer sets up routes and returns a configured *http.Server ready to start.
func (a *Application) HTTPServer(staticFS fs.FS) *http.Server {
	return a.srv.Build(staticFS)
}

// Close releases all resources (DB connection pool, etc.).
func (a *Application) Close() {
	if a.db != nil {
		_ = a.db.Close()
	}
}
