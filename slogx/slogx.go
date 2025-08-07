package slogx

import (
	"io"
	"log/slog"
	"os"
	"runtime"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Config 日志文件配置结构体
// Fields:
//   - Filename: 日志文件路径
//   - MaxSize: 单个日志文件最大大小(MB)
//   - MaxAge: 日志文件最大保存天数
//   - MaxBackups: 最大备份数量
//   - LocalTime: 是否使用本地时间
//   - Compress: 是否压缩历史日志
//   - Stdout: 是否同时输出到控制台
// Description: 日志文件滚动与输出配置
// "Fields description in English for cross-platform compatibility."
type Config struct {
	Filename   string `json:"filename" yaml:"filename"` // Log file path
	MaxSize    int    `json:"maxsize" yaml:"maxsize"`   // Max file size (MB)
	MaxAge     int    `json:"maxage" yaml:"maxage"`     // Max age (days)
	MaxBackups int    `json:"maxbackups" yaml:"maxbackups"` // Max backup count
	LocalTime  bool   `json:"localtime" yaml:"localtime"`   // Use local time
	Compress   bool   `json:"compress" yaml:"compress"`     // Compress backups
	Stdout     bool   `json:"stdout" yaml:"stdout"`         // Output to stdout
}

// newDefaultConfig returns a pointer to the default Config.
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

// Option is a function that modifies Config.
type Option func(c *Config)

// WithConfig sets the entire Config.
func WithConfig(cfg Config) Option {
	return func(c *Config) {
		*c = cfg
	}
}

// WithFilename sets the log file path. Handles cross-platform path separator.
func WithFilename(filename string) Option {
	return func(config *Config) {
		if runtime.GOOS == "windows" {
			filename = strings.ReplaceAll(filename, "/", "\\")
		} else {
			filename = strings.ReplaceAll(filename, "\\", "/")
		}
		config.Filename = filename
	}
}

// WithMaxSize sets the max file size (MB).
func WithMaxSize(maxSize int) Option {
	return func(c *Config) {
		c.MaxSize = maxSize
	}
}

// WithMaxBackups sets the max backup count.
func WithMaxBackups(maxBackups int) Option {
	return func(c *Config) {
		c.MaxBackups = maxBackups
	}
}

// WithMaxAge sets the max age (days).
func WithMaxAge(maxAge int) Option {
	return func(c *Config) {
		c.MaxAge = maxAge
	}
}

// WithLocalTime sets whether to use local time.
func WithLocalTime(localTime bool) Option {
	return func(c *Config) {
		c.LocalTime = localTime
	}
}

// WithCompress sets whether to compress backups.
func WithCompress(compress bool) Option {
	return func(c *Config) {
		c.Compress = compress
	}
}

// WithStdout sets whether to output to stdout.
func WithStdout(stdout bool) Option {
	return func(c *Config) {
		c.Stdout = stdout
	}
}

// NewFileLogger creates a lumberjack.Logger and io.Writer for file logging.
// Params: opts ...Option - config options
// Returns: *lumberjack.Logger, io.Writer
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

// newDefaultHandlerOptions returns default slog.HandlerOptions.
func newDefaultHandlerOptions() *slog.HandlerOptions {
	return &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}
}

// NewSlogTextLogger creates a slog.Logger with text handler.
// Params: w io.Writer, opts ...func(*slog.HandlerOptions)
// Returns: *slog.Logger
func NewSlogTextLogger(w io.Writer, opts ...func(*slog.HandlerOptions)) *slog.Logger {
	o := newDefaultHandlerOptions()
	for _, fn := range opts {
		fn(o)
	}
	return slog.New(slog.NewTextHandler(w, o))
}

// NewSlogJsonLogger creates a slog.Logger with JSON handler.
// Params: w io.Writer, opts ...func(*slog.HandlerOptions)
// Returns: *slog.Logger
func NewSlogJsonLogger(w io.Writer, opts ...func(*slog.HandlerOptions)) *slog.Logger {
	o := newDefaultHandlerOptions()
	for _, fn := range opts {
		fn(o)
	}
	return slog.New(slog.NewJSONHandler(w, o))
}

// WithHandlerOptions sets slog.HandlerOptions.
func WithHandlerOptions(opts slog.HandlerOptions) func(*slog.HandlerOptions) {
	return func(o *slog.HandlerOptions) {
		*o = opts
	}
}

// GetSlogLevel parses string to slog.Level.
// Params: s string
// Returns: slog.Level
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
