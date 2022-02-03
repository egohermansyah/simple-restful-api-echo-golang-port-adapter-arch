package mysql

import (
	"database/sql"
	"fmt"
	"golang-vscode-setup/config"
	"golang-vscode-setup/util/logger"

	_ "github.com/go-sql-driver/mysql"
)

const LOG_IDENTIFIER = "REPOSITORY_MYSQL"

var log = logger.SetIdentifierField(LOG_IDENTIFIER)

func Connect(config config.MariaDBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Username, config.Password, config.Host, config.Port, config.Name)
	DBCon, err := sql.Open(config.Driver, dsn)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	return DBCon, err
}
