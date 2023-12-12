package log

import (
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
)

func TestNewFileLogger(t *testing.T) {
	_ = NewFileLogger(
		WithFilename("D:\\1.log"),
		WithMaxSize(2),
		WithMaxBackups(5),
		WithMaxAge(3),
		WithLocalTime(true),
		WithCompress(true),
		WithStdout(true),
	)

	t.Log("create new file logger success")
}

func TestGetSlogLevel(t *testing.T) {
	s1 := "Debug"
	assert.Equal(t, slog.LevelDebug, GetSlogLevel(s1))
	s2 := "error"
	assert.Equal(t, slog.LevelError, GetSlogLevel(s2))
}
