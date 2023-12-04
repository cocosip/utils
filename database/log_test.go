package database

import (
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestNewGormLoggerWithSlog(t *testing.T) {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	_ = NewGormLoggerWithSlog(
		logger,
		WithSlowThreshold(time.Second),
		WithColorful(true),
		WithParameterizedQueries(true),
		WithSlogLevel(slog.LevelInfo))

	t.Log("new gorm logger with slog success")
}
