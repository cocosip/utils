package logging

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log/slog"
	"os"
)

type Config struct {
	Filename   string `json:"filename" yaml:"filename"`
	MaxSize    int    `json:"maxsize" yaml:"maxsize"`
	MaxAge     int    `json:"maxage" yaml:"maxage"`
	MaxBackups int    `json:"maxbackups" yaml:"maxbackups"`
	LocalTime  bool   `json:"localtime" yaml:"localtime"`
	Compress   bool   `json:"compress" yaml:"compress"`
	Stdout     bool   `json:"stdout" yaml:"stdout"`
}

func newDefaultConfig() *Config {
	return &Config{
		Filename:   "./logs/default.log",
		MaxSize:    10,
		MaxBackups: 3,
		MaxAge:     7,
		LocalTime:  true,
		Compress:   true,
		Stdout:     true,
	}
}

func NewFileLogger(opts ...func(*Config)) io.Writer {
	cfg := newDefaultConfig()
	for _, fn := range opts {
		fn(cfg)
	}
	out := &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		LocalTime:  cfg.LocalTime,
		Compress:   cfg.Compress,
	}

	if cfg.Stdout {
		return io.MultiWriter(os.Stdout, out)
	}
	return out
}

func newDefaultHandlerOptions() *slog.HandlerOptions {
	return &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}
}

func NewSlogTextLogger(w io.Writer, opts ...func(*slog.HandlerOptions)) *slog.Logger {
	o := newDefaultHandlerOptions()
	for _, fn := range opts {
		fn(o)
	}
	return slog.New(slog.NewTextHandler(w, o))
}

func NewSlogJsonLogger(w io.Writer, opts ...func(*slog.HandlerOptions)) *slog.Logger {
	o := newDefaultHandlerOptions()
	for _, fn := range opts {
		fn(o)
	}
	return slog.New(slog.NewJSONHandler(w, o))
}
