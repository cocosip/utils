package database

import (
	"log/slog"
	"os"
	"testing"
	"time"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

func TestNewDialectorAndNewDB(t *testing.T) {
	dial, err := NewDialector(Sqlite, ":memory:")
	if err != nil {
		t.Fatalf("NewDialector error: %v", err)
	}
	db, err := NewDB(dial)
	if err != nil {
		t.Fatalf("NewDB error: %v", err)
	}
	err = CloseDB(db)
	if err != nil {
		t.Fatalf("CloseDB error: %v", err)
	}
}

func TestErrNotSupport(t *testing.T) {
	_, err := NewDialector("unknown", "dsn")
	if err == nil {
		t.Error("should return error for unsupported driver")
	}
}
