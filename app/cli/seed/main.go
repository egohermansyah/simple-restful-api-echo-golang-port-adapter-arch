package main

import (
	"golang-vscode-setup/app/cli/seed/role"
	"golang-vscode-setup/config"
	"golang-vscode-setup/repository/mysql"
	"golang-vscode-setup/util/gorm"
	"golang-vscode-setup/util/logger"
	"time"

	_ "github.com/go-sql-driver/mysql"
	logRotator "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

const LOG_IDENTIFIER = "APP_CLI_SEED"

var log = logger.SetIdentifierField(LOG_IDENTIFIER)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "timestamp",
			logrus.FieldKeyMsg:  "message",
		},
	})
	logrus.SetLevel(logrus.InfoLevel)

	logFile, err := logRotator.New(
		"resource/log/log.%Y%m%d%H%M",
		logRotator.WithLinkName("resource/log/log"),
		logRotator.WithMaxAge(time.Duration(86400)*time.Second),
		logRotator.WithRotationTime(time.Duration(86400)*time.Second),
	)
	if err == nil {
		logrus.SetOutput(logFile)
	} else {
		log.Error(err)
		panic(err)
	}
	defer logFile.Close()
	dbCon, err := mysql.Connect(config.GetConfig().MariaDbMasterDatabase)
	if err != nil {
		log.Error(err)
		panic(err)
	}

	gormDb, err := gorm.InitGorm(dbCon)

	role.Seed(gormDb)

	dbCon.Close()
}
