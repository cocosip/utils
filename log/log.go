package log

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log/slog"
	"os"
	"strings"
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

type Option func(c *Config)

func WithConfig(cfg Config) Option {
	return func(c *Config) {
		*c = cfg
	}
}

func WithFilename(filename string) Option {
	return func(config *Config) {
		config.Filename = filename
	}
}

func WithMaxSize(maxSize int) Option {
	return func(c *Config) {
		c.MaxSize = maxSize
	}
}

func WithMaxBackups(maxBackups int) Option {
	return func(c *Config) {
		c.MaxBackups = maxBackups
	}
}

func WithMaxAge(maxAge int) Option {
	return func(c *Config) {
		c.MaxAge = maxAge
	}
}

func WithLocalTime(localTime bool) Option {
	return func(c *Config) {
		c.LocalTime = localTime
	}
}

func WithCompress(compress bool) Option {
	return func(c *Config) {
		c.Compress = compress
	}
}

func WithStdout(stdout bool) Option {
	return func(c *Config) {
		c.Stdout = stdout
	}
}

func NewFileLogger(opts ...Option) (*lumberjack.Logger, io.Writer) {
	c := newDefaultConfig()
	for _, fn := range opts {
		fn(c)
	}
	roller := &lumberjack.Logger{
		Filename:   c.Filename,
		MaxSize:    c.MaxSize,
		MaxBackups: c.MaxBackups,
		MaxAge:     c.MaxAge,
		LocalTime:  c.LocalTime,
		Compress:   c.Compress,
	}

	if c.Stdout {
		return roller, io.MultiWriter(os.Stdout, roller)
	}
	return roller, roller
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

func WithHandlerOptions(opts slog.HandlerOptions) func(*slog.HandlerOptions) {
	return func(o *slog.HandlerOptions) {
		*o = opts
	}
}

func GetSlogLevel(s string) slog.Level {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	}

	return slog.LevelDebug
}
