package slogx

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"log/slog"
)

func TestNewFileLogger(t *testing.T) {
	filename := "test.log"
	if runtime.GOOS == "windows" {
		filename = "D:/test.log"
	}
	roller, _ := NewFileLogger(
		WithFilename(filename),
		WithMaxSize(2),
		WithMaxBackups(5),
		WithMaxAge(3),
		WithLocalTime(true),
		WithCompress(true),
		WithStdout(true),
	)
	// 断言路径分隔符兼容
	assert.Equal(t, filepath.Clean(filename), filepath.Clean(roller.Filename))
	assert.Equal(t, 2, roller.MaxSize)
	assert.Equal(t, 5, roller.MaxBackups)
	assert.Equal(t, 3, roller.MaxAge)
	assert.True(t, roller.LocalTime)
	assert.True(t, roller.Compress)
	// 清理
	_ = os.Remove(filename)
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
	str := buf.String()
	assert.Contains(t, str, "level=INFO")
	assert.Contains(t, str, "msg=\"hello world\"")
}
