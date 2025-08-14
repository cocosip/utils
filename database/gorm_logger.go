package database

import (
	glog "gorm.io/gorm/logger"
	"log/slog"
	"time"
)

// GormLoggerOption defines a function type for customizing gorm logger configuration.
type GormLoggerOption func(c *glog.Config)

// WithSlowThreshold returns a GormLoggerOption that sets the slow SQL threshold.
// slowThreshold: the duration considered as slow SQL.
// Returns a GormLoggerOption.
func WithSlowThreshold(slowThreshold time.Duration) GormLoggerOption {
	return func(c *glog.Config) {
		c.SlowThreshold = slowThreshold
	}
}

// WithColorful returns a GormLoggerOption that enables or disables colorful log output.
// colorful: true to enable colorful output, false to disable.
// Returns a GormLoggerOption.
func WithColorful(colorful bool) GormLoggerOption {
	return func(c *glog.Config) {
		c.Colorful = colorful
	}
}

// WithIgnoreRecordNotFoundError returns a GormLoggerOption that sets whether to ignore record not found errors.
// ignoreRecordNotFoundError: true to ignore, false otherwise.
// Returns a GormLoggerOption.
func WithIgnoreRecordNotFoundError(ignoreRecordNotFoundError bool) GormLoggerOption {
	return func(c *glog.Config) {
		c.IgnoreRecordNotFoundError = ignoreRecordNotFoundError
	}
}

// WithParameterizedQueries returns a GormLoggerOption that enables or disables parameterized queries.
// parameterizedQueries: true to enable, false to disable.
// Returns a GormLoggerOption.
func WithParameterizedQueries(parameterizedQueries bool) GormLoggerOption {
	return func(c *glog.Config) {
		c.ParameterizedQueries = parameterizedQueries
	}
}

// WithLogLevel returns a GormLoggerOption that sets the gorm log level.
// level: the gorm log level to set.
// Returns a GormLoggerOption.
func WithLogLevel(level glog.LogLevel) GormLoggerOption {
	return func(c *glog.Config) {
		c.LogLevel = level
	}
}

// WithSlogLevel returns a GormLoggerOption that maps slog.Level to gorm LogLevel.
// level: the slog log level to map.
// Returns a GormLoggerOption.
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

// newDefaultConfig creates and returns a default gorm logger configuration.
// Returns a pointer to glog.Config.
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

// NewGormLoggerWithSlog creates a new gorm logger using slog.Logger and custom options.
// logger: the slog.Logger instance to use.
// opts: variadic GormLoggerOption for custom configuration.
// Returns a gorm logger implementing glog.Interface.
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