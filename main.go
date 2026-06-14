package main

import (
	"log"
	"os"

	"go-cloud-erp/database/migrations"
	"go-cloud-erp/internal/config"
	"go-cloud-erp/internal/platform/database"
	"go-cloud-erp/internal/platform/logger"
	"go-cloud-erp/internal/server"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg := config.Load()
	logger.Init(cfg.Env)

	// ── Database ──────────────────────────────────────────────────────────────
	db, err := database.NewMySQLConnection(cfg)
	if err != nil {
		logger.Log.Warn("MySQL unavailable, running without database: " + err.Error())
	} else {
		defer db.Close()
		logger.Log.Info("MySQL connected")

		if err := migrations.Run(db); err != nil {
			logger.Log.Error("migration failed: " + err.Error())
			os.Exit(1)
		}
	}

	// ── HTTP Server ───────────────────────────────────────────────────────────
	//
	//   make dev    →  ENV=development  →  proxies /* to Vite :5173
	//   make build  →  ENV=production   →  serves embedded ui/dist
	//
	srv := server.NewServer(cfg, db)

	if cfg.Env == "production" {
		if err := srv.Start(embeddedUI); err != nil {
			log.Fatalf("server: %v", err)
		}
	} else {
		if err := srv.StartDev("http://localhost:5173"); err != nil {
			log.Fatalf("dev server: %v", err)
		}
	}
}
