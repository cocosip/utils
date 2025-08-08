package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"time"
)

// ErrNotSupport returns an error indicating that the specified database driver is not supported.
func ErrNotSupport(driver Driver) error {
	return fmt.Errorf("not supported database <%s>", driver)
}

// Driver represents a database driver type.
type Driver string

const (
	// Sqlite driver type.
	Sqlite Driver = "sqlite"
	// Postgres driver type.
	Postgres Driver = "postgres"
	// SqlServer driver type.
	SqlServer Driver = "sqlserver"
	// Mysql driver type.
	Mysql Driver = "mysql"
)

// NewDialector creates a new GORM Dialector based on the provided driver type and DSN.
// It supports SQLite, PostgreSQL, SQL Server, and MySQL.
// Returns an error if the driver is not supported.
func NewDialector(driver Driver, dsn string) (gorm.Dialector, error) {
	var dialect gorm.Dialector
	switch driver {
	case Sqlite:
		dialect = sqlite.Open(dsn)
	case Postgres:
		dialect = postgres.Open(dsn)
	case SqlServer:
		dialect = sqlserver.Open(dsn)
	case Mysql:
		dialect = mysql.Open(dsn)
	default:
		return dialect, ErrNotSupport(driver)
	}
	return dialect, nil
}

// NewDB initializes a new GORM database connection using the provided Dialector and options.
// It configures connection pool settings (MaxOpenConns, MaxIdleConns, ConnMaxLifetime, ConnMaxIdleTime)
// and pings the database to ensure connectivity.
// For SQLite, MaxOpenConns and MaxIdleConns are set to 1.
func NewDB(dialector gorm.Dialector, opts ...gorm.Option) (*gorm.DB, error) {
	conn, err := gorm.Open(dialector, opts...)
	if err != nil {
		return nil, err
	}

	sql, err := conn.DB()
	if err != nil {
		return nil, err
	}

	maxOpen := 0 // Default to unlimited
	maxIdle := 50

	if dialector.Name() == string(Sqlite) {
		maxOpen = 1
		maxIdle = 1
	}

	sql.SetMaxOpenConns(maxOpen)
	sql.SetMaxIdleConns(maxIdle)

	err = sql.Ping()
	if err != nil {
		return nil, err
	}

	sql.SetConnMaxLifetime(10 * time.Minute)
	sql.SetConnMaxIdleTime(10 * time.Minute)

	return conn, nil
}

// CloseDB closes the underlying sql.DB connection from a gorm.DB instance.
// db: the gorm.DB instance to close.
// Returns error if closing fails.
func CloseDB(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}