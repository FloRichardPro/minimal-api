package services

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // mysql
	_ "github.com/lib/pq"              // postgres
)

type SQLConnectorConf struct {
	Driver string `mapstructure:"driver"`
	DSN    string `mapstructure:"dsn"`
}

func NewSQLConnector(c SQLConnectorConf) (*sql.DB, error) {
	sqlDB, err := sql.Open(c.Driver, c.DSN)
	if err != nil {
		return nil, fmt.Errorf("can't open connection to SqlDB(driver: %s): %w", c.Driver, err)
	}

	return sqlDB, nil
}
