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

func ErrNotSupport(driver Driver) error {
	return fmt.Errorf("not supported database <%s>", driver)
}

type Driver string

const (
	Sqlite    Driver = "sqlite"
	Postgres  Driver = "postgres"
	SqlServer Driver = "sqlserver"
	Mysql     Driver = "mysql"
)

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

func NewDB(dialector gorm.Dialector, opts ...gorm.Option) (*gorm.DB, error) {
	conn, err := gorm.Open(dialector, opts...)
	if err != nil {
		return nil, err
	}

	sql, err := conn.DB()
	if err != nil {
		return nil, err
	}

	maxOpen := 0
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
