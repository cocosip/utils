package database

import (
	glog "gorm.io/gorm/logger"
	"log/slog"
	"time"
)

type GormLoggerOption func(c *glog.Config)

func WithSlowThreshold(slowThreshold time.Duration) GormLoggerOption {
	return func(c *glog.Config) {
		c.SlowThreshold = slowThreshold
	}
}

func WithColorful(colorful bool) GormLoggerOption {
	return func(c *glog.Config) {
		c.Colorful = colorful
	}
}

func WithIgnoreRecordNotFoundError(ignoreRecordNotFoundError bool) GormLoggerOption {
	return func(c *glog.Config) {
		c.IgnoreRecordNotFoundError = ignoreRecordNotFoundError
	}
}

func WithParameterizedQueries(parameterizedQueries bool) GormLoggerOption {
	return func(c *glog.Config) {
		c.ParameterizedQueries = parameterizedQueries
	}
}

func WithLogLevel(level glog.LogLevel) GormLoggerOption {
	return func(c *glog.Config) {
		c.LogLevel = level
	}
}

func WithSlogLevel(level slog.Level) GormLoggerOption {
	return func(c *glog.Config) {
		l := glog.Info
		switch level {
		case slog.LevelDebug, slog.LevelInfo:
			l = glog.Info
		case slog.LevelWarn:
			l = glog.Warn
		case slog.LevelError:
			l = glog.Error
		}
		c.LogLevel = l
	}
}

func newDefaultConfig() *glog.Config {
	c := &glog.Config{
		SlowThreshold:             500 * time.Millisecond,
		Colorful:                  true,
		IgnoreRecordNotFoundError: true,
		LogLevel:                  glog.Warn,
		ParameterizedQueries:      true,
	}
	return c
}

func NewGormLoggerWithSlog(logger *slog.Logger, opts ...GormLoggerOption) glog.Interface {
	c := newDefaultConfig()
	for _, opt := range opts {
		opt(c)
	}

	l := slog.LevelInfo
	switch c.LogLevel {
	case glog.Info:
		l = slog.LevelInfo
	case glog.Warn:
		l = slog.LevelWarn
	case glog.Error:
		l = slog.LevelError
	}

	return glog.New(slog.NewLogLogger(logger.Handler(), l), *c)
}
