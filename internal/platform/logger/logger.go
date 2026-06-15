package logger

import (
	"log/slog"
	"os"
)

// Log is the global structured logger. Initialise once via Init().
var Log *slog.Logger

func Init(env string) {
	var handler slog.Handler
	if env == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})
	}
	Log = slog.New(handler)
	slog.SetDefault(Log)
}
