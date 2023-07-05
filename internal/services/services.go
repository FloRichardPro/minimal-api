package services

import (
	"database/sql"
	"fmt"
)

var (
	Config                     Conf
	sqlConnectorSharedInstance *sql.DB
)

type Conf struct {
	SqlConnectorConf SQLConnectorConf `mapstructure:"sql"`
}

func Init() (err error) {
	sqlConnectorSharedInstance, err = NewSQLConnector(Config.SqlConnectorConf)
	if err != nil {
		return fmt.Errorf("can't create sql connector : %w", err)
	}

	return nil
}
