package log

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
)

func TestNewFileLogger(t *testing.T) {
	roller, _ := NewFileLogger(
		WithFilename("D:\\1.log"),
		WithMaxSize(2),
		WithMaxBackups(5),
		WithMaxAge(3),
		WithLocalTime(true),
		WithCompress(true),
		WithStdout(true),
	)

	assert.Equal(t, "D:\\1.log", roller.Filename)
	assert.Equal(t, 2, roller.MaxSize)
	assert.Equal(t, 5, roller.MaxBackups)
	assert.Equal(t, 3, roller.MaxAge)
	assert.True(t, roller.LocalTime)
	assert.True(t, roller.Compress)

	t.Log("create new file logger success")
}

func TestGetSlogLevel(t *testing.T) {
	assert.Equal(t, slog.LevelDebug, GetSlogLevel("DEBUG"))
	assert.Equal(t, slog.LevelInfo, GetSlogLevel("INFO"))
	assert.Equal(t, slog.LevelWarn, GetSlogLevel("WARN"))
	assert.Equal(t, slog.LevelError, GetSlogLevel("ERROR"))
	assert.Equal(t, slog.LevelDebug, GetSlogLevel("invalid"))
}

func TestNewSlogLogger_WritesLog(t *testing.T) {
	var buf bytes.Buffer
	logger := NewSlogTextLogger(&buf)

	logger.Info("hello world")

	assert.Contains(t, buf.String(), "level=INFO")
	assert.Contains(t, buf.String(), "msg=\"hello world\"")
}
