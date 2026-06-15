package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"vuelang/bootstrap"
	"vuelang/internal/platform/logger"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file; silently ignored in production where real env vars are set.
	_ = godotenv.Load()

	app := bootstrap.New()
	defer app.Close()

	srv := app.HTTPServer(embeddedUI)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("server stopped: %v", err)
		}
	}()

	logger.Log.Info("✦ Vuelang V2 — press Ctrl+C to stop")
	<-ctx.Done()

	logger.Log.Info("shutdown signal received, draining connections…")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Log.Error("graceful shutdown error: " + err.Error())
		os.Exit(1)
	}
	logger.Log.Info("server stopped cleanly")
}
